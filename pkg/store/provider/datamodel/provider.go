package datamodel

import (
	"github.com/alanshaw/ucantone/did"
)

type ProviderModel struct {
	Provider did.DID `cborgen:"provider"`
	Endpoint string  `cborgen:"endpoint"`
	Weight   uint64  `cborgen:"weight"`
}
