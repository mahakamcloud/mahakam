package node

import "net"

// Node represents VM node metadata
type Node struct {
	Name string
	NetworkConfig
}

// NetworkConfig represents network config of a node
type NetworkConfig struct {
	MacAddress string
	IP         net.IP
	Gateway    net.IP
	Nameserver net.IP
}
