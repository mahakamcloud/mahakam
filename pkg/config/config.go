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
	KeyPathNetworkSubnet = KeyPathMahakam + "network/subnets/"

	// Helm default configuration
	HelmDefaultNamespace             = "default"
	HelmDefaultTillerNamespace       = "kube-system"
	HelmDefaultKubecontext           = "kubernetes-admin@kubernetes"
	HelmControllerWait               = false
	HelmControllerDefaultWaitTimeout = 300

	// Default mahakam config path to store multiple kubeconfig files
	MahakamMultiKubeconfigPath = "/opt/mahakamcloud/clusters"

	// Default terraform config
	TerraformDefaultDirectory = "/opt/mahakamcloud/terraform/"
	TerraformDefaultBucket    = "tf-mahakam"
	TerraformDefaultRegion    = "ap-southeast-1"

	// Various keys for storing metadata in map
	KeyControlPlaneIP = "key-control-plane-ip"
)

// Config represents mahakam configuration
type Config struct {
	KVStoreConfig   StorageBackendConfig `yaml:"storage_backend"`
	NetworkConfig   NetworkConfig        `yaml:"network"`
	GateConfig      GateConfig           `yaml:"gate"`
	TerraformConfig TerraformConfig      `yaml:"terraform"`
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

	if err := c.GateConfig.Validate(); err != nil {
		return fmt.Errorf("Error validating gate configuration: %s", err)
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
	// CIDR is datacenter network CIDR that Mahakam will use to provision cluster network from it
	CIDR string `yaml:"cidr"`
	// ClusterNetmask is subnet length that cluster network will be provisioned as
	ClusterNetmask int `yaml:"cluster_netmask"`
}

// Validate validates storage backend configuration
func (nc *NetworkConfig) Validate() error {
	if nc.CIDR == "" {
		return fmt.Errorf("Must provide non-empty network CIDR")
	}

	if nc.ClusterNetmask == 0 {
		return fmt.Errorf("Must provide non-empty cluster netmask")
	}

	if _, _, err := net.ParseCIDR(nc.CIDR); err != nil {
		return fmt.Errorf("Must provide valid CIDR format")
	}

	if nc.ClusterNetmask > 32 || nc.ClusterNetmask < 1 {
		return fmt.Errorf("Must provide valid cluster netmask between 0 and 32")
	}

	return nil
}

type GateConfig struct {
	GateNSSApiKey string `yaml:"gate_nss_api_key"`
}

func (gc *GateConfig) Validate() error {
	if gc.GateNSSApiKey == "" {
		return fmt.Errorf("Must provide non-empty Gate NSS API key")
	}
	return nil
}

type TerraformConfig struct {
	LibvirtModulePath string `yaml:"libvirt_module_path"`
	ImageSourcePath   string `yaml:"image_source_path"`
}

func (tc *TerraformConfig) Validate() error {
	if tc.LibvirtModulePath == "" {
		return fmt.Errorf("Must provide non-empty libvirt module path")
	}
	if tc.ImageSourcePath == "" {
		return fmt.Errorf("Must provide non-empty image source path")
	}
	return nil
}
