package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	bmhost "github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/bare_metal_hosts"
	"github.com/mahakamcloud/mahakam/pkg/service"
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

	s, err := service.NewBareMetalHostService()
	if err != nil {
		return h.handleError(err)
	}

	err = s.RegisterBareMetalHost(params.Body)
	if err != nil {
		return h.handleError(err)
	}

	return bmhost.NewRegisterBareMetalHostCreated()
}

func (h *RegisterBareMetalHost) handleError(err error) middleware.Responder {
	h.log.Errorf("error registering bare metal host: %s", err)
	return bmhost.NewRegisterBareMetalHostDefault(http.StatusInternalServerError).WithPayload(&models.Error{
		Message: fmt.Sprintf("failed registering bare metal host"),
	})
}
