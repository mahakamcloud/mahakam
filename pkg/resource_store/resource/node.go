package resource

import (
	"net"
	"strconv"
	"strings"

	"github.com/mahakamcloud/mahakam/pkg/config"
)

// Role is
type Role string

const (
	// RoleK8sControlPlane defines a k8s control plane node
	RoleK8sControlPlane Role = "k8s-control-plane"
	// RoleK8sWorker defines a k8s worker node
	RoleK8sWorker Role = "k8s-worker"
	// RoleNetworkDNS defines a dns node in a network
	RoleNetworkDNS Role = "network-dns"
	// RoleNetworkDHCP defines a dhcp node in a network
	RoleNetworkDHCP Role = "network-dhcp"
	// RoleNetworkGW defines a gateway node in a network
	RoleNetworkGW Role = "network-gw"
	// SmallNode depicts string value of small Node size
	SmallNode = "SMALL"
	// MediumNode depicts string value of medium Node size
	MediumNode = "MEDIUM"
	// LargeNode depicts string value of large Node size
	LargeNode = "LARGE"
	// ExtraLargeNode depicts string value of extra large Node size
	ExtraLargeNode = "EXTRALARGE"
	// DefaultK8sControlPlaneNode depicts default k8s control plane Node size
	DefaultK8sControlPlaneNode = MediumNode
	// DefaultK8sWorkerNode depicts default k8s worker Node size
	DefaultK8sWorkerNode = MediumNode
	// DefaultNetworkDNSNode depicts default DNS Node size
	DefaultNetworkDNSNode = SmallNode
	// DefaultNetworkDHCPNode depicts default DHCP Node size
	DefaultNetworkDHCPNode = MediumNode
	// DefaultNetworkGWNode depicts default Gateway Node size
	DefaultNetworkGWNode = MediumNode
)

var availableNodeSizes = []string{SmallNode, MediumNode, LargeNode, ExtraLargeNode}

// NodeSizeValidate verifies is Node size passed is valid
func NodeSizeValidate(size string) bool {
	for _, availableSize := range availableNodeSizes {
		if strings.ToUpper(size) == availableSize {
			return true
		}
	}
	return false
}

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
	macAddr    string
	ipMask     net.IPMask
	ipAddr     net.IP
	gateway    net.IP
	nameserver net.IP
	fqdn       string
}

// NewNetworkConfig returns the NetworkConfig for a Node
func NewNetworkConfig(macAddr, fqdn string, ipMask net.IPMask, ipAddr, gateway, nameserver net.IP) *NetworkConfig {
	return &NetworkConfig{
		macAddr:    macAddr,
		ipMask:     ipMask,
		ipAddr:     ipAddr,
		gateway:    gateway,
		nameserver: nameserver,
		fqdn:       fqdn,
	}
}

// MacAddress return mac address associated with a NetworkConfig
func (nc *NetworkConfig) MacAddress() string {
	return nc.macAddr
}

// IPAddress return IP address associated with a NetworkConfig
func (nc *NetworkConfig) IPAddress() string {
	return nc.ipAddr.String()
}

// IPMask return IP mask associated with a NetworkConfig
func (nc *NetworkConfig) IPMask() string {
	return nc.ipMask.String()
}

// Gateway return gateway associated with the network
func (nc *NetworkConfig) Gateway() string {
	return nc.gateway.String()
}

// Nameserver return Nameserver associated with a NetworkConfig
func (nc *NetworkConfig) Nameserver() string {
	return nc.nameserver.String()
}

// FQDN return Nameserver associated with a NetworkConfig
func (nc *NetworkConfig) FQDN() string {
	return nc.fqdn
}

// NodeStatus represents status of a node resource
type NodeStatus struct {
	host string
}

// NewNodeStatus return a Node Metadata
func NewNodeStatus(host string) *NodeStatus {
	return &NodeStatus{
		host: host,
	}
}

// Host return the host the Node is scheduled on
func (n *NodeStatus) Host() string {
	return n.host
}

// Metadata for node resource
type Metadata struct {
	userData    string
	sshKeys     []string
	extraConfig map[string]string
}

// NewMetadata return a Node Metadata
func NewMetadata(userData string, sskKeys []string, extraConfig map[string]string) *Metadata {
	return &Metadata{
		userData:    userData,
		sshKeys:     sskKeys,
		extraConfig: extraConfig,
	}
}

// SSHKeys return SSH Keys embedded in Metadata struct
func (m *Metadata) SSHKeys() []string {
	return m.sshKeys
}

// Node represents stored resource
type Node struct {
	BaseResource
	nodeSpec       NodeSpec
	networkConfigs []NetworkConfig
	status         NodeStatus
	metadata       Metadata
	disks          []Disk
}

// NewNode creates new Node resource
func NewNode(name string, ns NodeSpec, nc []NetworkConfig, m Metadata, s NodeStatus, role string) *Node {
	return &Node{
		BaseResource: BaseResource{
			Name:  name,
			Kind:  string(KindTerraformBackend),
			Owner: config.ResourceOwnerGojek,
			Labels: []Label{
				Label{
					Key:   "Role",
					Value: role,
				},
			},
		},
		nodeSpec:       ns,
		networkConfigs: nc,
		metadata:       m,
		status:         s,
	}
}

// NetworkConfigs returns Node's network configs
func (n *Node) NetworkConfigs() []NetworkConfig {
	return n.networkConfigs
}

// Role returns role of a node
func (n *Node) Role() string {
	for _, i := range n.BaseResource.Labels {
		if i.Key == "Role" {
			return i.Value
		}
	}

	return ""
}

// ExtraConfig returns extraconfig of a node
func (n *Node) ExtraConfig() map[string]string {
	return n.metadata.extraConfig
}

// Memory returns Memory associated with a node
func (n *Node) Memory() string {
	return strconv.FormatInt(n.nodeSpec.memoryInBytes, 10)
}

// CPUs returns numbner of CPUs associated with a node
func (n *Node) CPUs() string {
	return strconv.FormatInt(n.nodeSpec.nCPUCores, 10)
}

// Status returns Node status
func (n *Node) Status() *NodeStatus {
	return &n.status
}

// Metadata returns Node metadata
func (n *Node) Metadata() *Metadata {
	return &n.metadata
}

// WithLabels attaches labels metadata to Node resource
func (n *Node) WithLabels(labels []Label) *Node {
	n.Labels = labels
	return n
}
