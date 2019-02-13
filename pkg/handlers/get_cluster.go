package handlers

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/clusters"
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/resource_store/resource"

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

	clusterKind := resource.ResourceKind("cluster")
	clusterList := &resource.ResourceClusterList{Items: []*resource.ResourceCluster{}}

	err := h.Handlers.Store.List(swag.StringValue(params.Owner), clusterKind, clusterList)
	if err != nil {
		return clusters.NewGetClustersDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: fmt.Sprintf("error getting cluster for owner %s", err),
		})
	}
	res := clusterList.Items

	c := []*models.Cluster{}
	for _, v := range res {
		clusterModel := &models.Cluster{
			Name:     &v.Name,
			NumNodes: int64(v.NumNodes),
			Owner:    v.Owner,
			Status:   string(v.Status),
		}

		c = append(c, clusterModel)
	}
	return clusters.NewGetClustersOK().WithPayload(c)
}
