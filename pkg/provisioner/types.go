package provisioner

import "github.com/mahakamcloud/mahakam/pkg/resource_store/resource"

// Provisioner defines sets of methods to be called to cloud provider
type Provisioner interface {
	CreateNode(node resource.Node) error
}
