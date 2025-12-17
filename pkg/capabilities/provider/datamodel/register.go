package datamodel

import (
	"github.com/alanshaw/ucantone/did"
)

type RegisterArgumentsModel struct {
	Provider did.DID `cborgen:"provider"`
	Endpoint string  `cborgen:"endpoint"`
}

type RegisterOKModel struct{}
