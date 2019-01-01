package handlers

import (
	"fmt"
	"net"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/networks"
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/network"
	"github.com/mahakamcloud/mahakam/pkg/node"
	"github.com/mahakamcloud/mahakam/pkg/provisioner"
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

	nodePublicKey string
}

func newCreateNetworkWF(cluster *models.Network, cHandler *CreateNetwork) (*createNetworkWF, error) {
	clusterName := swag.StringValue(cluster.Name)

	gateway := node.Node{
		Name: fmt.Sprintf("%s-network-gw", clusterName),
	}

	return &createNetworkWF{
		handlers: cHandler.Handlers,
		log:      cHandler.log,
		gateway:  gateway,
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

	go func(taskList []provisioner.Task) {
		for _, t := range taskList {
			if err := t.Run(); err != nil {
				cn.log.Errorf("error running task %v: %s", t, err)
			}
		}
	}(tasks)

	return nil
}

func (cn *createNetworkWF) getCreateTask() ([]provisioner.Task, error) {
	var tasks []provisioner.Task
	tasks = cn.setupNetworkGateway(tasks)
	tasks = cn.setupNetworkDHCP(tasks)
	// tasks = cn.setupNetworkNameserver(tasks)
	return tasks, nil
}

func (cn *createNetworkWF) setupNetworkGateway(tasks []provisioner.Task) []provisioner.Task {
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
					// TODO(giri): pass proper public IP from config.yaml
					// IP:         net.ParseIP("1.2.3.4"),
					// Mask:       net.CIDRMask(28, 32),
					// Gateway:    net.ParseIP("1.2.3.1"),
					// Nameserver: net.ParseIP("8.8.8.8"),
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

func (cn *createNetworkWF) setupNetworkDHCP(tasks []provisioner.Task) []provisioner.Task {
	dhcpConfig := node.NodeCreateConfig{
		Host: net.ParseIP("10.30.0.1"),
		Role: node.RoleNetworkDHCP,
		Node: node.Node{
			Name:         cn.dhcp.Name,
			SSHPublicKey: cn.nodePublicKey,
			NetworkConfig: node.NetworkConfig{
				MacAddress: network.GenerateMacAddress(),
				IP:         cn.clusterNetwork.Dhcp,
				Mask:       cn.clusterNetwork.ClusterNetworkCIDR.Mask,
				Gateway:    cn.clusterNetwork.Gateway,
			},
		},
	}

	createDHCPNode := provisioner.NewCreateNode(dhcpConfig, cn.handlers.Provisioner, cn.log)

	tasks = append(tasks, createDHCPNode)

	return tasks
}

func (cn *createNetworkWF) setupNetworkNameserver(tasks []provisioner.Task) []provisioner.Task {
	return tasks
}
