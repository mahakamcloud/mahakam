package handlers

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/clusters"
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/kube"
	"github.com/mahakamcloud/mahakam/pkg/utils"
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
	owner := params.Body.Owner
	kubeconfig := utils.GenerateKubeconfigPath(config.MahakamMultiKubeconfigPath, owner, clusterName)

	_, kubeclient, err := kube.GetKubeClient(config.HelmDefaultKubecontext, kubeconfig)
	if err != nil {
		v.log.Warnf("error getting kubernetes client for %s: %s", clusterName, err)
		return clusters.NewValidateClusterDefault(http.StatusInternalServerError).WithPayload(&models.Error{
			Message: "cannot initialize kubernetes client",
		})
	}

	nodeFailures, err := validation.ValidateNodes(kubeclient)
	if err != nil {
		v.log.Errorf("error validating cluster nodes for %s: %s", clusterName, err)
		return clusters.NewValidateClusterDefault(http.StatusInternalServerError).WithPayload(&models.Error{
			Message: "fail validating cluster nodes",
		})
	}

	componentFailures, err := validation.ValidateComponents(kubeclient)
	if err != nil {
		v.log.Errorf("error validating cluster components for %s: %s", clusterName, err)
		return clusters.NewValidateClusterDefault(http.StatusInternalServerError).WithPayload(&models.Error{
			Message: "fail validating cluster components",
		})
	}

	podFailures, err := validation.ValidatePods(kubeclient)
	if err != nil {
		v.log.Errorf("error validating kube-system pods for %s: %s", clusterName, err)
		return clusters.NewValidateClusterDefault(http.StatusInternalServerError).WithPayload(&models.Error{
			Message: "fail validating kube-system pods",
		})
	}

	res := validationResult(clusterName, owner, nodeFailures, componentFailures, podFailures)
	return clusters.NewValidateClusterCreated().WithPayload(res)
}

func validationResult(clusterName, owner string, nodeFailures, componentFailures, podFailures []*validation.ValidationError) *models.Cluster {
	var nfailures, cfailures, pfailures []string
	for _, nf := range nodeFailures {
		nfailures = append(nfailures, nf.Message)
	}
	for _, cf := range componentFailures {
		cfailures = append(cfailures, cf.Message)
	}
	for _, pf := range podFailures {
		pfailures = append(pfailures, pf.Message)
	}

	res := &models.Cluster{
		Name:              swag.String(clusterName),
		Owner:             owner,
		NodeFailures:      nfailures,
		ComponentFailures: cfailures,
		PodFailures:       pfailures,
	}
	return res
}
