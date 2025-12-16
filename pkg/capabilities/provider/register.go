package provider

import (
	pdm "github.com/alanshaw/1up-service/pkg/capabilities/provider/datamodel"
	"github.com/alanshaw/ucantone/validator/bindcap"
)

const RegisterCommand = "/provider/register"

type (
	RegisterArguments = pdm.RegisterArgumentsModel
	RegisterOK        = pdm.RegisterOKModel
)

var Register, _ = bindcap.New[*RegisterArguments](RegisterCommand)
