package main

import (
	pdm "github.com/alanshaw/1up-service/pkg/store/provider/datamodel"
	cbg "github.com/whyrusleeping/cbor-gen"
)

func main() {
	if err := cbg.WriteMapEncodersToFile("../cbor_gen.go", "datamodel",
		pdm.ProviderModel{},
	); err != nil {
		panic(err)
	}
}
