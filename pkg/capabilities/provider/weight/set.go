package weight

import (
	cdm "github.com/alanshaw/1up-service/pkg/capabilities/datamodel"
	wdm "github.com/alanshaw/1up-service/pkg/capabilities/provider/weight/datamodel"
	"github.com/alanshaw/ucantone/ucan/delegation/policy"
	"github.com/alanshaw/ucantone/validator/bindcap"
	"github.com/alanshaw/ucantone/validator/capability"
)

const SetCommand = "/provider/weight/set"

type (
	SetArguments = wdm.SetArgumentsModel
	SetOK        = cdm.UnitModel
)

var Set, _ = bindcap.New[*SetArguments](
	SetCommand,
	capability.WithPolicyBuilder(
		policy.GreaterThanOrEqual(".weight", 0),
	),
)
