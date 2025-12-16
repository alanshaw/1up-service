package provider

import (
	"context"
	"iter"

	"github.com/alanshaw/1up-service/pkg/store/provider/datamodel"
	"github.com/alanshaw/ucantone/did"
)

type Store interface {
	Put(ctx context.Context, model datamodel.ProviderModel) error
	Del(ctx context.Context, did did.DID) error
	Get(ctx context.Context, did did.DID) (datamodel.ProviderModel, error)
	List(ctx context.Context) iter.Seq2[datamodel.ProviderModel, error]
}
