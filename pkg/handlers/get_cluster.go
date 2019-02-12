package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/clusters"
	"github.com/mahakamcloud/mahakam/pkg/config"

	// "github.com/mahakamcloud/mahakam/pkg/resource_store/resource_store"
	"github.com/sirupsen/logrus"
)

type GetCluster struct {
	Handlers
	KubernetesConfig config.KubernetesConfig
	log              logrus.FieldLogger
}

// NewGetClusterHandler creates new CreateCluster object
func NewGetClusterHandler(handlers Handlers) *GetCluster {
	return &GetCluster{
		Handlers:         handlers,
		KubernetesConfig: handlers.AppConfig.KubernetesConfig,
		log:              handlers.Log,
	}
}

// Handle is handler for create-cluster operation
func (h *GetCluster) Handle(params clusters.GetClustersParams) middleware.Responder {
	h.log.Infof("handling get cluster request: %v", params)

	// clusterKind := resource.ResourceKind("cluster")
	// var clusters interface{}
	// err := resource_store.List(params.owner, clusterKind, clusters)
	// if err != nil {
	// 	return clusters.NewGetClustersDefault(405).WithPayload(&models.Error{
	// 		Code:    405,
	// 		Message: fmt.Sprintf("error getting cluster for owner %s", err),
	// 	})
	// }
	res := []*models.Cluster{}
	return clusters.NewGetClustersOK().WithPayload(res)
}
