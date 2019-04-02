package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/bare_metal_hosts"
	"github.com/mahakamcloud/mahakam/pkg/service"
	"github.com/sirupsen/logrus"
)

type GetBareMetalHost struct {
	Handlers
	log logrus.FieldLogger
}

// NewGetBareMetalHostHandler creates new CreateCluster object
func NewGetBareMetalHostHandler(handlers Handlers) *GetBareMetalHost {
	return &GetBareMetalHost{
		Handlers: handlers,
		log:      handlers.Log,
	}
}

// Handle is handler for create-cluster operation
func (h *GetBareMetalHost) Handle(params bare_metal_hosts.GetBareMetalHostsParams) middleware.Responder {
	h.log.Infof("handling get baremetal host request: %v", params)

	service, err := service.NewBareMetalHostService()
	if err != nil {
		h.handleError(err)
	}

	hosts, err := service.GetAll()
	if err != nil {
		h.handleError(err)
	}
	return bare_metal_hosts.NewGetBareMetalHostsOK().WithPayload(hosts)
}

func (h *GetBareMetalHost) handleError(err error) middleware.Responder {
	h.log.Errorf("error getting bare-metal-host: %s", err)
	return bare_metal_hosts.NewGetBareMetalHostsDefault(http.StatusInternalServerError).WithPayload(&models.Error{
		Code:    http.StatusInternalServerError,
		Message: fmt.Sprintf("error getting bare-metal-hosts %s: ", err),
	})
}
