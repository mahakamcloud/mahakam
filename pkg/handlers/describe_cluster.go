package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/clusters"
	rs "github.com/mahakamcloud/mahakam/pkg/resource_store"
)

// DescribeCluster is handler for describe-cluster operation
type DescribeCluster struct {
	Handlers
}

// Handle is handler for describe-cluster operation
func (h *DescribeCluster) Handle(params clusters.DescribeClustersParams) middleware.Responder {

	c := rs.NewResourceCluster(swag.StringValue(params.Name))

	err := h.Handlers.Store.Get(c)
	if err != nil {
		return clusters.NewDescribeClustersDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: swag.String(err.Error()),
		})
	}

	return clusters.NewDescribeClustersOK().WithPayload(&models.Cluster{
		Name:        swag.String(c.Name),
		Owner:       c.Owner,
		ClusterPlan: string(c.Plan),
		NumNodes:    int64(c.NumNodes),
		Status:      string(c.Status),
	})
}
