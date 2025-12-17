package provider

import (
	"context"
	"iter"
	"net/url"

	"github.com/alanshaw/1up-service/pkg/capabilities/provider/datamodel"
	"github.com/alanshaw/ucantone/did"
)

type Store interface {
	Put(ctx context.Context, id did.DID, endpoint *url.URL) error
	Del(ctx context.Context, id did.DID) error
	Get(ctx context.Context, id did.DID) (datamodel.ProviderModel, error)
	List(ctx context.Context) iter.Seq2[datamodel.ProviderModel, error]
}
