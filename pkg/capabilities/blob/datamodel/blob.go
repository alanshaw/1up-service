package datamodel

import (
	"github.com/alanshaw/1up-service/pkg/capabilities"
	"github.com/alanshaw/ucantone/ucan"
	"github.com/alanshaw/ucantone/ucan/promise"
	"github.com/multiformats/go-multihash"
)

type AllocateArgumentsModel struct {
	Blob  BlobModel `cborgen:"blob"`
	Cause ucan.Link `cborgen:"cause"`
}

type BlobModel struct {
	Digest multihash.Multihash `cborgen:"digest"`
	Size   uint64              `cborgen:"size"`
}

type AllocateOKModel struct {
	Size    uint64            `cborgen:"size"`
	Address *BlobAddressModel `cborgen:"address,omitempty"`
}

type BlobAddressModel struct {
	URL     capabilities.CborURL  `cborgen:"url"`
	Headers map[string]string     `cborgen:"headers"`
	Expires capabilities.CborTime `cborgen:"expires"`
}

type AcceptArgumentsModel struct {
	Blob BlobModel       `cborgen:"blob"`
	Put  promise.AwaitOK `cborgen:"_put"`
}

type AcceptOKModel struct {
	Site ucan.Link `cborgen:"site"`
}
