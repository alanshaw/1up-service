package provider

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"iter"
	"net/url"

	"github.com/alanshaw/1up-service/pkg/store"
	"github.com/alanshaw/libracha/capabilities/provider/datamodel"
	"github.com/alanshaw/ucantone/did"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
)

type DSProviderStore struct {
	ds datastore.Datastore
}

func NewDSProviderStore(dstore datastore.Datastore) *DSProviderStore {
	return &DSProviderStore{dstore}
}

func (d *DSProviderStore) Del(ctx context.Context, did did.DID) error {
	return d.ds.Delete(ctx, datastore.NewKey(did.String()))
}

func (d *DSProviderStore) Get(ctx context.Context, did did.DID) (datamodel.ProviderModel, error) {
	buf, err := d.ds.Get(ctx, datastore.NewKey(did.String()))
	if err != nil {
		if errors.Is(err, datastore.ErrNotFound) {
			return datamodel.ProviderModel{}, store.ErrNotFound
		}
		return datamodel.ProviderModel{}, err
	}
	var model datamodel.ProviderModel
	if err := model.UnmarshalCBOR(bytes.NewReader(buf)); err != nil {
		return datamodel.ProviderModel{}, err
	}
	return model, nil
}

func (d *DSProviderStore) List(ctx context.Context) iter.Seq2[datamodel.ProviderModel, error] {
	return func(yield func(datamodel.ProviderModel, error) bool) {
		results, err := d.ds.Query(ctx, query.Query{})
		if err != nil {
			yield(datamodel.ProviderModel{}, fmt.Errorf("querying datastore: %w", err))
			return
		}
		for entry := range results.Next() {
			if entry.Error != nil {
				yield(datamodel.ProviderModel{}, fmt.Errorf("iterating query results: %w", err))
				return
			}
			var model datamodel.ProviderModel
			if err := model.UnmarshalCBOR(bytes.NewReader(entry.Value)); err != nil {
				yield(datamodel.ProviderModel{}, fmt.Errorf("unmarshaling provider: %w", err))
				return
			}
			if !yield(model, nil) {
				return
			}
		}
	}
}

func (d *DSProviderStore) Put(ctx context.Context, id did.DID, endpoint *url.URL) error {
	model := datamodel.ProviderModel{
		Provider: id,
		Endpoint: endpoint.String(),
		Weight:   0,
	}
	var buf bytes.Buffer
	if err := model.MarshalCBOR(&buf); err != nil {
		return err
	}
	return d.ds.Put(ctx, datastore.NewKey(id.String()), buf.Bytes())
}

func (d *DSProviderStore) Update(ctx context.Context, did did.DID, update func(datamodel.ProviderModel) (datamodel.ProviderModel, error)) error {
	model, err := d.Get(ctx, did)
	if err != nil {
		if errors.Is(err, datastore.ErrNotFound) {
			return store.ErrNotFound
		}
		return err
	}
	updatedModel, err := update(model)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := updatedModel.MarshalCBOR(&buf); err != nil {
		return err
	}
	return d.ds.Put(ctx, datastore.NewKey(did.String()), buf.Bytes())
}

var _ Store = (*DSProviderStore)(nil)
