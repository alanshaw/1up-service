package provider

import (
	cdm "github.com/alanshaw/1up-service/pkg/capabilities/datamodel"
	pdm "github.com/alanshaw/1up-service/pkg/capabilities/provider/datamodel"
	"github.com/alanshaw/ucantone/validator/bindcap"
)

const DeregisterCommand = "/provider/deregister"

type (
	DeregisterArguments = pdm.DeregisterArgumentsModel
	DeregisterOK        = cdm.UnitModel
)

var Deregister, _ = bindcap.New[*DeregisterArguments](DeregisterCommand)
