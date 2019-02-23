package handlers

import (
	"fmt"
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/nodes"
	r "github.com/mahakamcloud/mahakam/pkg/resource_store/resource"
	"github.com/sirupsen/logrus"
)

// CreateNode is handlers for create-node operation
type CreateNode struct {
	Handlers
	log logrus.FieldLogger
}

// NewCreateNodeHandler returns CreateNode handler
func NewCreateNodeHandler(handlers Handlers) *CreateNode {
	return &CreateNode{
		Handlers: handlers,
		log:      handlers.Log,
	}
}

// Handle is handler for create-node operation
func (h *CreateNode) Handle(params nodes.CreateNodeParams) middleware.Responder {
	h.log.Infof("handling create node request: %v", params)

	id, err := h.storeNodeResource(params)
	if err != nil {
		h.log.Errorf("error storing node '%v': %s", params, err)
		return nodes.NewCreateNodeDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: fmt.Sprintf("error creating node %s", err),
		})
	}

	numericID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		h.log.Errorf("error parsing id returned from store '%v': %s", params, err)
		return nodes.NewCreateNodeDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: fmt.Sprintf("error creating ip pool %s", err),
		})
	}

	res := &models.Node{
		ID:   numericID,
		Name: params.Body.Name,
	}

	return nodes.NewCreateNodeCreated().WithPayload(res)
}

func (h *CreateNode) storeNodeResource(params nodes.CreateNodeParams) (string, error) {

	var labels r.Labels
	for _, l := range params.Body.Labels {
		labels = append(labels, r.Label{
			Key:   l.Key,
			Value: l.Value,
		})
	}

	node := r.NewResourceNode(*params.Body.Name).WithLabels(labels)

	nodeID, err := h.Handlers.Store.Add(node)
	if err != nil {
		h.log.Errorf("error storing new ip pool resource '%v': %s", node, err)
		return "", fmt.Errorf("error storing new ip pool resource %s", err)
	}

	return nodeID, nil
}
