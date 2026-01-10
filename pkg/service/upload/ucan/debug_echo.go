package ucan

import (
	"github.com/alanshaw/1up-service/pkg/service"
	"github.com/alanshaw/libracha/capabilities/debug"
	"github.com/alanshaw/ucantone/execution/bindexec"
	"github.com/alanshaw/ucantone/principal"
	logging "github.com/ipfs/go-log/v2"
)

var echoDebugLog = logging.Logger("service/upload/ucan" + debug.EchoCommand)

func NewDebugEchoHandler(id principal.Signer) *service.Handler {
	return &service.Handler{
		Capability: debug.Echo,
		Handler: bindexec.NewHandler(
			func(req *bindexec.Request[*debug.EchoArguments], res *bindexec.Response[*debug.EchoOK]) error {
				task := req.Task()
				echoDebugLog.Debugf("%+v", task.BindArguments())
				return res.SetSuccess(&debug.EchoOK{Message: task.BindArguments().Message})
			},
		),
	}
}
