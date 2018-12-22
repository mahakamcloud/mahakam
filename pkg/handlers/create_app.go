package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/apps"
)

// CreateApp is handlers for create-app operatin
type CreateApp struct {
	Handlers
}

// Handle is handler for create-app operation
func (h *CreateApp) Handle(params apps.CreateAppParams) middleware.Responder {
	// TODO(giri): implement proper create app response by calling helmController
	return apps.NewCreateAppCreated().WithPayload(params.Body)
}
