package resource

import (
	"net"
	"strings"

	"github.com/mahakamcloud/mahakam/pkg/config"
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
			Name:  cidrToKeyString(cidr),
			Kind:  string(KindNetwork),
			Owner: config.ResourceOwnerMahakam,
		},
		Subnet:    cidr.String(),
		SubnetLen: cidr.Mask.String(),
	}
}

func cidrToKeyString(cidr net.IPNet) string {
	return strings.Replace(cidr.String(), "/", "-", -1)
}
