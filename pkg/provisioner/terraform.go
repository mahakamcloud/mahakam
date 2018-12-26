package provisioner

import (
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/node"
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
)

type terraformProvisioner struct {
	config config.TerraformConfig
}

func NewTerraformProvisioner(config config.TerraformConfig) Provisioner {
	return &terraformProvisioner{
		config: config,
	}
}

func (tp *terraformProvisioner) CreateNode(nconfig node.NodeCreateConfig) error {
	data := tp.getTerraformData(nconfig)
	log.Infof("terraform data to render files: %v\n", data)

	err := tfmodule.CreateNode(nconfig.Name, config.TerraformDefaultDirectory+nconfig.Name, data)
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
	}
	return data
}
