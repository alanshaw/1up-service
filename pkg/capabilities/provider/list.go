package provider

import (
	pdm "github.com/alanshaw/1up-service/pkg/capabilities/provider/datamodel"
	"github.com/alanshaw/ucantone/validator/bindcap"
)

const ListCommand = "/provider/list"

type (
	ListArguments = pdm.ListArgumentsModel
	ListOK        = pdm.ListOKModel
	Provider      = pdm.ProviderModel
)

var List, _ = bindcap.New[*ListArguments](ListCommand)
