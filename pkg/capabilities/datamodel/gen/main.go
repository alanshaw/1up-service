package main

import (
	cdm "github.com/alanshaw/1up-service/pkg/capabilities/datamodel"
	cbg "github.com/whyrusleeping/cbor-gen"
)

func main() {
	if err := cbg.WriteMapEncodersToFile("../cbor_gen.go", "datamodel",
		cdm.UnitModel{},
	); err != nil {
		panic(err)
	}
}
