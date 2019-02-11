package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/clusters"
	"github.com/sirupsen/logrus"
)

type GetCluster struct {
	Handlers
	log logrus.FieldLogger
}

// NewGetClusterHandler creates new CreateCluster object
func NewGetClusterHandler(handlers Handlers) *GetCluster {
	return &GetCluster{
		Handlers: handlers,
		log:      handlers.Log,
	}
}

// Handle is handler for create-cluster operation
func (h *GetCluster) Handle(params clusters.GetClustersParams) middleware.Responder {
	h.log.Infof("handling get cluster request: %v", params)

	owners, err := h.Handlers.Store.ListKeysFromPath("/cluster/")
	if err != nil {
		return clusters.NewGetClustersDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: err.Error(),
		})
	}

	return clusters.NewDescribeClustersOK().WithPayload([]*models.Cluster{})
}
