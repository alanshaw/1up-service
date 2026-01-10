package ucan

import (
	"github.com/alanshaw/1up-service/pkg/service"
	"github.com/alanshaw/1up-service/pkg/store/provider"
	provider_caps "github.com/alanshaw/libracha/capabilities/provider"
	"github.com/alanshaw/ucantone/errors"
	"github.com/alanshaw/ucantone/execution/bindexec"
	"github.com/alanshaw/ucantone/principal"
)

func NewProviderListHandler(id principal.Signer, providerStore provider.Store) *service.Handler {
	return &service.Handler{
		Capability: provider_caps.List,
		Handler: bindexec.NewHandler(
			func(req *bindexec.Request[*provider_caps.ListArguments], res *bindexec.Response[*provider_caps.ListOK]) error {
				if req.Invocation().Issuer().DID() != id.DID() {
					return res.SetFailure(errors.New("Unauthorized", "only the service identity can list providers"))
				}
				providers := []provider_caps.Provider{}
				for p, err := range providerStore.List(req.Context()) {
					if err != nil {
						return err
					}
					providers = append(providers, p)
				}
				return res.SetSuccess(&provider_caps.ListOK{Providers: providers})
			},
		),
	}
}
