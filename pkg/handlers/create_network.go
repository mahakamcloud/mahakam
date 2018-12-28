package handlers

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/networks"
	"github.com/mahakamcloud/mahakam/pkg/network"
	log "github.com/sirupsen/logrus"
)

// CreateNetwork is handlers for create-network operation
type CreateNetwork struct {
	Handlers
}

// Handle is handler for create-network operation
func (h *CreateNetwork) Handle(params networks.CreateNetworkParams) middleware.Responder {
	nwf, err := newCreateNetworkWF(params.Body, h.Handlers)
	if err != nil {
		log.Errorf("error creating network workflow %s", err)
		return networks.NewCreateNetworkDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: fmt.Sprintf("error creating network workflow %s", err),
		})
	}

	// Blocking allocate cluster network,
	// then parallel network nodes provisioning
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
	}
	return networks.NewCreateNetworkCreated().WithPayload(res)
}

// TODO(giri/vijay): create network
type createNetworkWF struct {
	handlers       Handlers
	clusterNetwork *network.ClusterNetwork
}

func newCreateNetworkWF(cluster *models.Network, handlers Handlers) (*createNetworkWF, error) {
	return &createNetworkWF{
		handlers: handlers,
	}, nil
}

func (cn *createNetworkWF) Run() error {
	// TODO(giri/vijay): create network components
	n, err := cn.handlers.Network.AllocateClusterNetwork()
	if err != nil {
		log.Errorf("cluster network creation failed %v: %s", cn, err)
		return err
	}
	log.Infof("cluster network has been created %v", n)

	cn.clusterNetwork = n

	return nil
}
