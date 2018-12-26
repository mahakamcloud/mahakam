package node

import "net"

// Node represents VM node metadata
type Node struct {
	Name       string
	NumCPUs    int32
	MemoryMB   int64
	DiskSizeGB int32
	NetworkConfig
}

// NetworkConfig represents network config of a node
type NetworkConfig struct {
	MacAddress string
	IP         net.IP
	Gateway    net.IP
	Nameserver net.IP
}

// NodeCreateConfig defines config for creating node within libvirt Datacenter
type NodeCreateConfig struct {
	Host net.IP
	Node
	ExtraConfig map[string]string
}
