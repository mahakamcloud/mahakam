package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/bare_metal_hosts"
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/model"
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

	bareMetalHostBuilderList := &model.BareMetalHostBuilderList{Items: []*model.BareMetalHost{}}

	err := h.Handlers.Store.ListV1(config.ResourceOwnerMahakam, model.KindBareMetalHost, bareMetalHostBuilderList)
	if err != nil {
		return bare_metal_hosts.NewGetBareMetalHostsDefault(http.StatusInternalServerError).WithPayload(&models.Error{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("error getting bare-metal-hosts %s: ", err),
		})
	}

	b := bareMetalHostBuilderList.GetBareMetalHosts()
	return bare_metal_hosts.NewGetBareMetalHostsOK().WithPayload(b)
}
