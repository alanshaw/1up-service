package ucan

import (
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"

	"github.com/alanshaw/1up-service/pkg/service"
	"github.com/alanshaw/1up-service/pkg/service/router"
	"github.com/alanshaw/1up-service/pkg/store/delegation"
	blob_caps "github.com/alanshaw/libracha/capabilities/blob"
	http_caps "github.com/alanshaw/libracha/capabilities/http"
	"github.com/alanshaw/libracha/digestutil"
	ucanlib "github.com/alanshaw/libracha/ucan"
	"github.com/alanshaw/ucantone/client"
	"github.com/alanshaw/ucantone/execution"
	"github.com/alanshaw/ucantone/execution/bindexec"
	"github.com/alanshaw/ucantone/ipld"
	"github.com/alanshaw/ucantone/ipld/datamodel"
	"github.com/alanshaw/ucantone/principal"
	ucan_ed "github.com/alanshaw/ucantone/principal/ed25519"
	"github.com/alanshaw/ucantone/result"
	"github.com/alanshaw/ucantone/ucan"
	"github.com/alanshaw/ucantone/ucan/container"
	"github.com/alanshaw/ucantone/ucan/invocation"
	"github.com/alanshaw/ucantone/ucan/promise"
	"github.com/alanshaw/ucantone/ucan/receipt"
	logging "github.com/ipfs/go-log/v2"
	"github.com/multiformats/go-multihash"
)

var blobAddLog = logging.Logger("service/upload/ucan" + blob_caps.AddCommand)

func NewBlobAddHandler(id principal.Signer, rt *router.Router) *service.Handler {
	return &service.Handler{
		Capability: blob_caps.Add,
		Handler: bindexec.NewHandler(
			func(req *bindexec.Request[*blob_caps.AddArguments]) (*bindexec.Response[*blob_caps.AddOK], error) {
				args := req.Task().BindArguments()
				space := req.Invocation().Subject()
				blob := args.Blob
				cause := req.Invocation().Task().Link()
				pstore := delegation.NewMapDelegationStore(req.Metadata().Delegations())
				log := blobAddLog.With("space", space.DID(), "digest", digestutil.Format(blob.Digest))

				provider, allocInv, allocRcpt, allocOK, err := doAllocate(req.Context(), id, rt, pstore, space, blob, cause)
				if err != nil {
					if errors.Is(err, router.ErrCandidateUnavailable) {
						log.Errorw("unable to select a storage provider", "error", err)
						return bindexec.NewResponse(bindexec.WithFailure[*blob_caps.AddOK](err))
					}
					return nil, err
				}

				putInv, putRcpt, err := genPut(blob, allocInv, allocOK)
				if err != nil {
					return nil, err
				}

				accInv, accRcpt, err := maybeAccept(req.Context(), id, provider, space, pstore, blob, putInv, putRcpt)
				if err != nil {
					return nil, err
				}

				invocations := []ucan.Invocation{allocInv, putInv}
				if accRcpt != nil {
					// if the accept receipt has been issued, add the issued invocation to
					// the response.
					invocations = append(invocations, accInv)
				}

				receipts := []ucan.Receipt{allocRcpt}
				if putRcpt != nil {
					receipts = append(receipts, putRcpt)
				}
				if accRcpt != nil {
					receipts = append(receipts, accRcpt)
				}

				meta := container.New(
					container.WithInvocations(invocations...),
					container.WithReceipts(receipts...),
				)

				return bindexec.NewResponse(
					bindexec.WithSuccess(&blob_caps.AddOK{
						Site: promise.AwaitOK{Task: accInv.Task().Link()},
					}),
					bindexec.WithMetadata[*blob_caps.AddOK](meta),
				)
			},
		),
	}
}

func doAllocate(
	ctx context.Context,
	id principal.Signer,
	rt *router.Router,
	pstore delegation.Store,
	space ucan.Subject,
	blob blob_caps.Blob,
	cause ucan.Link,
) (router.ProviderInfo, ucan.Invocation, ucan.Receipt, blob_caps.AllocateOK, error) {
	log := blobAddLog.With("space", space.DID(), "digest", digestutil.Format(blob.Digest))

	var exclusions []ucan.Principal
	for {
		candidate, err := rt.Select(ctx, blob.Digest, blob.Size, router.WithExclusions(exclusions...))
		if err != nil {
			return router.ProviderInfo{}, nil, nil, blob_caps.AllocateOK{}, err
		}
		log := log.With("candidate", candidate.ID)
		log.Infow("selected storage provider candidate")

		proofs, proofLinks, err := ucanlib.ProofChain(ctx, pstore, id, blob_caps.AllocateCommand, space)
		if err != nil {
			log.Warnw("failed to construct proof chain", "error", err)
			continue
		}

		if len(proofs) == 0 {
			log.Warnw("no proofs found for selected provider, excluding")
			exclusions = append(exclusions, candidate.ID)
			continue
		}

		inv, err := blob_caps.Allocate.Invoke(
			id,
			space,
			&blob_caps.AllocateArguments{
				Blob:  blob_caps.Blob(blob),
				Cause: cause,
			},
			invocation.WithProofs(proofLinks...),
		)
		if err != nil {
			return router.ProviderInfo{}, nil, nil, blob_caps.AllocateOK{}, err
		}

		c, err := client.NewHTTP(candidate.Endpoint)
		if err != nil {
			return router.ProviderInfo{}, nil, nil, blob_caps.AllocateOK{}, err
		}

		res, err := c.Execute(execution.NewRequest(ctx, inv, execution.WithProofs(proofs...)))
		if err != nil {
			log.Errorw("executing allocation invocation", "error", err)
			exclusions = append(exclusions, candidate.ID)
			continue
		}

		o, x := result.Unwrap(res.Result())
		if x != nil {
			log.Errorw("failure result for allocation", "error", x)
			exclusions = append(exclusions, candidate.ID)
			continue
		}

		var allocOK blob_caps.AllocateOK
		err = datamodel.Rebind(datamodel.NewAny(o), &allocOK)
		if err != nil {
			log.Errorw("rebinding allocation result", "error", err)
			exclusions = append(exclusions, candidate.ID)
			continue
		}

		rcpt, ok := res.Metadata().Receipt(inv.Task().Link())
		if !ok {
			log.Errorw("missing receipt for allocation task")
			exclusions = append(exclusions, candidate.ID)
			continue
		}

		return candidate, inv, rcpt, allocOK, nil
	}
}

// Generates an invocation to put the blob to the storage provider. It MAY
// return a receipt if the allocation result indicates that the provider already
// has the blob.
func genPut(blob blob_caps.Blob, allocInv ucan.Invocation, allocOK blob_caps.AllocateOK) (ucan.Invocation, ucan.Receipt, error) {
	log := blobAddLog.With(
		"space", allocInv.Subject().DID(),
		"digest", digestutil.Format(blob.Digest),
		"provider", allocInv.Audience().DID(),
	)
	log.Info("generating put invocation")

	// Derive the principal that will provide the blob from the blob digest.
	// we do this so that any actor with a blob could issue a receipt for the
	// `/http/put` invocation.
	blobProvider, err := deriveDID(blob.Digest)
	if err != nil {
		return nil, nil, err
	}

	putInv, err := http_caps.Put.Invoke(
		blobProvider,
		blobProvider,
		&http_caps.PutArguments{
			Body: blob,
			Destination: promise.AwaitOK{
				Task: allocInv.Task().Link(),
			},
		},
		// We encode the keys for the blob provider principal that can be used
		// by the client to use in order to sign a receipt. Client could
		// actually derive the same principal from the blob digest like we did
		// above, however by embedding the keys we make API more flexible and
		// could in the future generate one-off principals instead.
		invocation.WithMetadata(ipld.Map{
			"keys": ipld.Map{
				blobProvider.DID().String(): blobProvider.Bytes(),
			},
		}),
		// We use nonce-less invocation to make structure deterministic.
		invocation.WithNoNonce(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("invoking %q: %w", http_caps.PutCommand, err)
	}

	var putRcpt ucan.Receipt

	// If no address was provided we have a blob in store already and we can issue
	// a receipt for the `/http/put` without requiring blob to be provided.
	if allocOK.Address == nil {
		log.Infow("blob present on provider, issuing receipt for put")
		putRcpt, err = receipt.Issue(
			blobProvider,
			putInv.Task().Link(),
			result.OK[ipld.Map, ipld.Any](ipld.Map{}),
		)
		if err != nil {
			return nil, nil, fmt.Errorf("issuing %q receipt: %w", http_caps.PutCommand, err)
		}
	}

	return putInv, putRcpt, nil
}

// Derives did:key principal from (blob) multihash that can be used to
// sign ucan invocations/receipts for the the subject (blob) multihash.
func deriveDID(digest multihash.Multihash) (ucan_ed.Ed25519Signer, error) {
	if len(digest) < 32 {
		return nil, fmt.Errorf("expected []byte with length %d, got %d", ed25519.SeedSize, len(digest))
	}
	seed := digest[len(digest)-32:]
	pk := ed25519.NewKeyFromSeed(seed)
	return ucan_ed.FromRaw(pk[:ed25519.SeedSize])
}

// maybeAccept generates and possibly executes a `/blob/accept` invocation if
// the provided put receipt is non-nil and non-failure.
func maybeAccept(
	ctx context.Context,
	id principal.Signer,
	provider router.ProviderInfo,
	space ucan.Principal,
	pstore delegation.Store,
	blob blob_caps.Blob,
	putInv ucan.Invocation,
	putRcpt ucan.Receipt,
) (ucan.Invocation, ucan.Receipt, error) {
	log := blobAddLog.With(
		"space", space.DID(),
		"digest", digestutil.Format(blob.Digest),
		"provider", provider.ID.DID(),
	)
	log.Info("generating accept invocation")

	proofs, proofLinks, err := ucanlib.ProofChain(ctx, pstore, id, blob_caps.AcceptCommand, space)
	if err != nil {
		return nil, nil, err
	}

	allocInv, err := blob_caps.Accept.Invoke(
		id,
		space,
		&blob_caps.AcceptArguments{
			Blob: blob_caps.Blob(blob),
			Put: promise.AwaitOK{
				Task: putInv.Task().Link(),
			},
		},
		invocation.WithProofs(proofLinks...),
	)
	if err != nil {
		return nil, nil, err
	}

	var allocRcpt ucan.Receipt

	// If put has already succeeded, we can execute `/blob/accept` right away.
	if putRcpt != nil {
		_, x := result.Unwrap(putRcpt.Out())
		if x == nil {
			c, err := client.NewHTTP(provider.Endpoint)
			if err != nil {
				return nil, nil, err
			}

			res, err := c.Execute(execution.NewRequest(ctx, allocInv, execution.WithProofs(proofs...)))
			if err != nil {
				return nil, nil, err
			}

			rcpt, ok := res.Metadata().Receipt(allocInv.Task().Link())
			if !ok {
				log.Errorw("missing receipt for allocation task")
				return nil, nil, err
			}
			allocRcpt = rcpt

			// TODO: add to blob registry
		}
	}

	return allocInv, allocRcpt, nil
}
