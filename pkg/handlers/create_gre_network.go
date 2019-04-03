package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	grenet "github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/gre_networks"
	"github.com/mahakamcloud/mahakam/pkg/service"
	"github.com/sirupsen/logrus"
)

// CreateGreNetwork is handlers for register-bare-metal-host operation
type CreateGreNetwork struct {
	Handlers
	log logrus.FieldLogger
}

// NewCreateGreNetworkHandler registers new bare metal host
func NewCreateGreNetwork(handlers Handlers) *CreateGreNetwork {
	return &CreateGreNetwork{
		Handlers: handlers,
		log:      handlers.Log,
	}
}

func (g *CreateGreNetwork) Handle(params grenet.CreateGreNetworkParams) middleware.Responder {
	g.log.Infof("handling creation of Gre network request: %v", params)

	s, err := service.NewGreNetworkService()
	if err != nil {
		return g.handleError(err)
	}

	err = s.CreateGreNetwork(params.Body)
	if err != nil {
		return g.handleError(err)
	}

	return grenet.NewCreateGreNetworkCreated()
}

func (g *CreateGreNetwork) handleError(err error) middleware.Responder {
	g.log.Errorf("error creating gre network: %s", err)
	return grenet.NewCreateGreNetworkDefault(http.StatusInternalServerError).WithPayload(&models.Error{
		Message: fmt.Sprintf("failed creating the gre network"),
	})
}
