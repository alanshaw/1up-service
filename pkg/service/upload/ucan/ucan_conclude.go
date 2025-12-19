package ucan

import (
	"github.com/alanshaw/1up-service/pkg/service"
	ucan_caps "github.com/alanshaw/libracha/capabilities/ucan"
	"github.com/alanshaw/ucantone/execution/bindexec"
	logging "github.com/ipfs/go-log/v2"
)

var ucanConcludeLog = logging.Logger("service/upload/ucan/ucan/conclude")

func NewUCANConcludeHandler() *service.Handler {
	return &service.Handler{
		Capability: ucan_caps.Conclude,
		Handler: bindexec.NewHandler(
			func(req *bindexec.Request[*ucan_caps.ConcludeArguments]) (*bindexec.Response[*ucan_caps.ConcludeOK], error) {
				ucanConcludeLog.Debugf("%+v", req.Task().BindArguments())
				return bindexec.NewResponse(bindexec.WithSuccess(&ucan_caps.ConcludeOK{}))
			},
		),
	}
}
