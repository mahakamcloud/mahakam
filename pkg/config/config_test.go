package config_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v2"

	. "github.com/mahakamcloud/mahakam/pkg/config"
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
			CIDR:           "1.2.3.4/16",
			ClusterNetmask: 24,
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
	}
)

var _ = Describe("LoadConfig", func() {
	var (
		err    error
		input  Config
		output *Config

		dir             string
		configFilePath  string
		configFileBytes []byte
	)

	BeforeEach(func() {
		input = validConfig
		dir, err = ioutil.TempDir("", "mahakam-config-")
		Expect(err).ToNot(HaveOccurred())
		configFilePath = filepath.Join(dir, "config-file")
	})

	JustBeforeEach(func() {
		output, err = LoadConfig(configFilePath)
	})

	AfterEach(func() {
		os.RemoveAll(dir)
	})

	Context("when having valid config file bytes", func() {
		BeforeEach(func() {
			configFileBytes, err = yaml.Marshal(input)
			Expect(err).ToNot(HaveOccurred())
			err = ioutil.WriteFile(configFilePath, configFileBytes, 0644)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns config file bytes", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(*output).To(Equal(input))
		})
	})
})
