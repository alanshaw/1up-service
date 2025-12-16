package datamodel

import (
	"github.com/alanshaw/ucantone/did"
	"github.com/ipfs/go-cid"
)

type RegisterArgumentsModel struct {
	Provider did.DID `cborgen:"provider"`
	Endpoint string  `cborgen:"endpoint"`
	Proof    cid.Cid `cborgen:"proof"`
}

type RegisterOKModel struct{}
