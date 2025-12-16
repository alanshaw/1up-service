package main

import (
	bdm "github.com/alanshaw/1up-service/pkg/capabilities/space/blob/datamodel"
	cbg "github.com/whyrusleeping/cbor-gen"
)

func main() {
	if err := cbg.WriteMapEncodersToFile("../cbor_gen.go", "datamodel",
		bdm.AddArgumentsModel{},
		bdm.BlobModel{},
		bdm.AddOKModel{},
	); err != nil {
		panic(err)
	}
}
