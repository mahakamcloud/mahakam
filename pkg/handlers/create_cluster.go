package handlers

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/clusters"
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/provisioner"
	r "github.com/mahakamcloud/mahakam/pkg/resource_store/resource"
)

// CreateCluster is handlers for create-cluster operation
type CreateCluster struct {
	Handlers
}

// Handle is handler for create-cluster operation
func (h *CreateCluster) Handle(params clusters.CreateClusterParams) middleware.Responder {
	b := params.Body
	c := r.NewResourceCluster(swag.StringValue(b.Name))
	c.NumNodes = int(b.NumNodes)
	c.Status = r.StatusPending

	// TODO(giri): create cluster workflow should pull
	// /etc/kubernetes/admin.conf into this kubeconfig path
	c.KubeconfigPath = h.generateKubeconfigPath(c.Owner, c.Name)

	_, err := h.Handlers.Store.Add(c)
	if err != nil {
		return clusters.NewCreateClusterDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: swag.String(err.Error()),
		})
	}

	// TODO(giri/iqbal): run this provisioner from another routine,
	// must update resource status to creating and success accordingly
	err = provisioner.CreateCluster(b)
	if err != nil {
		return clusters.NewCreateClusterDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: swag.String(err.Error()),
		})
	}
	return clusters.NewCreateClusterCreated().WithPayload(b)
}

func (h *CreateCluster) generateKubeconfigPath(owner, clusterName string) string {
	return fmt.Sprintf(config.MahakamMultiKubeconfigPath + "/" + owner + "-" + clusterName + "-kubeconfig")
}
