package main

import (
	pdm "github.com/alanshaw/1up-service/pkg/capabilities/provider/datamodel"
	cbg "github.com/whyrusleeping/cbor-gen"
)

func main() {
	if err := cbg.WriteMapEncodersToFile("../cbor_gen.go", "datamodel",
		pdm.ListArgumentsModel{},
		pdm.ProviderModel{},
		pdm.ListOKModel{},
		pdm.RegisterArgumentsModel{},
		pdm.DeregisterArgumentsModel{},
	); err != nil {
		panic(err)
	}
}
