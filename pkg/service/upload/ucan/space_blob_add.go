package ucan

import (
	space_blob_caps "github.com/alanshaw/1up-service/pkg/capabilities/space/blob"
	"github.com/alanshaw/1up-service/pkg/service"
	"github.com/alanshaw/1up-service/pkg/service/router"
	"github.com/alanshaw/ucantone/execution/bindexec"
	"github.com/alanshaw/ucantone/principal"
)

func NewSpaceBlobAddHandler(id principal.Signer, router *router.Router) *service.Handler {
	return &service.Handler{
		Capability: space_blob_caps.Add,
		Handler: bindexec.NewHandler(
			func(req *bindexec.Request[*space_blob_caps.AddArguments]) (*bindexec.Response[*space_blob_caps.AddOK], error) {
				// TODO
				return bindexec.NewResponse(bindexec.WithSuccess(&space_blob_caps.AddOK{}))
			},
		),
	}
}
