package provisioner

import (
	"github.com/mahakamcloud/mahakam/pkg/node"
)

// Task defines task to be performed in a job
type Task interface {
	Run() error
}

// Provisioner defines sets of methods to be called to cloud provider
type Provisioner interface {
	CreateNode(config node.NodeCreateConfig) error
}
