package handlers

import (
	"fmt"
	"net"

	"github.com/go-openapi/runtime/middleware"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/networks"
	"github.com/mahakamcloud/mahakam/pkg/config"
	r "github.com/mahakamcloud/mahakam/pkg/resource_store/resource"
	"github.com/sirupsen/logrus"
)

// CreateIPPool is handlers for create-ip-pool operation,
// useful for registering public IP pool that we have
type CreateIPPool struct {
	Handlers
	log logrus.FieldLogger
}

// NewCreateNodeHandler returns CreateIPPool handler
func NewCreateIPPoolHandler(handlers Handlers) *CreateIPPool {
	return &CreateIPPool{
		Handlers: handlers,
		log:      handlers.Log,
	}
}

// Handle is handler for create-ip-pool operation
func (h *CreateIPPool) Handle(params networks.CreateIPPoolParams) middleware.Responder {
	h.log.Infof("handling create ip pool request: %v", params)

	_, ipnet, err := net.ParseCIDR(params.IPPool.Cidr)
	if err != nil {
		h.log.Errorf("error creating ip pool '%v': %s", params, err)
		return networks.NewCreateIPPoolDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: fmt.Sprintf("error creating ip pool %s", err),
		})
	}

	id, err := h.storeIPPoolResource(*ipnet, params)
	if err != nil {
		h.log.Errorf("error storing ip pool '%v': %s", params, err)
		return networks.NewCreateIPPoolDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: fmt.Sprintf("error creating ip pool %s", err),
		})
	}

	res := &models.IPPool{
		ID:              id,
		Cidr:            ipnet.String(),
		ReservedIPPools: params.IPPool.ReservedIPPools,
	}
	return networks.NewCreateIPPoolCreated().WithPayload(res)
}

func (h *CreateIPPool) storeIPPoolResource(ipnet net.IPNet, params networks.CreateIPPoolParams) (string, error) {

	var labels r.Labels
	for _, l := range params.IPPool.Labels {
		labels = append(labels, r.Label{
			Key:   l.Key,
			Value: l.Value,
		})
	}

	ipPool := r.NewResourceIPPool(ipnet).WithLabels(labels)
	ipPool.Subnet = ipnet.IP.String()
	ipPool.SubnetLen = ipnet.Mask.String()
	ipPool.IPPoolRangeStart = params.IPPool.IPPoolRangeStart
	ipPool.IPPoolRangeEnd = params.IPPool.IPPoolRangeEnd
	ipPool.AllocatedIPPools = params.IPPool.ReservedIPPools

	pool := h.Handlers.Network.AllocateIPPools(ipnet, 1)
	var availableIPPools []string
	for _, ip := range pool {
		if containIP(ipPool.AllocatedIPPools, ip) {
			continue
		}
		availableIPPools = append(availableIPPools, ip)
	}
	ipPool.AvailableIPPools = availableIPPools

	_, err := h.Handlers.Store.AddFromPath(config.KeyPathNetworkIPPool+ipPool.Name, ipPool)
	if err != nil {
		h.log.Errorf("error storing new ip pool resource '%v': %s", ipPool, err)
		return "", fmt.Errorf("error storing new ip pool resource %s", err)
	}
	return ipPool.ID, nil
}

func containIP(ips []string, ip string) bool {
	for _, a := range ips {
		if a == ip {
			return true
		}
	}
	return false
}
