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
	logging "github.com/ipfs/go-log/v2"
)

var provRegLog = logging.Logger("service/upload/ucan" + provider_caps.RegisterCommand)

func NewProviderRegisterHandler(id principal.Signer, providerStore provider.Store) *service.Handler {
	return &service.Handler{
		Capability: provider_caps.Register,
		Handler: bindexec.NewHandler(
			func(req *bindexec.Request[*provider_caps.RegisterArguments]) (*bindexec.Response[*provider_caps.RegisterOK], error) {
				args := req.Task().BindArguments()
				endpoint, err := url.Parse(args.Endpoint)
				if err != nil {
					return bindexec.NewResponse(bindexec.WithFailure[*provider_caps.RegisterOK](
						errors.New("InvalidEndpoint", fmt.Sprintf("parsing endpoint: %s", err.Error())),
					))
				}
				if req.Invocation().Issuer().DID() != id.DID() && req.Invocation().Issuer().DID() != args.Provider {
					return bindexec.NewResponse(bindexec.WithFailure[*provider_caps.RegisterOK](
						errors.New("Unauthorized", "only the service identity or the provider itself can register a provider"),
					))
				}
				provRegLog.Infow(
					"registering storage provider",
					"id", args.Provider,
					"endpoint", args.Endpoint,
				)
				err = providerStore.Put(req.Context(), args.Provider, endpoint)
				if err != nil {
					return nil, err
				}
				return bindexec.NewResponse(bindexec.WithSuccess(&provider_caps.RegisterOK{}))
			},
		),
	}
}
