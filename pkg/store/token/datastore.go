package token

import (
	"bytes"
	"context"
	"iter"

	"github.com/alanshaw/ucantone/ipld/codec/dagcbor"
	"github.com/alanshaw/ucantone/ucan"
	"github.com/alanshaw/ucantone/ucan/container"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	"github.com/ipfs/go-datastore/query"
	"github.com/multiformats/go-multihash"
)

type DSTokenStore struct {
	index datastore.Datastore
	data  datastore.Datastore
}

func (d *DSTokenStore) FindByTask(ctx context.Context, task cid.Cid) iter.Seq2[ucan.Container, error] {
	return func(yield func(ucan.Container, error) bool) {
		key := datastore.NewKey(task.String())
		results, err := d.index.Query(ctx, query.Query{Prefix: key.String()})
		if err != nil {
			yield(nil, err)
			return
		}
		for entry := range results.Next() {
			if entry.Error != nil {
				yield(nil, err)
				return
			}
			_, ctLink, err := cid.CidFromBytes(entry.Value)
			if err != nil {
				yield(nil, err)
				return
			}
			ctBytes, err := d.data.Get(ctx, datastore.NewKey(ctLink.String()))
			if err != nil {
				yield(nil, err)
				return
			}
			var ct container.Container
			err = ct.UnmarshalCBOR(bytes.NewReader(ctBytes))
			if err != nil {
				yield(nil, err)
				return
			}
			if !yield(&ct, nil) {
				return
			}
		}
	}
}

func (d *DSTokenStore) Put(ctx context.Context, tokens ucan.Container) error {
	ct := container.New(
		container.WithInvocations(tokens.Invocations()...),
		container.WithDelegations(tokens.Delegations()...),
		container.WithReceipts(tokens.Receipts()...),
	)

	var ctBuf bytes.Buffer
	err := ct.MarshalCBOR(&ctBuf)
	if err != nil {
		return err
	}
	ctLink, err := cid.V1Builder{
		Codec:  dagcbor.Code,
		MhType: multihash.SHA2_256,
	}.Sum(ctBuf.Bytes())
	if err != nil {
		return err
	}
	err = d.data.Put(ctx, datastore.NewKey(ctLink.String()), ctBuf.Bytes())
	if err != nil {
		return err
	}

	// index the data by <task-cid>/<container-cid>/<invocation-cid> and <task-cid>/<container-cid>/<receipt-cid>
	for _, inv := range ct.Invocations() {
		key := datastore.NewKey(inv.Task().Link().String()).ChildString(ctLink.String()).ChildString(inv.Link().String())
		err := d.index.Put(ctx, key, ctLink.Bytes())
		if err != nil {
			return err
		}
	}
	for _, rcpt := range ct.Receipts() {
		key := datastore.NewKey(rcpt.Ran().String()).ChildString(ctLink.String()).ChildString(rcpt.Link().String())
		err := d.index.Put(ctx, key, ctLink.Bytes())
		if err != nil {
			return err
		}
	}
	return nil
}

var _ Store = (*DSTokenStore)(nil)

func NewDSTokenStore(ds datastore.Datastore) *DSTokenStore {
	return &DSTokenStore{
		index: namespace.Wrap(ds, datastore.NewKey("/index")),
		data:  namespace.Wrap(ds, datastore.NewKey("/data")),
	}
}
