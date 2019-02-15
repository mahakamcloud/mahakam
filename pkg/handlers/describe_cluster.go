package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/clusters"
	r "github.com/mahakamcloud/mahakam/pkg/resource_store/resource"
)

// DescribeCluster is handler for describe-cluster operation
type DescribeCluster struct {
	Handlers
}

// Handle is handler for describe-cluster operation
func (h *DescribeCluster) Handle(params clusters.DescribeClustersParams) middleware.Responder {

	c := r.NewCluster(swag.StringValue(params.Name))

	err := h.Handlers.Store.Get(c)
	if err != nil {
		return clusters.NewDescribeClustersDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: err.Error(),
		})
	}

	return clusters.NewDescribeClustersOK().WithPayload(&models.Cluster{
		Name:     swag.String(c.Name),
		Owner:    c.Owner,
		NumNodes: int64(c.NumNodes),
		NodeSize: swag.String(c.NodeSize),
		Status:   string(c.Status),
	})
}
