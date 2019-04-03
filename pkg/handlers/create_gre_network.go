package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	grenet "github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/gre_networks"
	"github.com/sirupsen/logrus"
)

// CreateGreNetwork is handlers for register-bare-metal-host operation
type CreateGreNetwork struct {
	Handlers
	log logrus.FieldLogger
}

// NewRegisterBareMetalHostHandler registers new bare metal host
func NewCreateGreNetwork(handlers Handlers) *CreateGreNetwork {
	return &CreateGreNetwork{
		Handlers: handlers,
		log:      handlers.Log,
	}
}

func (g *CreateGreNetwork) Handle(params grenet.GetGreNetworksParams) middleware.Responder {
	return nil
}
