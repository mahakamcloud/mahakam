package provisioner

import (
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/node"
	"github.com/mahakamcloud/mahakam/pkg/tfmodule"
)

// Task defines task to be performed in a job
type Task interface {
	Run() error
}

// Provisioner defines sets of methods to be called to cloud provider
type Provisioner interface {
	CreateNode(config node.NodeCreateConfig) error
}

type terraformProvisioner struct{}

func NewTerraformProvisioner() Provisioner {
	return &terraformProvisioner{}
}

func (tp *terraformProvisioner) CreateNode(nconfig node.NodeCreateConfig) error {
	err := tfmodule.CreateNode(nconfig.Name, config.TerraformDefaultDirectory+nconfig.Name, nconfig.ExtraConfig)
	if err != nil {
		return err
	}
	return nil
}
