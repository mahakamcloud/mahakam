package handlers

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/apps"
	"github.com/mahakamcloud/mahakam/pkg/kube"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

func configForContext(context string, kubeconfig string) (*rest.Config, error) {
	config, err := kube.GetConfig(context, kubeconfig).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get kubernetes config for context %q: %s", context, err)
	}
	return config, nil
}

func getKubeClient(context string, kubeconfig string) (*rest.Config, kubernetes.Interface, error) {
	config, err := configForContext(context, kubeconfig)
	if err != nil {
		return nil, nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get kubernetes client: %s", err)
	}
	return config, client, nil
}
