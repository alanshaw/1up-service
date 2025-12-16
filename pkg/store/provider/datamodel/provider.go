package datamodel

import (
	"github.com/alanshaw/ucantone/did"
	"github.com/alanshaw/ucantone/ucan/delegation"
)

type ProviderModel struct {
	Provider did.DID                `cborgen:"provider"`
	Endpoint string                 `cborgen:"endpoint"`
	Proof    *delegation.Delegation `cborgen:"proof"`
	Weight   uint64                 `cborgen:"weight"`
}
