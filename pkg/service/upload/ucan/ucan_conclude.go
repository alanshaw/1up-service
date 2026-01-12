package ucan

import (
	"context"
	"fmt"

	"github.com/alanshaw/1up-service/pkg/service"
	"github.com/alanshaw/1up-service/pkg/service/routing"
	delegation_store "github.com/alanshaw/1up-service/pkg/store/delegation"
	"github.com/alanshaw/1up-service/pkg/store/token"
	blob_caps "github.com/alanshaw/libracha/capabilities/blob"
	http_caps "github.com/alanshaw/libracha/capabilities/http"
	ucan_caps "github.com/alanshaw/libracha/capabilities/ucan"
	ucanlib "github.com/alanshaw/libracha/ucan"
	"github.com/alanshaw/ucantone/client"
	"github.com/alanshaw/ucantone/did"
	"github.com/alanshaw/ucantone/errors"
	"github.com/alanshaw/ucantone/execution"
	"github.com/alanshaw/ucantone/execution/bindexec"
	"github.com/alanshaw/ucantone/principal"
	"github.com/alanshaw/ucantone/ucan"
	"github.com/alanshaw/ucantone/ucan/invocation"
	"github.com/alanshaw/ucantone/ucan/promise"
	logging "github.com/ipfs/go-log/v2"
)

var ucanConcludeLog = logging.Logger("service/upload/ucan/ucan/conclude")

func NewUCANConcludeHandler(id principal.Signer, router *routing.Router, tokens token.Store) *service.Handler {
	return &service.Handler{
		Capability: ucan_caps.Conclude,
		Handler: bindexec.NewHandler(
			func(req *bindexec.Request[*ucan_caps.ConcludeArguments], res *bindexec.Response[*ucan_caps.ConcludeOK]) error {
				args := req.Task().BindArguments()
				var receipt ucan.Receipt
				for _, r := range req.Metadata().Receipts() {
					if r.Link() == args.Receipt {
						receipt = r
						break
					}
				}
				if receipt == nil {
					return res.SetFailure(errors.New("MissingReceipt", "missing receipt for conclude"))
				}

				// find the allocation target - the storage provider for the blob
				var space did.DID
				var storageProvider ucan.Principal
				var putTask *bindexec.Task[*http_caps.PutArguments]
				for _, inv := range req.Metadata().Invocations() {
					task := inv.Task()
					if receipt.Ran() == task.Link() && inv.Command() == http_caps.PutCommand {
						t, err := bindexec.NewTask[*http_caps.PutArguments](task.Subject(), task.Command(), task.Arguments(), task.Nonce())
						if err != nil {
							return err
						}
						putTask = t
						blobAllocateTaskLink := putTask.BindArguments().Destination.Task
						for ct, err := range tokens.FindByTask(req.Context(), blobAllocateTaskLink) {
							if err != nil {
								return err
							}
							for _, inv := range ct.Invocations() {
								if inv.Command() == blob_caps.AllocateCommand {
									storageProvider = inv.Audience()
									if storageProvider == nil {
										return res.SetFailure(errors.New("MissingAudience", "allocate invocation did not specify audience, and subject must be space"))
									}
									space = inv.Subject().DID()
									break
								}
							}
						}
						break
					}
				}
				if storageProvider == nil {
					return res.SetFailure(errors.New("InvalidConclusion", "could not find storage provider blob was allocated with"))
				}

				pstore := delegation_store.NewMapDelegationStore(req.Metadata().Delegations())
				matcher := ucanlib.NewDelegationMatcher(pstore)
				proofs, proofLinks, err := ucanlib.ProofChain(req.Context(), matcher, id, blob_caps.AcceptCommand, space)
				if err != nil {
					return err
				}
				if len(proofs) == 0 {
					return res.SetFailure(errors.New("MissingProofs", fmt.Sprintf("missing proof(s) for %q invocation", blob_caps.AcceptCommand)))
				}

				storageProviderInfo, err := router.Provider(req.Context(), storageProvider)
				if err != nil {
					if errors.Is(err, routing.ErrNotFound) {
						ucanConcludeLog.Errorw("storage provider not found", "provider", storageProvider.DID())
						return res.SetFailure(errors.New("ProviderNotFound", "storage provider not found"))
					}
					return err
				}
				accInv, err := blob_caps.Accept.Invoke(
					id,
					space,
					&blob_caps.AcceptArguments{
						Blob: putTask.BindArguments().Body,
						Put:  promise.AwaitOK{Task: putTask.Link()},
					},
					invocation.WithAudience(storageProvider),
					invocation.WithProofs(proofLinks...),
					invocation.WithNoNonce(),
				)
				if err != nil {
					return err
				}
				c, err := client.NewHTTP(
					storageProviderInfo.Endpoint,
					client.WithEventListener(&requestLogger{tokens}),
				)
				if err != nil {
					return err
				}
				// We don't care about the response - it is logged. The client can
				// inspect the receipt later.
				_, err = c.Execute(execution.NewRequest(req.Context(), accInv, execution.WithDelegations(proofs...)))
				if err != nil {
					return err
				}

				return res.SetSuccess(&ucan_caps.ConcludeOK{})
			},
		),
	}
}

type requestLogger struct {
	tokens token.Store
}

func (rl *requestLogger) OnRequestEncode(ctx context.Context, ct ucan.Container) error {
	return rl.tokens.Put(ctx, ct)
}

func (rl *requestLogger) OnResponseDecode(ctx context.Context, ct ucan.Container) error {
	return rl.tokens.Put(ctx, ct)
}

var _ client.RequestEncodeListener = (*requestLogger)(nil)
var _ client.ResponseDecodeListener = (*requestLogger)(nil)
