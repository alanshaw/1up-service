package token

import (
	"context"
	"iter"

	"github.com/alanshaw/ucantone/ucan"
	"github.com/ipfs/go-cid"
)

// Store is a store for UCAN tokens.
type Store interface {
	// Put a container of tokens to the store.
	Put(ctx context.Context, container ucan.Container) error
	// FindByTask returns all UCAN containers that have tokens related to the
	// given task i.e. invocations or receipts for the task.
	FindByTask(ctx context.Context, task cid.Cid) iter.Seq2[ucan.Container, error]
}
