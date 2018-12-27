package node

import "net"

// Role of Kubernetes node
type Role string

const (
	RoleControlPlane Role = "control-plane"
	RoleWorker       Role = "worker"
)

// Node represents VM node metadata
type Node struct {
	Name         string
	NumCPUs      int32
	MemoryMB     int64
	DiskSizeGB   int32
	SSHPublicKey string
	NetworkConfig
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
