package routing

import (
	"context"
	"errors"
	"math/rand"
	"net/url"
	"slices"

	"github.com/alanshaw/1up-service/pkg/store"
	"github.com/alanshaw/1up-service/pkg/store/provider"
	"github.com/alanshaw/ucantone/ucan"
	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("service/router")

type ProviderInfo struct {
	ID       ucan.Principal
	Endpoint *url.URL
}

type Router struct {
	providerStore provider.Store
}

func NewRouter(providerStore provider.Store) *Router {
	return &Router{providerStore}
}

// Provider returns information about a registered storage provider.
func (r *Router) Provider(ctx context.Context, id ucan.Principal) (ProviderInfo, error) {
	p, err := r.providerStore.Get(ctx, id.DID())
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return ProviderInfo{}, ErrNotFound
		}
		return ProviderInfo{}, err
	}
	endpoint, err := url.Parse(p.Endpoint)
	if err != nil {
		return ProviderInfo{}, err
	}
	return ProviderInfo{
		ID:       p.Provider,
		Endpoint: endpoint,
	}, nil
}

// Select chooses a registered storage provider based on provider weight.
func (r *Router) Select(ctx context.Context, options ...SelectOption) (ProviderInfo, error) {
	cfg := &selectConfig{}
	for _, opt := range options {
		opt(cfg)
	}

	var candidates []ProviderInfo
	var weights []uint64
	for p, err := range r.providerStore.List(ctx) {
		if err != nil {
			return ProviderInfo{}, err
		}
		if slices.ContainsFunc(cfg.exclusions, func(e ucan.Principal) bool {
			return e.DID() == p.Provider
		}) {
			continue
		}
		endpoint, err := url.Parse(p.Endpoint)
		if err != nil {
			log.Warnf("provider %q has invalid endpoint: %w", p.Provider.DID(), err)
			continue
		}
		if p.Weight == 0 {
			continue
		}
		candidates = append(candidates, ProviderInfo{
			ID:       p.Provider,
			Endpoint: endpoint,
		})
		weights = append(weights, p.Weight)
	}
	switch len(candidates) {
	case 0:
		return ProviderInfo{}, ErrCandidateUnavailable
	case 1:
		return candidates[0], nil
	}
	return candidates[getWeightedRandomInt(weights)], nil
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
