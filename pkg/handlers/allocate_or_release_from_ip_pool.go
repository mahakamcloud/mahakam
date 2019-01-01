package handlers

import (
	"fmt"

	"github.com/go-openapi/swag"

	"github.com/go-openapi/runtime/middleware"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/networks"
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/network"
	r "github.com/mahakamcloud/mahakam/pkg/resource_store/resource"
	log "github.com/sirupsen/logrus"
)

// AllocateOrReleaseFromIPPool is handlers for allocate-or-release from ip pool operation
type AllocateOrReleaseFromIPPool struct {
	Handlers
	log log.FieldLogger
}

func NewAllocateOrReleaseFromIPPool(handlers Handlers) *AllocateOrReleaseFromIPPool {
	log := log.WithField("allocate-or-release", "ip-pool")
	return &AllocateOrReleaseFromIPPool{
		Handlers: handlers,
		log:      log,
	}
}

// Handle is handler for allocate-or-release from ip pool operation
func (h *AllocateOrReleaseFromIPPool) Handle(params networks.AllocateOrReleaseFromIPPoolParams) middleware.Responder {
	if swag.StringValue(params.Action) == config.IPPoolActionAllocate {
		allocatedIP, err := h.allocateIP()
		if err != nil {
			return networks.NewAllocateOrReleaseFromIPPoolDefault(405).WithPayload(&models.Error{
				Code:    405,
				Message: fmt.Sprintf("error allocating ip from ip pool '%v': %s", params, err),
			})
		}
		return networks.NewAllocateOrReleaseFromIPPoolCreated().WithPayload(&models.AllocatedIPPool{
			AllocatedIP: allocatedIP,
		})
	}
	return networks.NewAllocateOrReleaseFromIPPoolDefault(405)
}

func (h *AllocateOrReleaseFromIPPool) allocateIP() (string, error) {
	// TODO(giri): replace first IP pool with the one from IP pool ID
	ipPoolPath, err := h.getFirstIPPool()
	if err != nil {
		log.Errorf("error retrieving ip pool: %s", err)
		return "", err
	}

	p := r.NewResourceIPPool(network.ParseSubnetCIDR(ipPoolPath))
	err = h.Handlers.Store.GetFromPath(ipPoolPath, p)
	if err != nil {
		return "", fmt.Errorf("error getting ip pool resource from kvstore: %s", err)
	}

	if len(p.AvailableIPPools) == 0 {
		return "", fmt.Errorf("running out of available ip pools %v", p)
	}

	ipPools := p.AvailableIPPools
	allocatedIP, ipPools := ipPools[len(ipPools)-1], ipPools[:len(ipPools)-1]
	p.AvailableIPPools = ipPools
	p.AllocatedIPPools = append(p.AllocatedIPPools, allocatedIP)

	_, err = h.Handlers.Store.UpdateFromPath(ipPoolPath, p)
	if err != nil {
		return "", fmt.Errorf("Error updating network subnet resource into kvstore: %s", err)
	}

	return allocatedIP, nil
}

func (h *AllocateOrReleaseFromIPPool) getFirstIPPool() (string, error) {
	keys, err := h.Handlers.Store.ListKeysFromPath(config.KeyPathNetworkIPPool)
	if err != nil || len(keys) == 0 {
		return "", fmt.Errorf("no ip pool exists in kvstore")
	}

	return keys[0], nil
}
