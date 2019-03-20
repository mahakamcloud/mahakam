package provisioner

import (
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/resource_store/resource"
	"github.com/mahakamcloud/mahakam/pkg/tfmodule"
	log "github.com/sirupsen/logrus"
)

const (
	TerraformBucket            = "Bucket"
	TerraformKey               = "Key"
	TerraformRegion            = "Region"
	TerraformName              = "Name"
	TerraformLibvirtModulePath = "LibvirtModulePath"
	TerraformHost              = "Host"
	TerraformImageSourcePath   = "ImageSourcePath"
	TerraformSSHPublicKey      = "SSHPublicKey"
	TerraformMacAddress        = "MacAddress"
	TerraformIPAddress         = "IPAddress"
	TerraformNetmask           = "NetMask"
	TerraformGateway           = "Gateway"
	TerraformDNSAddress        = "DNSAddress"
	TerraformDNSDomainName     = "DNSDomainName"
	TerraformPublicIPAddress   = "PublicIPAddress"
	TerraformPublicNetmask     = "PublicNetmask"
	TerraformNetworkCIDR       = "NetworkCIDR"
	TerraformControlPlaneIP    = "ControlPlaneIP"
	TerraformPodNetworkCidr    = "PodNetworkCidr"
	TerraformKubeadmToken      = "KubeadmToken"
	TerraformBroadcastAddress  = "BroadcastAddress"
	TerraformSubnetAddress     = "SubnetAddress"
	TerraformSubnetMask        = "SubnetMask"
	TerraformMemory            = "MemorySize"
	TerraformCPU               = "NumCpu"
)

type terraformProvisioner struct {
	config config.TerraformConfig
}

// NewTerraformProvisioner returns a terraform provisioner based on passed config
func NewTerraformProvisioner(config config.TerraformConfig) Provisioner {
	return &terraformProvisioner{
		config: config,
	}
}

func (tp *terraformProvisioner) CreateNode(node resource.Node) error {
	data := tp.getTerraformData(node)

	log.Infof("terraform raw data to render files: %v\n", node)
	log.Infof("terraform data to render files: %v\n", data)

	var err error

	switch role := node.Role(); role {
	case string(resource.RoleK8sControlPlane):
		data = tp.overrideControlPlaneData(node, data)
		log.Infof("terraform data for control plane nodes to render files: %v\n", data)
		err = tfmodule.CreateControlPlaneNode(node.Name, config.TerraformDefaultDirectory+node.Name, data)
	case string(resource.RoleK8sWorker):
		data = tp.overrideWorkerData(node, data)
		log.Infof("terraform data for worker nodes to render files: %v\n", data)
		err = tfmodule.CreateWorkerNode(node.Name, config.TerraformDefaultDirectory+node.Name, data)
	case string(resource.RoleNetworkDNS):
		log.Infof("terraform data for network DNS to render files: %v\n", data)
		err = tfmodule.CreateNetworkDNS(node.Name, config.TerraformDefaultDirectory+node.Name, data)
	case string(resource.RoleNetworkDHCP):
		data := tp.overrideNetworkDHCPData(node, data)
		log.Infof("terraform data for network dhcp to render files: %v\n", data)
		err = tfmodule.CreateNetworkDHCP(node.Name, config.TerraformDefaultDirectory+node.Name, data)
	case string(resource.RoleNetworkGW):
		data = tp.overrideNetworkGWData(node, data)
		log.Infof("terraform data for network gateway to render files: %v\n", data)
		err = tfmodule.CreateNetworkGWNode(node.Name, config.TerraformDefaultDirectory+node.Name, data)
	}

	return err
}

func (tp *terraformProvisioner) getTerraformData(node resource.Node) map[string]string {
	networkConfigs := node.NetworkConfigs()
	extraConfig := node.ExtraConfig()

	data := map[string]string{
		TerraformBucket:            config.TerraformDefaultBucket,
		TerraformKey:               node.Name + config.TerraformTFState,
		TerraformRegion:            config.TerraformDefaultRegion,
		TerraformName:              node.Name,
		TerraformLibvirtModulePath: tp.config.LibvirtModulePath,
		TerraformHost:              node.Status().Host(),
		TerraformImageSourcePath:   tp.config.ImageSourcePath,
		TerraformSSHPublicKey:      node.Metadata().SSHKeys()[0],
		TerraformMacAddress:        networkConfigs[0].MacAddress(),
		TerraformIPAddress:         networkConfigs[0].IPAddress(),
		TerraformNetmask:           networkConfigs[0].IPMask(),
		TerraformGateway:           networkConfigs[0].Gateway(),
		TerraformDNSAddress:        networkConfigs[0].Nameserver(),
		TerraformDNSDomainName:     networkConfigs[0].FQDN(),
		TerraformControlPlaneIP:    extraConfig[config.KeyControlPlaneIP],
		TerraformPodNetworkCidr:    extraConfig[config.KeyPodNetworkCidr],
		TerraformKubeadmToken:      extraConfig[config.KeyKubeadmToken],
	}
	return data
}

func (tp *terraformProvisioner) overrideNetworkDHCPData(node resource.Node, data map[string]string) map[string]string {
	extraConfig := node.ExtraConfig()

	data[TerraformNetworkCIDR] = extraConfig[config.KeyClusterNetworkCidr]
	data[TerraformBroadcastAddress] = extraConfig[config.KeyBroadcastAddress]
	data[TerraformSubnetAddress] = extraConfig[config.KeySubnetAddress]
	data[TerraformSubnetMask] = extraConfig[config.KeySubnetMask]

	return data
}

func (tp *terraformProvisioner) overrideNetworkGWData(node resource.Node, data map[string]string) map[string]string {
	publicNetworkConfig := node.NetworkConfigs()[1]
	extraConfig := node.ExtraConfig()

	data[TerraformPublicIPAddress] = publicNetworkConfig.IPAddress()
	data[TerraformPublicNetmask] = publicNetworkConfig.IPMask()
	data[TerraformGateway] = publicNetworkConfig.Gateway()
	data[TerraformDNSAddress] = publicNetworkConfig.Nameserver()
	data[TerraformNetworkCIDR] = extraConfig[config.KeyClusterNetworkCidr]
	data[TerraformLibvirtModulePath] = tp.config.PublicLibvirtModulePath

	return data
}

func (tp *terraformProvisioner) overrideWorkerData(node resource.Node, data map[string]string) map[string]string {
	data[TerraformMemory] = node.Memory()
	data[TerraformCPU] = node.CPUs()

	return data
}

func (tp *terraformProvisioner) overrideControlPlaneData(node resource.Node, data map[string]string) map[string]string {
	data[TerraformMemory] = node.Memory()
	data[TerraformCPU] = node.CPUs()

	return data
}

func (tp *terraformProvisioner) mergeTerraformData(tfData ...map[string]string) map[string]string {
	var mergedData map[string]string
	for _, data := range tfData {
		for k, v := range data {
			mergedData[k] = v
		}
	}
	return mergedData
}
