package provisioner

import (
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/node"
	"github.com/mahakamcloud/mahakam/pkg/tfmodule"
	"github.com/mahakamcloud/mahakam/pkg/utils"
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

func (tp *terraformProvisioner) CreateNode(nconfig node.NodeCreateConfig) error {
	data := tp.getTerraformData(nconfig)
	log.Infof("terraform raw data to render files: %v\n", nconfig)
	log.Infof("terraform data to render files: %v\n", data)

	var err error
	switch role := nconfig.Role; role {
	case node.RoleControlPlane:
		err = tfmodule.CreateControlPlaneNode(nconfig.Name, config.TerraformDefaultDirectory+nconfig.Name, data)
	case node.RoleWorker:
		err = tfmodule.CreateWorkerNode(nconfig.Name, config.TerraformDefaultDirectory+nconfig.Name, data)
	case node.RoleNetworkDNS:
		err = tfmodule.CreateNetworkDNS(nconfig.Name, config.TerraformDefaultDirectory+nconfig.Name, data)
	case node.RoleNetworkDHCP:
		err = tfmodule.CreateNetworkDHCP(nconfig.Name, config.TerraformDefaultDirectory+nconfig.Name, data)
	case node.RoleNetworkGW:
		data = tp.overrideNetworkGWData(nconfig, data)
		log.Infof("terraform data for network gateway to render files: %v\n", data)

		err = tfmodule.CreateNetworkGWNode(nconfig.Name, config.TerraformDefaultDirectory+nconfig.Name, data)
	default:
		err = tfmodule.CreateNode(nconfig.Name, config.TerraformDefaultDirectory+nconfig.Name, data)
	}

	if err != nil {
		return err
	}
	return nil
}

func (tp *terraformProvisioner) getTerraformData(nconfig node.NodeCreateConfig) map[string]string {
	data := map[string]string{
		TerraformBucket:            config.TerraformDefaultBucket,
		TerraformKey:               nconfig.Name + "/terraform.tfstate",
		TerraformRegion:            config.TerraformDefaultRegion,
		TerraformName:              nconfig.Name,
		TerraformLibvirtModulePath: tp.config.LibvirtModulePath,
		TerraformHost:              nconfig.Host.String(),
		TerraformImageSourcePath:   tp.config.ImageSourcePath,
		TerraformSSHPublicKey:      nconfig.SSHPublicKey,
		TerraformMacAddress:        nconfig.MacAddress,
		TerraformIPAddress:         nconfig.IP.String(),
		TerraformNetmask:           utils.IPv4MaskString(nconfig.Node.Mask),
		TerraformGateway:           nconfig.Gateway.String(),
		TerraformDNSAddress:        nconfig.Nameserver.String(),
		TerraformDNSDomainName:     nconfig.Name + ".gocloud.io",
		TerraformControlPlaneIP:    nconfig.ExtraConfig[config.KeyControlPlaneIP],
		TerraformPodNetworkCidr:    nconfig.ExtraConfig[config.KeyPodNetworkCidr],
		TerraformKubeadmToken:      nconfig.ExtraConfig[config.KeyKubeadmToken],
	}
	return data
}

func (tp *terraformProvisioner) overrideNetworkGWData(nconfig node.NodeCreateConfig, data map[string]string) map[string]string {
	if len(nconfig.ExtraNetworks) == 0 {
		return data
	}

	data[TerraformPublicIPAddress] = nconfig.ExtraNetworks[0].IP.String()
	data[TerraformPublicNetmask] = utils.IPv4MaskString(nconfig.ExtraNetworks[0].Mask)
	data[TerraformGateway] = nconfig.ExtraNetworks[0].Gateway.String()
	data[TerraformDNSAddress] = nconfig.ExtraNetworks[0].Nameserver.String()
	data[TerraformNetworkCIDR] = nconfig.ExtraConfig[config.KeyClusterNetworkCidr]
	data[TerraformLibvirtModulePath] = tp.config.PublicLibvirtModulePath

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
