package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	bmhost "github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/bare_metal_hosts"
	"github.com/mahakamcloud/mahakam/pkg/model"
	"github.com/sirupsen/logrus"
)

// RegisterBareMetalHost is handlers for register-bare-metal-host operation
type RegisterBareMetalHost struct {
	Handlers
	log logrus.FieldLogger
}

// NewRegisterBareMetalHostHandler registers new bare metal host
func NewRegisterBareMetalHostHandler(handlers Handlers) *RegisterBareMetalHost {
	return &RegisterBareMetalHost{
		Handlers: handlers,
		log:      handlers.Log,
	}
}

// Handle is handler for create-cluster operation
func (h *RegisterBareMetalHost) Handle(params bmhost.RegisterBareMetalHostParams) middleware.Responder {
	h.log.Infof("handling register bare metal host request: %v", params)

	bm := &model.BareMetalHostWrapper{}
	res := bm.BuildWithModel(params.Body)

	_, err := h.Handlers.Store.AddV1(bm)
	if err != nil {
		h.log.Errorf("error registering bare metal host: %s", err)
		return bmhost.NewRegisterBareMetalHostDefault(http.StatusInternalServerError).WithPayload(&models.Error{
			Message: fmt.Sprintf("failed registering bare metal host"),
		})
	}

	return bmhost.NewRegisterBareMetalHostCreated().WithPayload(res)
}
