package datamodel

import "github.com/alanshaw/ucantone/did"

type ListArgumentsModel struct{}

type ProviderModel struct {
	Provider did.DID `cborgen:"provider"`
	Endpoint string  `cborgen:"endpoint"`
	Weight   uint64  `cborgen:"weight"`
}

type ListOKModel struct {
	Providers []ProviderModel `cborgen:"providers"`
}
