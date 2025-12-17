package router

import (
	"context"
	"math/rand"

	"github.com/alanshaw/1up-service/pkg/store/provider"
	"github.com/alanshaw/ucantone/ucan"
	"github.com/multiformats/go-multihash"
)

type Router struct {
	providerStore provider.Store
}

func NewRouter(providerStore provider.Store) *Router {
	return &Router{}
}

// Select chooses a registered storage provider based on provider weight.
func (r *Router) Select(ctx context.Context, digest multihash.Multihash, size uint64) (ucan.Principal, error) {
	ids := []ucan.Principal{}
	weights := []uint64{}
	for p, err := range r.providerStore.List(ctx) {
		if err != nil {
			return nil, err
		}
		ids = append(ids, p.Provider)
		weights = append(weights, p.Weight)
	}
	if len(ids) == 0 {
		return nil, NewCandidateUnavailableError("no storage providers available")
	}
	return ids[getWeightedRandomInt(weights)], nil
}

func getWeightedRandomInt(weights []uint64) int {
	totalWeight := uint64(0)
	for _, weight := range weights {
		totalWeight += weight
	}
	random := uint64(rand.Int63n(int64(totalWeight)))
	for i, weight := range weights {
		random -= weight
		if random <= 0 {
			return i
		}
	}
	panic("did not find a weight - should never reach here")
}
