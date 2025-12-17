package main

import (
	hdm "github.com/alanshaw/1up-service/pkg/capabilities/http/datamodel"
	cbg "github.com/whyrusleeping/cbor-gen"
)

func main() {
	if err := cbg.WriteMapEncodersToFile("../cbor_gen.go", "datamodel",
		hdm.PutArgumentsModel{},
	); err != nil {
		panic(err)
	}
}
