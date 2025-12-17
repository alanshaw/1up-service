package app

type StorageConfig struct {
	DataDir string
	// Service-specific storage configurations
	Provider ProviderStorageConfig
}

type ProviderStorageConfig struct {
	Dir string
}
