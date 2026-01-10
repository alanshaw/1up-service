package routing

import "github.com/alanshaw/ucantone/ucan"

type SelectOption func(*selectConfig)

type selectConfig struct {
	exclusions []ucan.Principal
}

// WithExclusions specifies a list of providers to exclude from selection.
func WithExclusions(exclusions ...ucan.Principal) SelectOption {
	return func(cfg *selectConfig) {
		cfg.exclusions = exclusions
	}
}
