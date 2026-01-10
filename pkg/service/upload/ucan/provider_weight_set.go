package ucan

import (
	"github.com/alanshaw/1up-service/pkg/service"
	"github.com/alanshaw/1up-service/pkg/store/provider"
	provider_caps "github.com/alanshaw/libracha/capabilities/provider"
	weight_caps "github.com/alanshaw/libracha/capabilities/provider/weight"
	ucan_errors "github.com/alanshaw/ucantone/errors"
	"github.com/alanshaw/ucantone/execution/bindexec"
	"github.com/alanshaw/ucantone/principal"
	logging "github.com/ipfs/go-log/v2"
)

var provWeightSetLog = logging.Logger("service/upload/ucan" + weight_caps.SetCommand)

func NewProviderWeightSetHandler(id principal.Signer, providerStore provider.Store) *service.Handler {
	return &service.Handler{
		Capability: weight_caps.Set,
		Handler: bindexec.NewHandler(
			func(req *bindexec.Request[*weight_caps.SetArguments], res *bindexec.Response[*weight_caps.SetOK]) error {
				args := req.Task().BindArguments()
				if req.Invocation().Issuer().DID() != id.DID() {
					return res.SetFailure(ucan_errors.New("Unauthorized", "only the service identity can set provider weight"))
				}
				provWeightSetLog.Infow("setting provider weight", "id", args.Provider, "weight", args.Weight)
				err := providerStore.Update(req.Context(), args.Provider, func(p provider_caps.Provider) (provider_caps.Provider, error) {
					p.Weight = args.Weight
					return p, nil
				})
				if err != nil {
					return err
				}
				return res.SetSuccess(&weight_caps.SetOK{})
			},
		),
	}
}
