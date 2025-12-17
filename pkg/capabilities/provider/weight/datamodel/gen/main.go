package main

import (
	wdm "github.com/alanshaw/1up-service/pkg/capabilities/provider/weight/datamodel"
	cbg "github.com/whyrusleeping/cbor-gen"
)

func main() {
	if err := cbg.WriteMapEncodersToFile("../cbor_gen.go", "datamodel",
		wdm.SetArgumentsModel{},
	); err != nil {
		panic(err)
	}
}
