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
	"github.com/mahakamcloud/mahakam/pkg/scheduler"
	"github.com/mahakamcloud/mahakam/pkg/task"
	"github.com/mahakamcloud/mahakam/pkg/utils"
	"github.com/sirupsen/logrus"
)

// CreateNetwork is handlers for create-network operation
type CreateNetwork struct {
	Handlers
	hosts []config.Host
	log   logrus.FieldLogger
}

// NewCreateNetworkHandler creates a CreateNetwork object
func NewCreateNetworkHandler(handlers Handlers) *CreateNetwork {
	return &CreateNetwork{
		Handlers: handlers,
		hosts:    handlers.AppConfig.HostsConfig,
		log:      handlers.Log,
	}
}

// Handle is handler for create-network operation
func (h *CreateNetwork) Handle(params networks.CreateNetworkParams) middleware.Responder {
	h.log.Infof("handling create network request: %v", params)

	nwf, err := newCreateNetworkWF(params.Body, h)
	if err != nil {
		h.log.Errorf("error creating network workflow %s", err)
		return networks.NewCreateNetworkDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: fmt.Sprintf("error creating network workflow %s", err),
		})
	}

	// Blocking allocate cluster network,
	// then async network nodes provisioning
	err = nwf.Run()
	if err != nil {
		h.log.Errorf("error creating network components %s", err)
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
	log            logrus.FieldLogger
	clusterNetwork *network.ClusterNetwork
	gateway        node.Node
	nameserver     node.Node
	dhcp           node.Node
	hosts          []config.Host

	dcPublicNetworkCIDR    net.IPNet
	dcGatewayIP            net.IP
	dcNameserverIP         net.IP
	clusterGatewayPublicIP net.IP

	nodePublicKey string
}

func newCreateNetworkWF(cluster *models.Network, cHandler *CreateNetwork) (*createNetworkWF, error) {
	nwfLog := cHandler.log.WithField("workflow", "create-network")
	nwfLog.Debugf("init create network workflow")

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
		hosts:                  cHandler.hosts,
		dcPublicNetworkCIDR:    *dcPublicNetworkCIDR,
		dcGatewayIP:            dcGatewayIP,
		dcNameserverIP:         dcNameserverIP,
		clusterGatewayPublicIP: allocatedPublicIP,
	}, nil
}

func (cn *createNetworkWF) Run() error {
	cn.log.Infof("running create network workflow: %v", cn)

	pretasks, err := cn.getPreCreateTask()
	if err != nil {
		return err
	}

	// blocking pre-create tasks
	for _, t := range pretasks {
		cn.log.Infof("running pre-create task %v", t)
		if err := t.Run(); err != nil {
			cn.log.Errorf("error running pre-create task %v: %s", t, err)
		}
	}

	tasks, err := cn.getCreateTask()
	if err != nil {
		return err
	}

	go func(taskList []task.Task) {
		for _, t := range taskList {
			cn.log.Infof("running task %v", t)
			if err := t.Run(); err != nil {
				cn.log.Errorf("error running task %v: %s", t, err)
			}
		}
	}(tasks)

	return nil
}

func (cn *createNetworkWF) getPreCreateTask() ([]task.Task, error) {
	cn.log.Debugf("getting pre-create task for network workflow")

	n, err := cn.handlers.Network.AllocateClusterNetwork()
	if err != nil {
		cn.log.Errorf("cluster network allocation failed %v: %s", cn, err)
		return nil, err
	}
	cn.log.Infof("cluster network has been allocated %v", n)
	cn.clusterNetwork = n

	var tasks []task.Task
	tasks = cn.setupNetworkPreCreateTasks(tasks)
	return tasks, nil
}

func (cn *createNetworkWF) setupNetworkPreCreateTasks(tasks []task.Task) []task.Task {
	cn.log.Debugf("setup network pre-create tasks for network %s", cn.clusterNetwork.Name)

	mahakamServerIP := cn.clusterNetwork.MahakamServer
	mahakamServerMask := cn.clusterNetwork.ClusterNetworkCIDR.Mask
	mahakamNetIf := config.MahakamDefaultNetworkInterface

	networkReachability := provisioner.NewClusterNetworkReachability(utils.NewIPUtil(), mahakamServerIP, mahakamServerMask, mahakamNetIf)
	tasks = append(tasks, networkReachability)
	return tasks
}

func (cn *createNetworkWF) getCreateTask() ([]task.Task, error) {
	cn.log.Debugf("getting create task for network %s", cn.clusterNetwork.Name)

	var tasks []task.Task
	tasks = cn.setupNetworkGatewayTasks(tasks)
	tasks = cn.setupNetworkDHCPTasks(tasks)
	tasks = cn.setupNetworkNameserverTasks(tasks)
	return tasks, nil
}

func (cn *createNetworkWF) setupNetworkGatewayTasks(tasks []task.Task) []task.Task {
	cn.log.Debugf("setup network gateway tasks for network %s", cn.clusterNetwork.Name)

	host, err := scheduler.GetHost(cn.hosts)
	if err != nil {
		cn.log.Errorf("error getting scheduled a host to provision network gateway for %s: %v", cn.clusterNetwork.Name, err)
		return nil
	}

	gwConfig := node.NodeCreateConfig{
		Host: host,
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

func (cn *createNetworkWF) setupNetworkDHCPTasks(tasks []task.Task) []task.Task {
	cn.log.Debugf("setup network dhcp tasks for network %s", cn.clusterNetwork.Name)

	netCIDR := cn.clusterNetwork.ClusterNetworkCIDR

	host, err := scheduler.GetHost(cn.hosts)
	if err != nil {
		cn.log.Errorf("error getting scheduled a host to provision network dhcp for %s: %v", cn.clusterNetwork.Name, err)
		return nil
	}

	dhcpConfig := node.NodeCreateConfig{
		Host: host,
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

	checkNetworkGWNode := provisioner.NewCheckNode(cn.clusterNetwork.Gateway, cn.log, utils.NewPingCheck())
	createDHCPNode := provisioner.NewCreateNode(dhcpConfig, cn.handlers.Provisioner, cn.log)

	dhcpSeqTasks := task.NewSeqTask(cn.log, checkNetworkGWNode, createDHCPNode)

	tasks = append(tasks, dhcpSeqTasks)

	return tasks
}

func (cn *createNetworkWF) setupNetworkNameserverTasks(tasks []task.Task) []task.Task {
	cn.log.Debugf("setup network nameserver tasks for network %s", cn.clusterNetwork.Name)

	netCIDR := cn.clusterNetwork.ClusterNetworkCIDR

	host, err := scheduler.GetHost(cn.hosts)
	if err != nil {
		cn.log.Errorf("error getting scheduled a host to provision network nameserver for %s: %v", cn.clusterNetwork.Name, err)
		return nil
	}

	dnsConfig := node.NodeCreateConfig{
		Host: host,
		Role: node.RoleNetworkDNS,
		Node: node.Node{
			Name:         cn.nameserver.Name,
			SSHPublicKey: cn.nodePublicKey,
			NetworkConfig: node.NetworkConfig{
				MacAddress: network.GenerateMacAddress(),
				IP:         cn.clusterNetwork.Nameserver,
				Mask:       netCIDR.Mask,
				Gateway:    cn.clusterNetwork.Gateway,
				Nameserver: cn.dcNameserverIP,
			},
		},
		ExtraConfig: map[string]string{
			config.KeyClusterNetworkCidr: netCIDR.String(),
		},
	}

	checkNetworkGWNode := provisioner.NewCheckNode(cn.clusterNetwork.Gateway, cn.log, utils.NewPingCheck())
	checkNetworkDHCPNode := provisioner.NewCheckNode(cn.clusterNetwork.Dhcp, cn.log, utils.NewPingCheck())
	createDNSNode := provisioner.NewCreateNode(dnsConfig, cn.handlers.Provisioner, cn.log)

	dnsSeqTasks := task.NewSeqTask(cn.log, checkNetworkGWNode, checkNetworkDHCPNode, createDNSNode)

	tasks = append(tasks, dnsSeqTasks)

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
