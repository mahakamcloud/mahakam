package resource

import "github.com/mahakamcloud/mahakam/pkg/config"

// Node represents stored resource with terraform backend kind
type Node struct {
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

// NewNode creates new resource node
func NewNode(name string) *Node {
	return &Node{
		BaseResource: BaseResource{
			Name:  name,
			Kind:  string(KindTerraformBackend),
			Owner: config.ResourceOwnerGojek,
		},
	}
}
