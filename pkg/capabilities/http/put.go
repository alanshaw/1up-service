package http

import (
	cdm "github.com/alanshaw/1up-service/pkg/capabilities/datamodel"
	hdm "github.com/alanshaw/1up-service/pkg/capabilities/http/datamodel"
	"github.com/alanshaw/ucantone/validator/bindcap"
)

const PutCommand = "/http/put"

type (
	PutArguments = hdm.PutArgumentsModel
	PutOK        = cdm.UnitModel
)

var Put, _ = bindcap.New[*PutArguments](PutCommand)
