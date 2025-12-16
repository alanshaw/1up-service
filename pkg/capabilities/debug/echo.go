package debug

import (
	ddm "github.com/alanshaw/1up-service/pkg/capabilities/debug/datamodel"
	"github.com/alanshaw/ucantone/ucan/delegation/policy"
	"github.com/alanshaw/ucantone/validator/bindcap"
	"github.com/alanshaw/ucantone/validator/capability"
)

const EchoCommand = "/debug/echo"

type (
	EchoArguments = ddm.EchoArgumentsModel
	EchoOK        = EchoArguments
)

var Echo, _ = bindcap.New[*EchoArguments](
	EchoCommand,
	capability.WithPolicyBuilder(policy.NotEqual(".message", "")),
)
