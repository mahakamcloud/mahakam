package config

// Some config constants or environment variables will go away
// once we populate this through kind of config.yaml
const (
	// ResourceOwner hardcodes all tenant resources to be owned by gojek
	// since we don't have auth mechanism yet
	ResourceOwnerGojek = "gojek"
)

// Config represents mahakam configuration
type Config struct {
	KVStoreConfig StorageBackendConfig
}

// StorageBackendConfig stores metadata for storage backend that we use
type StorageBackendConfig struct {
	BackendType string
	Address     string
	Username    string
	Password    string
	Bucket      string
}
