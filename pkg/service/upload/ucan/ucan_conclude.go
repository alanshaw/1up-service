package ucan

import (
	"github.com/alanshaw/1up-service/pkg/service"
	"github.com/alanshaw/1up-service/pkg/service/router"
	ucan_caps "github.com/alanshaw/libracha/capabilities/ucan"
	"github.com/alanshaw/ucantone/execution/bindexec"
	"github.com/alanshaw/ucantone/principal"
	logging "github.com/ipfs/go-log/v2"
)

var ucanConcludeLog = logging.Logger("service/upload/ucan/ucan/conclude")

func NewUCANConcludeHandler(id principal.Signer, router *router.Router) *service.Handler {
	return &service.Handler{
		Capability: ucan_caps.Conclude,
		Handler: bindexec.NewHandler(
			func(req *bindexec.Request[*ucan_caps.ConcludeArguments]) (*bindexec.Response[*ucan_caps.ConcludeOK], error) {
				// var accDlg ucan.Delegation
				// for _, d := range req.Metadata().Delegations() {
				// 	if d.Audience().DID() == id.DID() && d.Command() == blob_caps.AcceptCommand {
				// 		accDlg = d
				// 		break
				// 	}
				// }
				// if accDlg == nil {
				// 	return nil, fmt.Errorf("missing blob accept delegation")
				// }
				// // TODO: use the correct provider
				// prov, err := router.Select(req.Context())
				// if err != nil {
				// 	return nil, err
				// }
				// blob_caps.Accept.Invoke(
				// 	id,
				// 	accDlg.Subject(),
				// 	&blob_caps.AcceptArguments{},
				// 	ucan.WithProofs(accDlg.Link()),
				// )
				return bindexec.NewResponse(bindexec.WithSuccess(&ucan_caps.ConcludeOK{}))
			},
		),
	}
}
