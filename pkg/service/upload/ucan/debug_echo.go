package ucan

import (
	"github.com/alanshaw/1up-service/pkg/capabilities/debug"
	"github.com/alanshaw/1up-service/pkg/service"
	"github.com/alanshaw/ucantone/execution/bindexec"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("service/upload/ucan")

func NewDebugEchoHandler() *service.Handler {
	return &service.Handler{
		Capability: debug.Echo,
		Handler: bindexec.NewHandler(
			func(req *bindexec.Request[*debug.EchoArguments]) (*bindexec.Response[*debug.EchoOK], error) {
				log.Debugf("%+v", req.Task().BindArguments())
				return bindexec.NewResponse(bindexec.WithSuccess(req.Task().BindArguments()))
			},
		),
	}
}
