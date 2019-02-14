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
	// Name is identifier for a libvirt VM
	Name string
	// NumCPUs is number of CPU's assigned to a VM
	NumCPUs string
	// Memory is amount of memory allocated for
	// a VM in Gigabytes
	Memory string
	// DiskSize is the amount of disk for a VM
	// in Gigabytes
	DiskSize int32
	// SSHPublicKey is ssh public key for a VM
	SSHPublicKey string
	// NetworkConfig is network configs for a VM
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
