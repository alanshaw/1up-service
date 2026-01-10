package app

type StorageConfig struct {
	DataDir string
	// Service-specific storage configurations
	Provider ProviderStorageConfig
	Token    TokenStorageConfig
}

type ProviderStorageConfig struct {
	Dir string
}

type TokenStorageConfig struct {
	Dir string
}
