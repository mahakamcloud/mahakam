package resource

import (
	"net"

	"github.com/mahakamcloud/mahakam/pkg/config"
)

const (
	SmallNode      = "SmallNode"
	MediumNode     = "MediumNode"
	LargeNode      = "LargeNode"
	ExtraLargeNode = "ExtraLargeNode"
)

// NodeSpec represents types of nodes
type NodeSpec struct {
	nCPUCores     int64
	memoryInBytes int64
}

// NewNodeSpec creates new NodeSpec
func NewNodeSpec(nodeType string) *NodeSpec {
	if nodeType == SmallNode {
		return &NodeSpec{2, 4}
	} else if nodeType == MediumNode {
		return &NodeSpec{4, 16}
	} else if nodeType == LargeNode {
		return &NodeSpec{8, 32}
	} else if nodeType == ExtraLargeNode {
		return &NodeSpec{16, 64}
	}

	return nil
}

// Disk represents disk attached to a node
type Disk struct {
	size     string
	disktype string
}

// NetworkConfig represents network config of a node
type NetworkConfig struct {
	macAddress string
	ipAddress  net.IP
	ipMask     net.IPMask
	gateway    net.IP
	nameserver net.IP
}

// NodeStatus represents status of a node resource
type NodeStatus struct {
	host  string
	state string
}

// Metadata for node resource
type Metadata struct {
	userData    string
	sshKeys     []string
	extraConfig map[string]string
}

// Node represents stored resource
type Node struct {
	BaseResource
	nodeSpec        NodeSpec
	networkConfig   NetworkConfig
	nodeStatus      NodeStatus
	metadata        Metadata
	bootDisk        Disk
	additionalDisks []Disk
}

// NewNode creates new Node resource
func NewNode(name string, ns NodeSpec, nc NetworkConfig) *Node {
	return &Node{
		BaseResource: BaseResource{
			Name:  name,
			Kind:  string(KindTerraformBackend),
			Owner: config.ResourceOwnerGojek,
		},
	}
}

// WithLabels attaches labels metadata to Node resource
func (p *Node) WithLabels(labels []Label) *Node {
	p.Labels = labels
	return p
}
