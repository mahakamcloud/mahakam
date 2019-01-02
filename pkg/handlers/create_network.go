package handlers

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	netclient "github.com/mahakamcloud/mahakam/pkg/api/v1/client/networks"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/networks"
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/network"
	"github.com/mahakamcloud/mahakam/pkg/node"
	"github.com/mahakamcloud/mahakam/pkg/provisioner"
	"github.com/mahakamcloud/mahakam/pkg/task"
	"github.com/mahakamcloud/mahakam/pkg/utils"
	log "github.com/sirupsen/logrus"
)

// CreateNetwork is handlers for create-network operation
type CreateNetwork struct {
	Handlers
	log log.FieldLogger
}

func NewCreateNetworkHandler(handlers Handlers) *CreateNetwork {
	log := log.WithField("create", "network")
	return &CreateNetwork{
		Handlers: handlers,
		log:      log,
	}
}

// Handle is handler for create-network operation
func (h *CreateNetwork) Handle(params networks.CreateNetworkParams) middleware.Responder {
	nwf, err := newCreateNetworkWF(params.Body, h)
	if err != nil {
		log.Errorf("error creating network workflow %s", err)
		return networks.NewCreateNetworkDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: fmt.Sprintf("error creating network workflow %s", err),
		})
	}

	// Blocking allocate cluster network,
	// then async network nodes provisioning
	err = nwf.Run()
	if err != nil {
		log.Errorf("error creating network components %s", err)
		return networks.NewCreateNetworkDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: fmt.Sprintf("error creating network components %s", err),
		})
	}

	res := &models.Network{
		Name:        params.Body.Name,
		NetworkName: nwf.clusterNetwork.Name,
		NetworkCIDR: nwf.clusterNetwork.ClusterNetworkCIDR.String(),
		Gateway:     nwf.clusterNetwork.Gateway.String(),
		Nameserver:  nwf.clusterNetwork.Nameserver.String(),
		Dhcp:        nwf.clusterNetwork.Dhcp.String(),
	}
	return networks.NewCreateNetworkCreated().WithPayload(res)
}

type createNetworkWF struct {
	handlers       Handlers
	log            log.FieldLogger
	clusterNetwork *network.ClusterNetwork
	gateway        node.Node
	nameserver     node.Node
	dhcp           node.Node

	dcPublicNetworkCIDR    net.IPNet
	dcGatewayIP            net.IP
	dcNameserverIP         net.IP
	clusterGatewayPublicIP net.IP

	nodePublicKey string
}

func newCreateNetworkWF(cluster *models.Network, cHandler *CreateNetwork) (*createNetworkWF, error) {
	clusterName := swag.StringValue(cluster.Name)

	gateway := node.Node{
		Name: fmt.Sprintf("%s-network-gw", clusterName),
	}

	dhcp := node.Node{
		Name: fmt.Sprintf("%s-network-dhcp", clusterName),
	}

	dns := node.Node{
		Name: fmt.Sprintf("%s-network-dns", clusterName),
	}

	dcGatewayIP, dcPublicNetworkCIDR, _ := net.ParseCIDR(cHandler.AppConfig.NetworkConfig.DatacenterGatewayCIDR)

	dcNameserverIP := net.ParseIP(cHandler.AppConfig.NetworkConfig.DatacenterNameserver)

	allocatedPublicIP, err := getPublicIP()
	if err != nil {
		return nil, fmt.Errorf("error getting public IP allocation: %s", err)
	}

	return &createNetworkWF{
		handlers:               cHandler.Handlers,
		log:                    cHandler.log,
		gateway:                gateway,
		dhcp:                   dhcp,
		nameserver:             dns,
		dcPublicNetworkCIDR:    *dcPublicNetworkCIDR,
		dcGatewayIP:            dcGatewayIP,
		dcNameserverIP:         dcNameserverIP,
		clusterGatewayPublicIP: allocatedPublicIP,
	}, nil
}

func (cn *createNetworkWF) Run() error {
	n, err := cn.handlers.Network.AllocateClusterNetwork()
	if err != nil {
		log.Errorf("cluster network allocation failed %v: %s", cn, err)
		return err
	}
	log.Infof("cluster network has been allocated %v", n)
	cn.clusterNetwork = n

	tasks, err := cn.getCreateTask()
	if err != nil {
		return err
	}

	go func(taskList []task.Task) {
		for _, t := range taskList {
			if err := t.Run(); err != nil {
				cn.log.Errorf("error running task %v: %s", t, err)
			}
		}
	}(tasks)

	return nil
}

func (cn *createNetworkWF) getCreateTask() ([]task.Task, error) {
	var tasks []task.Task
	tasks = cn.setupNetworkGateway(tasks)
	tasks = cn.setupNetworkDHCP(tasks)
	// tasks = cn.setupNetworkNameserver(tasks)
	return tasks, nil
}

func (cn *createNetworkWF) setupNetworkGateway(tasks []task.Task) []task.Task {
	gwConfig := node.NodeCreateConfig{
		Host: net.ParseIP("10.30.0.1"),
		Role: node.RoleNetworkGW,
		Node: node.Node{
			Name:         cn.gateway.Name,
			SSHPublicKey: cn.nodePublicKey,
			NetworkConfig: node.NetworkConfig{
				MacAddress: network.GenerateMacAddress(),
				IP:         cn.clusterNetwork.Gateway,
				Mask:       cn.clusterNetwork.ClusterNetworkCIDR.Mask,
			},
			ExtraNetworks: []node.NetworkConfig{
				node.NetworkConfig{
					IP:         cn.clusterGatewayPublicIP,
					Mask:       cn.dcPublicNetworkCIDR.Mask,
					Gateway:    cn.dcGatewayIP,
					Nameserver: cn.dcNameserverIP,
				},
			},
		},
		ExtraConfig: map[string]string{
			config.KeyClusterNetworkCidr: cn.clusterNetwork.ClusterNetworkCIDR.String(),
		},
	}

	createGatewayNode := provisioner.NewCreateNode(gwConfig, cn.handlers.Provisioner, cn.log)

	tasks = append(tasks, createGatewayNode)

	return tasks
}

func (cn *createNetworkWF) setupNetworkDHCP(tasks []task.Task) []task.Task {
	netCIDR := cn.clusterNetwork.ClusterNetworkCIDR
	dhcpConfig := node.NodeCreateConfig{
		Host: net.ParseIP("10.30.0.1"),
		Role: node.RoleNetworkDHCP,
		Node: node.Node{
			Name:         cn.dhcp.Name,
			SSHPublicKey: cn.nodePublicKey,
			NetworkConfig: node.NetworkConfig{
				MacAddress: network.GenerateMacAddress(),
				IP:         cn.clusterNetwork.Dhcp,
				Mask:       netCIDR.Mask,
				Gateway:    cn.clusterNetwork.Gateway,
				Nameserver: cn.dcNameserverIP,
			},
		},
		ExtraConfig: map[string]string{
			config.KeyClusterNetworkCidr: netCIDR.String(),
			config.KeySubnetAddress:      netCIDR.IP.String(),
			config.KeySubnetMask:         utils.IPv4MaskString(netCIDR.Mask),
			config.KeyBroadcastAddress:   broadcastAddr(netCIDR),
		},
	}

	createDHCPNode := provisioner.NewCreateNode(dhcpConfig, cn.handlers.Provisioner, cn.log)

	tasks = append(tasks, createDHCPNode)

	return tasks
}

func (cn *createNetworkWF) setupNetworkNameserver(tasks []task.Task) []task.Task {
	return tasks
}

func getPublicIP() (net.IP, error) {
	client := GetMahakamClient(":" + strconv.Itoa(config.MahakamAPIDefaultPort))
	res, err := client.Networks.AllocateOrReleaseFromIPPool(netclient.NewAllocateOrReleaseFromIPPoolParams().
		WithAction(swag.String(config.IPPoolActionAllocate)))
	if err != nil {
		return nil, err
	}

	return net.ParseIP(res.Payload.AllocatedIP), nil
}

func broadcastAddr(n net.IPNet) string {
	if n.IP.To4() == nil {
		return ""
	}
	ip := make(net.IP, len(n.IP.To4()))
	binary.BigEndian.PutUint32(ip, binary.BigEndian.Uint32(n.IP.To4())|^binary.BigEndian.Uint32(net.IP(n.Mask).To4()))
	return ip.String()
}
