package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/clusters"
	"github.com/sirupsen/logrus"
)

type ValidateCluster struct {
	Handlers
	log logrus.FieldLogger
}

func NewValidateClusterHandler(handlers Handlers) *ValidateCluster {
	return &ValidateCluster{
		Handlers: handlers,
		log:      handlers.Log,
	}
}

func (h *ValidateCluster) Handle(params clusters.ValidateClusterParams) middleware.Responder {
	h.log.Infof("handling validate cluster request: %v", params)

	// TODO(giri): call proper cluster validation and return result
	res := &models.Cluster{
		Failures: []string{},
	}
	return clusters.NewValidateClusterCreated().WithPayload(res)
}
