package ucan

import (
	"fmt"
	"net/url"

	provider_caps "github.com/alanshaw/1up-service/pkg/capabilities/provider"
	"github.com/alanshaw/1up-service/pkg/service"
	"github.com/alanshaw/1up-service/pkg/store/provider"
	"github.com/alanshaw/ucantone/errors"
	"github.com/alanshaw/ucantone/execution/bindexec"
	"github.com/alanshaw/ucantone/principal"
)

func NewProviderRegisterHandler(id principal.Signer, providerStore provider.Store) *service.Handler {
	return &service.Handler{
		Capability: provider_caps.Register,
		Handler: bindexec.NewHandler(
			func(req *bindexec.Request[*provider_caps.RegisterArguments]) (*bindexec.Response[*provider_caps.RegisterOK], error) {
				args := req.Task().BindArguments()
				_, err := url.Parse(args.Endpoint)
				if err != nil {
					return bindexec.NewResponse(bindexec.WithFailure[*provider_caps.RegisterOK](
						errors.New("InvalidEndpoint", fmt.Sprintf("parsing endpoint: %s", err.Error())),
					))
				}
				proof, ok := req.Metadata().Delegation(args.Proof)
				if !ok {
					return bindexec.NewResponse(bindexec.WithFailure[*provider_caps.RegisterOK](
						errors.New("MissingProof", "missing proof delegation"),
					))
				}
				if proof.Issuer().DID() != args.Provider {
					return bindexec.NewResponse(bindexec.WithFailure[*provider_caps.RegisterOK](
						errors.New("InvalidProof", "proof issuer does not match provider DID"),
					))
				}
				if proof.Audience().DID() != req.Principal().DID() {
					return bindexec.NewResponse(bindexec.WithFailure[*provider_caps.RegisterOK](
						errors.New("InvalidProof", "proof audience does not match request principal DID"),
					))
				}

				log.Infow(
					"registering storage provider",
					"id", args.Provider,
					"endpoint", args.Endpoint,
					"proof", args.Proof.String(),
				)
				return bindexec.NewResponse(bindexec.WithSuccess(req.Task().BindArguments()))
			},
		),
	}
}
