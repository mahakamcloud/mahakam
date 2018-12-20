package network

import (
	"context"
	"fmt"
	"net"

	"github.com/giantswarm/ipam"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/microstorage/memory"
)

const (
	subnetMaskBits = 32
)

// SubnetManager represents manager that can create/delete subnet
type SubnetManager struct {
	subnet *ipam.Service
	ctx    context.Context
}

// NewSubnetManager depends on giantswarm's library
// for create and delete subnet from specific network cidr.
// This library has hard dependencies on micrologger and microstorage.
func NewSubnetManager(networkCIDR *net.IPNet) *SubnetManager {
	ctx := context.Background()

	config := subnetManagerConfig(networkCIDR)
	subnetService, _ := ipam.New(*config)

	return &SubnetManager{
		subnet: subnetService,
		ctx:    ctx,
	}
}

// subnetManagerConfig reprensets giantswarm's library config
func subnetManagerConfig(networkCIDR *net.IPNet) *ipam.Config {
	l, _ := micrologger.New(micrologger.Config{})
	s, _ := memory.New(memory.Config{})

	return &ipam.Config{
		Logger:  l,
		Storage: s,
		Network: networkCIDR,
	}
}

// CreateSubnet creates new subnet from registered networkCIDR.
// Can also avoid subnets that are already reserved.
func (sm *SubnetManager) CreateSubnet(mask int, reserved []net.IPNet) (net.IPNet, error) {
	network, err := sm.subnet.CreateSubnet(sm.ctx, net.CIDRMask(mask, subnetMaskBits), "", reserved)
	if err != nil {
		return net.IPNet{}, fmt.Errorf("Error creating new subnet: %s", err)
	}
	return network, nil
}

// DeleteSubnet deletes specified subnet from the pools.
func (sm *SubnetManager) DeleteSubnet(subnet net.IPNet) error {
	err := sm.subnet.DeleteSubnet(sm.ctx, subnet)
	if err != nil {
		return fmt.Errorf("Error creating new subnet: %s", err)
	}
	return nil
}
