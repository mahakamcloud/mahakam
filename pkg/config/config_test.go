package config_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	. "github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

var (
	validConfig = Config{
		KVStoreConfig: StorageBackendConfig{
			BackendType: "fake-storage-backend-type",
			Address:     "fake-storage-backend-address",
			Username:    "fake-storage-backend-username",
			Password:    "fake-storage-backend-password",
			Bucket:      "fake-storage-backend-bucket",
		},
		NetworkConfig: NetworkConfig{
			CIDR:                  "1.2.3.4/16",
			ClusterNetmask:        24,
			DatacenterGatewayCIDR: "1.2.3.4/28",
			DatacenterNameserver:  "1.2.3.4",
		},
		KubernetesConfig: KubernetesConfig{
			PodNetworkCidr: "1.2.3.4/16",
			KubeadmToken:   "fake-kubeadm-token",
			SSHPublicKey:   "fake-ssh-public-key",
		},
		TerraformConfig: TerraformConfig{
			LibvirtModulePath:       "fake-libvirt-module-path",
			PublicLibvirtModulePath: "fake-public-libvirt-module-path",
			ImageSourcePath:         "fake-image-source-path",
		},
		GateConfig: GateConfig{
			GateNSSApiKey: "fake-gate-nss-api-key",
		},
		HostsConfig: HostsConfig{
			Hosts: []Host{
				Host{Name: "server01", IPAddress: "10.30.0.1"},
				Host{Name: "server02", IPAddress: "10.30.02"}},
		},
	}
)

func TestLoadConfig(t *testing.T) {
	var (
		err    error
		input  Config
		output *Config

		dir             string
		configFilePath  string
		configFileBytes []byte
	)

	// Setup config file directory
	input = validConfig
	dir, err = ioutil.TempDir("", "mahakam-config-")
	assert.Nil(t, err)
	configFilePath = filepath.Join(dir, "config-file")

	// Setup config file data
	configFileBytes, err = yaml.Marshal(input)
	assert.Nil(t, err)
	err = ioutil.WriteFile(configFilePath, configFileBytes, 0644)
	assert.Nil(t, err)

	output, err = LoadConfig(configFilePath)

	assert.Equal(t, *output, input, "they should be equal")
}
