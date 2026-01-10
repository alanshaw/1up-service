package store

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alanshaw/1up-service/pkg/store/provider"
	"github.com/alanshaw/1up-service/pkg/store/token"
	leveldb "github.com/ipfs/go-ds-leveldb"
	"go.uber.org/fx"

	"github.com/alanshaw/1up-service/pkg/config/app"
)

var Module = fx.Module("store",
	fx.Provide(
		ProvideConfigs,
		NewProviderStore,
		NewTokenStore,
	),
)

type Configs struct {
	fx.Out
	Provider app.ProviderStorageConfig
	Token    app.TokenStorageConfig
}

// ProvideConfigs provides the fields of a storage config
func ProvideConfigs(cfg app.StorageConfig) Configs {
	return Configs{
		Provider: cfg.Provider,
		Token:    cfg.Token,
	}
}

func NewProviderStore(cfg app.ProviderStorageConfig, lc fx.Lifecycle) (provider.Store, error) {
	if cfg.Dir == "" {
		return nil, fmt.Errorf("no data dir provided for provider store")
	}

	ds, err := newDatastore(cfg.Dir)
	if err != nil {
		return nil, fmt.Errorf("creating provider store: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return ds.Close()
		},
	})

	return provider.NewDSProviderStore(ds), nil
}

func NewTokenStore(cfg app.TokenStorageConfig, lc fx.Lifecycle) (token.Store, error) {
	if cfg.Dir == "" {
		return nil, fmt.Errorf("no data dir provided for token store")
	}

	ds, err := newDatastore(cfg.Dir)
	if err != nil {
		return nil, fmt.Errorf("creating token store: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return ds.Close()
		},
	})

	return token.NewDSTokenStore(ds), nil
}

func newDatastore(path string) (*leveldb.Datastore, error) {
	dirPath, err := mkdirp(path)
	if err != nil {
		return nil, fmt.Errorf("creating leveldb for store at path %s: %w", path, err)
	}
	return leveldb.NewDatastore(dirPath, nil)
}

func mkdirp(dirpath ...string) (string, error) {
	dir := filepath.Join(dirpath...)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", fmt.Errorf("creating directory: %s: %w", dir, err)
	}
	return dir, nil
}
