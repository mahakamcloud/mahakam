package resource

import (
	"net"

	"github.com/mahakamcloud/mahakam/pkg/config"
)

// Type represents types of nodes
type Type string

// Disk represents disk attached to a node
type Disk struct {
	Size string
	Type string
}

// NetworkConfig represents network config of a node
type NetworkConfig struct {
	MacAddress string
	IP         net.IP
	Mask       net.IPMask
	Gateway    net.IP
	Nameserver net.IP
}

// NodeStatus represents status of a node resource
type NodeStatus struct {
	Host  string
	State string
}

// Metadata for node resource
type Metadata struct {
	UserData    string
	SSHKeys     []string
	ExtraConfig map[string]string
}

// Node represents stored resource
type Node struct {
	BaseResource
	Type
	NetworkConfig
	NodeStatus
	Metadata
	BootDisk        Disk
	AdditionalDisks []Disk
	Name            string
}

// NewNode creates new Node resource
func NewNode(name string) *Node {
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
