package router

import "github.com/alanshaw/ucantone/ucan"

type SelectOption func(*selectConfig)

type selectConfig struct {
	exclusions []ucan.Principal
}

func WithExclusions(exclusions ...ucan.Principal) SelectOption {
	return func(cfg *selectConfig) {
		cfg.exclusions = exclusions
	}
}
