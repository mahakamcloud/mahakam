package resource

import (
	"net"

	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/utils"
)

// ResourceIPPool represents stored ip pool kind
type ResourceIPPool struct {
	*BaseResource
	Subnet           string   `json:"subnet"`
	SubnetLen        string   `json:"subnet_len"`
	IPPoolRangeStart string   `json:"ip_pool_range_start"`
	IPPoolRangeEnd   string   `json:"ip_pool_range_end"`
	AvailableIPPools []string `json:"available_ip_pools"`
	AllocatedIPPools []string `json:"allocated_ip_pools"`
}

func NewResourceIPPool(cidr net.IPNet) *ResourceIPPool {
	return &ResourceIPPool{
		BaseResource: &BaseResource{
			Name:  utils.CidrToKeyString(cidr),
			Kind:  string(KindIPPool),
			Owner: config.ResourceOwnerMahakam,
		},
		Subnet:    cidr.String(),
		SubnetLen: cidr.Mask.String(),
	}
}

func (p *ResourceIPPool) WithLabels(labels []Label) *ResourceIPPool {
	p.Labels = labels
	return p
}

type ResourceIPPoolList struct {
	Items []*ResourceIPPool
}

func NewResourceIPPoolList() *ResourceIPPoolList {
	return &ResourceIPPoolList{}
}

func (l *ResourceIPPoolList) Resource() Resource {
	return &ResourceIPPool{}
}

func (l *ResourceIPPoolList) WithItems(items []Resource) *ResourceIPPoolList {
	for _, i := range items {
		ipPool := i.(*ResourceIPPool)
		l.Items = append(l.Items, ipPool)
	}
	return l
}
