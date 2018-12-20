package resource

import (
	"net"

	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/utils"
)

// ResourceNetwork represents stored resource with network kind
type ResourceNetwork struct {
	BaseResource
	Subnet           string   `json:"subnet"`
	SubnetLen        string   `json:"subnet_len"`
	AvailableIPPools []string `json:"available_ip_pools"`
	AllocatedIPPools []string `json:"allocated_ip_pools"`
}

// NewResourceNetwork creates new resource network
func NewResourceNetwork(cidr net.IPNet) *ResourceNetwork {
	return &ResourceNetwork{
		BaseResource: BaseResource{
			Name:  utils.CidrToKeyString(cidr),
			Kind:  string(KindNetwork),
			Owner: config.ResourceOwnerMahakam,
		},
		Subnet:    cidr.String(),
		SubnetLen: cidr.Mask.String(),
	}
}
