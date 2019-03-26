package provisioner

import "github.com/mahakamcloud/mahakam/pkg/node"

// Provisioner defines sets of methods to be called to cloud provider
type Provisioner interface {
	CreateNode(nconfig node.NodeCreateConfig) error
}
