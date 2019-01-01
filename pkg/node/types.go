package node

import "net"

// Role of Kubernetes node
type Role string

const (
	// RoleControlPlane defines a k8s control plane node
	RoleControlPlane Role = "control-plane"
	// RoleWorker defines a k8s worker node
	RoleWorker Role = "worker"
	// RoleNetworkDNS defines a dns node in a network
	RoleNetworkDNS Role = "network-dns"
	// RoleNetworkDHCP defines a dhcp node in a network
	RoleNetworkDHCP Role = "network-dhcp"
	// RoleNetworkGW defines a gateway node in a network
	RoleNetworkGW Role = "network-gw"
)

// Node represents VM node metadata
type Node struct {
	Name         string
	NumCPUs      int32
	MemoryMB     int64
	DiskSizeGB   int32
	SSHPublicKey string
	NetworkConfig
	ExtraNetworks []NetworkConfig
}

// NetworkConfig represents network config of a node
type NetworkConfig struct {
	MacAddress string
	IP         net.IP
	Mask       net.IPMask
	Gateway    net.IP
	Nameserver net.IP
}

// NodeCreateConfig defines config for creating node within libvirt Datacenter
type NodeCreateConfig struct {
	Host net.IP
	Node
	Role        Role
	ExtraConfig map[string]string
}
