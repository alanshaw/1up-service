package ucan

import (
	provider_caps "github.com/alanshaw/1up-service/pkg/capabilities/provider"
	weight_caps "github.com/alanshaw/1up-service/pkg/capabilities/provider/weight"
	"github.com/alanshaw/1up-service/pkg/service"
	"github.com/alanshaw/1up-service/pkg/store/provider"
	ucan_errors "github.com/alanshaw/ucantone/errors"
	"github.com/alanshaw/ucantone/execution/bindexec"
	"github.com/alanshaw/ucantone/principal"
)

func NewProviderWeightSetHandler(id principal.Signer, providerStore provider.Store) *service.Handler {
	return &service.Handler{
		Capability: weight_caps.Set,
		Handler: bindexec.NewHandler(
			func(req *bindexec.Request[*weight_caps.SetArguments]) (*bindexec.Response[*weight_caps.SetOK], error) {
				args := req.Task().BindArguments()
				if req.Invocation().Issuer().DID() != id.DID() {
					return bindexec.NewResponse(bindexec.WithFailure[*weight_caps.SetOK](
						ucan_errors.New("Unauthorized", "only the service identity can set provider weight"),
					))
				}
				log.Infow("setting provider weight", "id", args.Provider, "weight", args.Weight)
				err := providerStore.Update(req.Context(), args.Provider, func(p provider_caps.Provider) (provider_caps.Provider, error) {
					p.Weight = args.Weight
					return p, nil
				})
				if err != nil {
					return nil, err
				}
				return bindexec.NewResponse(bindexec.WithSuccess(&weight_caps.SetOK{}))
			},
		),
	}
}
