package ucan

import (
	"github.com/alanshaw/1up-service/pkg/service"
	"github.com/alanshaw/libracha/capabilities/debug"
	"github.com/alanshaw/ucantone/execution/bindexec"
	logging "github.com/ipfs/go-log/v2"
)

var echoDebugLog = logging.Logger("service/upload/ucan" + debug.EchoCommand)

func NewDebugEchoHandler() *service.Handler {
	return &service.Handler{
		Capability: debug.Echo,
		Handler: bindexec.NewHandler(
			func(req *bindexec.Request[*debug.EchoArguments]) (*bindexec.Response[*debug.EchoOK], error) {
				echoDebugLog.Debugf("%+v", req.Task().BindArguments())
				return bindexec.NewResponse(bindexec.WithSuccess(req.Task().BindArguments()))
			},
		),
	}
}
