package config

import (
	"fmt"

	"github.com/mahakamcloud/mahakam/pkg/utils"
	yaml "gopkg.in/yaml.v2"
)

// Some config constants or environment variables will go away
// once we populate this through kind of config.yaml
const (
	// ResourceOwner hardcodes all tenant resources to be owned by gojek
	// since we don't have auth mechanism yet
	ResourceOwnerGojek = "gojek"
)

// Config represents mahakam configuration
type Config struct {
	KVStoreConfig StorageBackendConfig `yaml:"storage_backend"`
}

// LoadConfig loads a configuration file
func LoadConfig(configFilePath string) (*Config, error) {
	var config *Config
	if configFilePath == "" {
		return config, fmt.Errorf("Must provide non-empty configuration file path")
	}

	bytes, err := utils.ReadFile(configFilePath)
	if err != nil {
		return config, err
	}

	if err = yaml.Unmarshal(bytes, &config); err != nil {
		return config, fmt.Errorf("Error unmarshaling configuration file: %s", err)
	}

	if err = config.Validate(); err != nil {
		return config, fmt.Errorf("Error validating configuration file: %s", err)
	}

	return config, nil
}

// Validate validates mahakam configuration
func (c *Config) Validate() error {
	if err := c.KVStoreConfig.Validate(); err != nil {
		return fmt.Errorf("Error validating KV store configuration: %s", err)
	}
	return nil
}

// StorageBackendConfig stores metadata for storage backend that we use
type StorageBackendConfig struct {
	BackendType string `yaml:"backend_type"`
	Address     string `yaml:"address"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Bucket      string `yaml:"bucket"`
}

// Validate validates storage backend configuration
func (sbc *StorageBackendConfig) Validate() error {
	if sbc.Address == "" {
		return fmt.Errorf("Must provide non-empty storage backend address")
	}

	return nil
}
