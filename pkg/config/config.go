package config

import (
	"fmt"
	"net"

	"github.com/mahakamcloud/mahakam/pkg/utils"
	yaml "gopkg.in/yaml.v2"
)

// Some config constants or environment variables will go away
// once we populate this through kind of config.yaml
const (
	// ResourceOwner hardcodes all tenant resources to be owned by gojek
	// since we don't have auth mechanism yet
	ResourceOwnerGojek   = "gojek"
	ResourceOwnerMahakam = "mahakam"

	// Custom path for storing resource in kvstore
	KeyPathMahakam       = "mahakamcloud/"
	KeyPathNetworkSubnet = KeyPathMahakam + "network/subnet/"
)

// Config represents mahakam configuration
type Config struct {
	KVStoreConfig StorageBackendConfig `yaml:"storage_backend"`
	NetworkConfig NetworkConfig        `yaml:"network"`
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

	if err := c.NetworkConfig.Validate(); err != nil {
		return fmt.Errorf("Error validating network configuration: %s", err)
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

// NetworkConfig stores metadata for network that mahakam will configure
type NetworkConfig struct {
	// Base CIDR that Mahakam will use to provision small network from it
	CIDR string `yaml:"cidr"`
}

// Validate validates storage backend configuration
func (nc *NetworkConfig) Validate() error {
	if nc.CIDR == "" {
		return fmt.Errorf("Must provide non-empty network CIDR")
	}

	if _, _, err := net.ParseCIDR(nc.CIDR); err != nil {
		return fmt.Errorf("Must provide valid CIDR format")
	}

	return nil
}
