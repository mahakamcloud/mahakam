package resource

import "github.com/mahakamcloud/mahakam/pkg/config"

// ResourceNode represents stored resource with terraform backend kind
type ResourceNode struct {
	BaseResource
	Name                  string
	IPAddress             string
	MacAddress            string
	NetMask               string
	Host                  string
	ImageSourcePath       string
	DNSDhcpServerUsername string
	UserData              string
	DNSAddress            string
	GatewayAddress        string
	GateNssAPIKEY         string
	CPU                   int
	Memory                int
	DiskSize              int
	Password              string
}

// NewResourceNode creates new resource node
func NewResourceNode(name string) *ResourceNode {
	return &ResourceNode{
		BaseResource: BaseResource{
			Name:  name,
			Kind:  string(KindTerraformBackend),
			Owner: config.ResourceOwnerGojek,
		},
	}
}
