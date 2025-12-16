package main

import (
	ddm "github.com/alanshaw/1up-service/pkg/capabilities/debug/datamodel"
	jsg "github.com/alanshaw/dag-json-gen"
	cbg "github.com/whyrusleeping/cbor-gen"
)

func main() {
	if err := cbg.WriteMapEncodersToFile("../cbor_gen.go", "datamodel",
		ddm.EchoArgumentsModel{},
	); err != nil {
		panic(err)
	}
	if err := jsg.WriteMapEncodersToFile("../dag_json_gen.go", "datamodel",
		ddm.EchoArgumentsModel{},
	); err != nil {
		panic(err)
	}
}
