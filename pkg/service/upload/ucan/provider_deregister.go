package ucan

import (
	"github.com/alanshaw/1up-service/pkg/service"
	"github.com/alanshaw/1up-service/pkg/store/provider"
	provider_caps "github.com/alanshaw/libracha/capabilities/provider"
	"github.com/alanshaw/ucantone/errors"
	"github.com/alanshaw/ucantone/execution/bindexec"
	"github.com/alanshaw/ucantone/principal"
	logging "github.com/ipfs/go-log/v2"
)

var provDeregLog = logging.Logger("service/upload/ucan" + provider_caps.DeregisterCommand)

func NewProviderDeregisterHandler(id principal.Signer, providerStore provider.Store) *service.Handler {
	return &service.Handler{
		Capability: provider_caps.Deregister,
		Handler: bindexec.NewHandler(
			func(req *bindexec.Request[*provider_caps.DeregisterArguments]) (*bindexec.Response[*provider_caps.DeregisterOK], error) {
				args := req.Task().BindArguments()
				if req.Invocation().Issuer().DID() != id.DID() && req.Invocation().Issuer().DID() != args.Provider {
					return bindexec.NewResponse(bindexec.WithFailure[*provider_caps.DeregisterOK](
						errors.New("Unauthorized", "only the service identity or the provider itself can deregister a provider"),
					))
				}
				provDeregLog.Infow("deregistering storage provider", "id", args.Provider)
				err := providerStore.Del(req.Context(), args.Provider)
				if err != nil {
					return nil, err
				}
				return bindexec.NewResponse(bindexec.WithSuccess(&provider_caps.DeregisterOK{}))
			},
		),
	}
}
