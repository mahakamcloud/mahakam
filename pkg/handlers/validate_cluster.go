package handlers

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/clusters"
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/kube"
	"github.com/mahakamcloud/mahakam/pkg/resource_store/resource"
	"github.com/mahakamcloud/mahakam/pkg/validation"
	"github.com/sirupsen/logrus"
)

type ValidateCluster struct {
	Handlers
	log logrus.FieldLogger
}

func NewValidateClusterHandler(handlers Handlers) *ValidateCluster {
	return &ValidateCluster{
		Handlers: handlers,
		log:      handlers.Log,
	}
}

func (v *ValidateCluster) Handle(params clusters.ValidateClusterParams) middleware.Responder {
	v.log.Infof("handling validate cluster request: %v", params)

	clusterName := swag.StringValue(params.Body.Name)

	kubeconfig := resource.NewResourceCluster(clusterName).BuildKey()
	_, kubeclient, err := kube.GetKubeClient(config.HelmDefaultKubecontext, kubeconfig)
	if err != nil {
		v.log.Errorf("error getting kubernetes client for %s: %s", clusterName, err)
		return clusters.NewValidateClusterDefault(http.StatusInternalServerError).WithPayload(&models.Error{
			Message: "cannot initialize kubernetes client",
		})
	}

	podFailures, err := validation.ValidatePods(kubeclient)
	if err != nil {
		v.log.Errorf("error running cluster validation for %s: %s", clusterName, err)
		return clusters.NewValidateClusterDefault(http.StatusInternalServerError).WithPayload(&models.Error{
			Message: "cannot run cluster validation",
		})
	}

	var failures []string
	for _, pf := range podFailures {
		failures = append(failures, pf.Message)
	}

	res := &models.Cluster{
		Name:     params.Body.Name,
		Owner:    params.Body.Owner,
		Failures: failures,
	}
	return clusters.NewValidateClusterCreated().WithPayload(res)
}
