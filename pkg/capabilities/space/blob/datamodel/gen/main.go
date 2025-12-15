package main

import (
	bdm "github.com/alanshaw/1up-service/pkg/capabilities/space/blob/datamodel"
	jsg "github.com/alanshaw/dag-json-gen"
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
	if err := jsg.WriteMapEncodersToFile("../dag_json_gen.go", "datamodel",
		bdm.AddArgumentsModel{},
		bdm.BlobModel{},
		bdm.AddOKModel{},
	); err != nil {
		panic(err)
	}
}
