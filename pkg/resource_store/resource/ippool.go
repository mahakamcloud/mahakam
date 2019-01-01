package resource

import (
	"net"

	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/utils"
)

// ResourceIPPool represents stored ip pool kind
type ResourceIPPool struct {
	BaseResource
	Subnet           string   `json:"subnet"`
	SubnetLen        string   `json:"subnet_len"`
	IPPoolRangeStart string   `json:"ip_pool_range_start"`
	IPPoolRangeEnd   string   `json:"ip_pool_range_end"`
	AvailableIPPools []string `json:"available_ip_pools"`
	AllocatedIPPools []string `json:"allocated_ip_pools"`
}

func NewResourceIPPool(cidr net.IPNet) *ResourceIPPool {
	return &ResourceIPPool{
		BaseResource: BaseResource{
			Name:  utils.CidrToKeyString(cidr),
			Kind:  string(KindIPPool),
			Owner: config.ResourceOwnerMahakam,
		},
		Subnet:    cidr.String(),
		SubnetLen: cidr.Mask.String(),
	}
}
