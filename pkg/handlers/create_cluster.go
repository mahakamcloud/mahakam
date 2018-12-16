package handlers

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/clusters"
	"github.com/mahakamcloud/mahakam/pkg/provisioner"
	rs "github.com/mahakamcloud/mahakam/pkg/resource_store"
)

// CreateCluster is handlers for create-cluster operation
type CreateCluster struct {
	Handlers
}

// Handle is handler for create-cluster operation
func (h *CreateCluster) Handle(params clusters.CreateClusterParams) middleware.Responder {
	b := params.Body
	c := rs.NewResourceCluster(swag.StringValue(b.Name))
	c.NumNodes = int(b.NumNodes)
	c.Status = rs.StatusPending

	_, err := h.Handlers.Store.Add(c)
	if err != nil {
		fmt.Printf("Error storing: %s", err)
		return clusters.NewCreateClusterDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: swag.String(err.Error()),
		})
	}

	// TODO(giri/iqbal): run this provisioner from another routine,
	// must update resource status to creating and success accordingly
	err = provisioner.CreateCluster(b)
	if err != nil {
		fmt.Printf("Error creating: %s", err)
		return clusters.NewCreateClusterDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: swag.String(err.Error()),
		})
	}
	return clusters.NewCreateClusterCreated().WithPayload(b)
}
