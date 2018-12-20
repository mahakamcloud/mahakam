package network

import (
	"fmt"
	"net"
	"strings"

	"github.com/mahakamcloud/mahakam/pkg/config"
	store "github.com/mahakamcloud/mahakam/pkg/resource_store"
	r "github.com/mahakamcloud/mahakam/pkg/resource_store/resource"
)

type ClusterNetwork struct {
	ClusterNetworkCIDR net.IPNet
}

type NetworkManager struct {
	store store.ResourceStore
	net   config.NetworkConfig
	sm    *SubnetManager

	networkIP   net.IP
	networkCIDR *net.IPNet
}

func New(rs store.ResourceStore, nc config.NetworkConfig) (*NetworkManager, error) {
	nm := &NetworkManager{
		store: rs,
		net:   nc,
	}

	if err := parseConfig(nm); err != nil {
		return nil, fmt.Errorf("Error initializing network manager: %s", err)
	}

	nm.sm = NewSubnetManager(nm.networkCIDR)

	return nm, nil
}

func parseConfig(nm *NetworkManager) error {
	ip, ipnet, err := net.ParseCIDR(nm.net.CIDR)
	if err != nil {
		return fmt.Errorf("Error parsing network CIDR: %s", err)
	}
	nm.networkIP = ip
	nm.networkCIDR = ipnet

	return nil
}

func (nm *NetworkManager) AllocateClusterNetwork() (*ClusterNetwork, error) {

	reservedSubnets, err := nm.getReservedSubnets()
	if err != nil {
		return nil, fmt.Errorf("Error getting reserved subnets: %s", err)
	}

	clusterNetworkCIDR, err := nm.sm.CreateSubnet(nm.net.ClusterNetmask, reservedSubnets)
	if err != nil {
		return nil, fmt.Errorf("Error allocating cluster network: %s", err)
	}

	n := r.NewResourceNetwork(clusterNetworkCIDR)
	n.AvailableIPPools = nm.allocateIPPools(clusterNetworkCIDR)

	_, err = nm.store.AddFromPath(config.KeyPathNetworkSubnet+n.Name, n)
	if err != nil {
		return nil, fmt.Errorf("Error storing network resource to kvstore: %s", err)
	}

	return &ClusterNetwork{
		ClusterNetworkCIDR: clusterNetworkCIDR,
	}, nil
}

// TODO(giri): Implement release network
func (nm *NetworkManager) ReleaseClusterNetwork() error {
	return nil
}

func (nm *NetworkManager) getReservedSubnets() ([]net.IPNet, error) {
	reservedSubnets := []net.IPNet{}

	keys, err := nm.store.ListKeysFromPath(config.KeyPathNetworkSubnet)
	if err != nil {
		return nil, fmt.Errorf("Error listing keys from path: %s", err)
	}

	for _, key := range keys {
		reservedSubnets = append(reservedSubnets, nm.parseSubnetCIDR(key))
	}

	return reservedSubnets, nil
}

func (nm *NetworkManager) allocateIPPools(cidr net.IPNet) []string {
	var ips []string
	for ip := cidr.IP.Mask(cidr.Mask); cidr.Contains(ip); nm.nextIP(ip) {
		ips = append(ips, ip.String())
	}

	// First and last IP should not be on the IP pools
	// Start from behind to efficiently allocate/release IP
	ipPools := ips[1 : len(ips)-1]
	nm.reverseIPPools(ipPools)

	return ipPools
}

func (nm *NetworkManager) nextIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}

func (nm *NetworkManager) reverseIPPools(ipPools []string) {
	for i, j := 0, len(ipPools)-1; i < j; i, j = i+1, j-1 {
		ipPools[i], ipPools[j] = ipPools[j], ipPools[i]
	}
}

func (nm *NetworkManager) parseSubnetCIDR(key string) net.IPNet {
	cidrKeys := strings.Split(key, "/")
	cidr := strings.Replace(cidrKeys[len(cidrKeys)-1], "-", "/", -1)
	_, ipnet, _ := net.ParseCIDR(cidr)
	return *ipnet
}
