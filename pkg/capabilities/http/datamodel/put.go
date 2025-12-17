package datamodel

import (
	"github.com/alanshaw/1up-service/pkg/capabilities/space/blob"
	"github.com/alanshaw/ucantone/ucan/promise"
)

type PutArgumentsModel struct {
	Body blob.Blob `cborgen:"body"`
	// Destination is the promise that resolves to the upload destination
	// where the blob should be PUT to. It is the result of a /blob/allocate task.
	Destination promise.AwaitOK `cborgen:"destination"`
}
