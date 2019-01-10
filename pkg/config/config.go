package config

import (
	"fmt"
	"net"

	"github.com/mahakamcloud/mahakam/pkg/utils"
	yaml "gopkg.in/yaml.v2"
)

// Config represents mahakam configuration
type Config struct {
	KVStoreConfig    StorageBackendConfig `yaml:"storage_backend"`
	NetworkConfig    NetworkConfig        `yaml:"network"`
	KubernetesConfig KubernetesConfig     `yaml:"kubernetes"`
	GateConfig       GateConfig           `yaml:"gate"`
	TerraformConfig  TerraformConfig      `yaml:"terraform"`
	HostsConfig      HostsConfig          `yaml:"hosts"`
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
	fmt.Println(config)
	if err = yaml.Unmarshal(bytes, &config); err != nil {
		return config, fmt.Errorf("Error unmarshaling configuration file: %s", err)
	}
	fmt.Println(config)
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

	if err := c.TerraformConfig.Validate(); err != nil {
		return fmt.Errorf("Error validating terraform configuration: %s", err)
	}

	if err := c.GateConfig.Validate(); err != nil {
		return fmt.Errorf("Error validating gate configuration: %s", err)
	}

	if err := c.TerraformConfig.Validate(); err != nil {
		return fmt.Errorf("Error validating terraform configuration: %s", err)
	}

	if err := c.HostsConfig.Validate(); err != nil {
		return fmt.Errorf("Error validating hosts configuration: %s", err)
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
	// DatacenterGatewayCIDR is gateway public IP in the datacenter that can reach Internet
	DatacenterGatewayCIDR string `yaml:"datacenter_gateway_cidr"`
	// DatacenterNameserver is nameserver in datacenter that can solve domains in Internet
	DatacenterNameserver string `yaml:"datacenter_nameserver"`
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
		return fmt.Errorf("Must provide valid CIDR format for cluster network")
	}

	if nc.ClusterNetmask > 32 || nc.ClusterNetmask < 1 {
		return fmt.Errorf("Must provide valid cluster netmask between 0 and 32")
	}

	if nc.DatacenterGatewayCIDR == "" {
		return fmt.Errorf("Must provide non-empty datacenter gateway CIDR")
	}

	if _, _, err := net.ParseCIDR(nc.DatacenterGatewayCIDR); err != nil {
		return fmt.Errorf("Must provide valid CIDR format for datacenter gateway")
	}

	if nc.DatacenterNameserver == "" {
		return fmt.Errorf("Must provide non-empty datacenter nameserver")
	}

	if validIP := net.ParseIP(nc.DatacenterNameserver); validIP == nil {
		return fmt.Errorf("Must provide valid IP format for datacenter nameserver")
	}

	return nil
}

// KubernetesConfig stores metadata for kubernetes cluster that mahakam will configure
type KubernetesConfig struct {
	// PodNetworkCidr is pod network that CNI will provision inside Kubernetes cluster
	PodNetworkCidr string `yaml:"pod_network_cidr"`
	// KubeadmToken is token secret used for workers to join Kubernetes cluster
	KubeadmToken string `yaml:"kubeadm_token"`
	// SSHPublicKey is public key that will be authorized to access Kubernetes nodes from its private key pair
	SSHPublicKey string `yaml:"ssh_public_key"`
}

// Validate validates storage backend configuration
func (kc *KubernetesConfig) Validate() error {
	if kc.PodNetworkCidr == "" {
		return fmt.Errorf("Must provide non-empty pod network CIDR")
	}

	if kc.KubeadmToken == "" {
		return fmt.Errorf("Must provide non-empty kubeadm token")
	}

	if kc.SSHPublicKey == "" {
		return fmt.Errorf("Must provide non-empty SSH public key")
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
	LibvirtModulePath       string `yaml:"libvirt_module_path"`
	PublicLibvirtModulePath string `yaml:"public_libvirt_module_path"`
	ImageSourcePath         string `yaml:"image_source_path"`
}

func (tc *TerraformConfig) Validate() error {
	if tc.LibvirtModulePath == "" {
		return fmt.Errorf("Must provide non-empty libvirt module path")
	}
	if tc.PublicLibvirtModulePath == "" {
		return fmt.Errorf("Must provide non-empty public libvirt module path")
	}
	if tc.ImageSourcePath == "" {
		return fmt.Errorf("Must provide non-empty image source path")
	}
	return nil
}

type Host struct {
	Name      string `yaml:"name"`
	IPAddress string `yaml:"ip_address"`
}

type HostsConfig struct {
	Hosts []Host `yaml:"hosts"`
}

func (hc *HostsConfig) Validate() error {
	if len(hc.Hosts) == 0 {
		return fmt.Errorf("Must provide atleast one host")
	}

	for _, host := range hc.Hosts {
		if host.Name == "" {
			return fmt.Errorf("Must provide non-empty host name for Host: [%s]", host)
		}

		if validIP := net.ParseIP(host.IPAddress); validIP == nil {
			return fmt.Errorf("Must provide valid IP format for Host: [%s]", host)
		}
	}
	return nil
}
