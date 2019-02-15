package handlers

import (
	"fmt"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/apps"
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/helmcontroller"
	"github.com/mahakamcloud/mahakam/pkg/helmcontroller/portforwarder"
	"github.com/mahakamcloud/mahakam/pkg/kube"
	r "github.com/mahakamcloud/mahakam/pkg/resource_store/resource"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	helm_env "k8s.io/helm/pkg/helm/environment"
)

// CreateApp is handlers for create-app operatin
type CreateApp struct {
	Handlers
	tillerTunnel *kube.Tunnel
	settings     helm_env.EnvSettings
	kubeclient   kubernetes.Interface
	chartValues  []string
	log          logrus.FieldLogger
}

// NewCreateAppHandler creates new CreateApp object
func NewCreateAppHandler(handlers Handlers) *CreateApp {
	return &CreateApp{
		Handlers: handlers,
		log:      handlers.Log,
	}
}

// Handle is handler for create-app operation
func (h *CreateApp) Handle(params apps.CreateAppParams) middleware.Responder {
	h.log.Infof("handling create app request: %v", params)

	b := params.Body
	cluster := r.NewCluster(b.ClusterName)
	err := h.Handlers.Store.Get(cluster)
	if err != nil {
		h.log.Errorf("error retrieving cluster info from kvstore '%v': %s\n", cluster, err)
		return apps.NewCreateAppDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: "cannot retrieve cluster info from kvstore",
		})
	}

	err = h.createHelmTillerTunnel(cluster.KubeconfigPath)
	if err != nil {
		h.log.Errorf("error creating tunnel to helm tiller: %v\n", err)
		return apps.NewCreateAppDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: "cannot create tunnel to helm tiller",
		})
	}

	releaseName := getReleaseName(b.Owner, swag.StringValue(b.Name))
	chartValuesFile := getChartValuesFile(b.Owner, b.ClusterName, swag.StringValue(b.Name))

	hc := helmcontroller.New(
		h.settings.TillerHost,
		b.ChartURL,
		getChartValues(config.HelmDefaultChartValuesDirectory+chartValuesFile),
		config.HelmDefaultNamespace,
		releaseName,
		config.HelmControllerWait,
		config.HelmControllerDefaultWaitTimeout,
		nil,
	)
	req := &models.App{
		Name:  b.Name,
		Owner: b.Owner,
	}

	// TODO(giri): run create app in separate routine and update app status
	// to ready when it's done
	err = hc.CreateApp(req)
	if err != nil {
		h.log.Errorf("error deploying app with helm chart '%v': %v\n", req, err)
		return apps.NewCreateAppDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: "cannot deploy application",
		})
	}

	serviceFQDN, err := h.getServiceFQDN(b.ChartURL, config.HelmDefaultNamespace,
		getReleaseName(b.Owner, swag.StringValue(b.Name)))
	if err != nil {
		h.log.Errorf("error getting service name with helm chart '%v': %v\n", req, err)
		return apps.NewCreateAppDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: "cannot retrieve service endpoint of application",
		})
	}

	res := &models.App{
		Name:        params.Body.Name,
		Owner:       params.Body.Owner,
		ClusterName: params.Body.ClusterName,
		ChartURL:    params.Body.ChartURL,
		ChartValues: params.Body.ChartValues,
		ServiceFQDN: serviceFQDN,
		Status:      string(r.StatusPending),
	}
	return apps.NewCreateAppCreated().WithPayload(res)
}

func (h *CreateApp) createHelmTillerTunnel(kubeconfig string) error {
	h.log.Debugf("create helm tiller tunnel with kubeconfig %s", kubeconfig)

	h.settings.TillerNamespace = config.HelmDefaultTillerNamespace
	h.settings.KubeConfig = kubeconfig
	h.settings.KubeContext = config.HelmDefaultKubecontext

	config, client, err := kube.GetKubeClient(h.settings.KubeContext, h.settings.KubeConfig)
	if err != nil {
		return fmt.Errorf("could not get kubernetes client for context %q: %s", h.settings.KubeContext, err)
	}
	h.kubeclient = client

	tillerTunnel, err := portforwarder.New(h.settings.TillerNamespace, client, config)
	if err != nil {
		return fmt.Errorf("could not create tiller tunnel: %s", err)
	}

	h.settings.TillerHost = fmt.Sprintf("127.0.0.1:%d", tillerTunnel.Local)

	h.log.Debugf("created helm tiller tunnel using local port %d", tillerTunnel.Local)
	return nil
}

func (h *CreateApp) getServiceFQDN(chartURL, namespace, releaseName string) (string, error) {
	h.log.Debugf("getting service fqdn for chart %s, namespace %s, release %s", chartURL, namespace, releaseName)

	app := strings.Split(chartURL, "/")
	if len(app) < 2 {
		return "", fmt.Errorf("invalid chart url %s", chartURL)
	}
	appName := app[1]

	serviceLabels := labels.Set{"app": appName, "release": releaseName, "heritage": "Tiller"}
	serviceName, err := kube.GetServiceName(h.kubeclient.CoreV1(), namespace, serviceLabels)
	if err != nil {
		return "", fmt.Errorf("could not retrieve endpoint service of deployed chart: %s", err)
	}

	serviceFQDN := fmt.Sprintf("%s.%s.svc.cluster.local", serviceName, namespace)
	return serviceFQDN, nil
}

func getChartValues(value string) []string {
	var valueArray []string
	if value != "" {
		valueArray = append(valueArray, value)
	}
	return valueArray
}

func getReleaseName(owner, appName string) string {
	return fmt.Sprintf("%s-%s", owner, appName)
}

func getChartValuesFile(owner, clusterName, appName string) string {
	return fmt.Sprintf("%s-%s-%s", owner, clusterName, appName)
}
