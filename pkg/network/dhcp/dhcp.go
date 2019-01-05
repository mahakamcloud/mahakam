package dhcp

import (
	"net"

	"strings"
	"text/template"

	"github.com/hashicorp/consul/api"
	"github.com/mahakamcloud/mahakam/pkg/node"
	log "github.com/sirupsen/logrus"
)

const (
	consulDHCPRootKey = "hosts/"
	consulPort        = "8500"
	dhcpTemplate      = `
host {{ .NodeHostName }} {
	hardware ethernet {{ .HardwareAddress }};
	fixed-address {{ .NodeAddress }};
	option host-name {{ .NodeHostName }};
	option domain-name-servers {{ .DNS }};
	option domain-name "mahakam.gocloud.io";
	option routers {{ .Gateway }};
}
`
)

// EndpointConfig describes the consul endpoint of DHCP node
type EndpointConfig struct {
	IPAddress  net.IP
	ConsulPort string
}

// NodeDHCPConfig describes dhcp config for a host to be appended in dhcpd.conf
type NodeDHCPConfig struct {
	NodeHostName    string
	NodeAddress     net.IP
	HardwareAddress string
	Gateway         net.IP
	DNS             net.IP
}

// NewEndpointConfig returns EndpointConfig for a DHCP
func NewEndpointConfig(dhcpIP net.IP) EndpointConfig {
	return EndpointConfig{
		IPAddress:  dhcpIP,
		ConsulPort: consulPort,
	}
}

// NewNodeConfig returns a NodeDHCPConfig to be registered on DHCP
func NewNodeConfig(nodeConfig node.NodeCreateConfig) *NodeDHCPConfig {
	return &NodeDHCPConfig{
		NodeHostName:    nodeConfig.Node.Name,
		NodeAddress:     nodeConfig.Node.NetworkConfig.IP,
		HardwareAddress: nodeConfig.Node.NetworkConfig.MacAddress,
		Gateway:         nodeConfig.Node.NetworkConfig.Gateway,
		DNS:             nodeConfig.Node.NetworkConfig.Nameserver,
	}
}

// Register registers a node on DHCP
func (nodeDhcpConfig *NodeDHCPConfig) Register(dhcpEpConfig EndpointConfig) error {

	nodeDhcpData, err := nodeDhcpConfig.parse()
	if err != nil {
		log.Errorf("Node dhcp config parse error : %v", err)
		return err
	}

	kv, err := dhcpEpConfig.getConsulEndpoint()
	if err != nil {
		log.Errorf("Error getting consul KV handle : %v", err)
		return err
	}

	// PUT a new KV pair
	p := &api.KVPair{Key: consulDHCPRootKey + nodeDhcpConfig.NodeHostName, Value: []byte(nodeDhcpData)}
	_, err = kv.Put(p, nil)
	if err != nil {
		log.Errorf("Error writing data to consul : %v", err)
		return err
	}
	return nil
}

func (nodeDhcpConfig *NodeDHCPConfig) parse() (string, error) {
	dhcpTmplVal := template.New("dhcp")
	parsedDHCPTmpl, err := dhcpTmplVal.Parse(dhcpTemplate)
	if err != nil {
		log.Errorf("Error parsing the DHCP Template: %v", err)
		return "", err
	}

	var data strings.Builder
	err = parsedDHCPTmpl.Execute(&data, nodeDhcpConfig)
	if err != nil {
		log.Errorf("Error executing the DHCP template: %v", err)
		return "", err
	}

	return data.String(), nil
}

func (dhcpEpConfig *EndpointConfig) getConsulEndpoint() (*api.KV, error) {
	config := &api.Config{
		Address: dhcpEpConfig.IPAddress.String() + ":" + dhcpEpConfig.ConsulPort,
		Scheme:  "http",
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	// Get a handle to the KV API
	kv := client.KV()
	return kv, nil
}
