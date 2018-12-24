package helmcontroller

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-openapi/swag"
	yaml "gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"

	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/config"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/release"
)

// HelmController works with helm to deploy helm charts
// into Kubernetes cluster
type HelmController struct {
	// Helm is client that talks to tiller through GRPC
	Helm helm.Interface
	// ChartPath is path to dir where chart archive is located
	ChartPath string
	// ValueFile is filepath that contains custom values to override default chart values
	ValueFiles valueFiles
	// Namespace is default namespace to deploy resources into
	Namespace string
	// ReleaseName is prefix to helm release name
	ReleaseName string
	// Wait is whether or not to wait for resources to complete before stating successful release
	Wait bool
	// WaitTimeout is time in seconds to wait for resources to be created before stating successful release
	WaitTimeout int64
	// logger holds the log
	logger *log.Logger
}

// New will return a configured helm controller
func New(tillerEndpoint, chartPath string, valueFiles valueFiles,
	namespace, releaseName string, wait bool, waitTimeout int64, logger *log.Logger) *HelmController {

	if logger == nil {
		logger = log.New()
		logger.Out = os.Stdout
		logger.WithFields(log.Fields{
			"module": "helm_controller",
		})
	}

	if namespace == "" {
		namespace = config.HelmDefaultNamespace
	}

	return &HelmController{
		Helm:        helm.NewClient(helm.Host(tillerEndpoint)),
		ChartPath:   chartPath,
		ValueFiles:  valueFiles,
		Namespace:   namespace,
		ReleaseName: releaseName,
		Wait:        wait,
		WaitTimeout: waitTimeout,
		logger:      logger,
	}
}

// CreateApp will kick off helm install of given chart and values
func (hc *HelmController) CreateApp(app *models.App) error {
	hc.logger.Infof("create app %s", swag.StringValue(app.Name))
	if err := hc.installOrUpdate(app); err != nil {
		hc.logger.Errorf("failed to create app %s: %s", swag.StringValue(app.Name), err)
		return err
	}
	return nil
}

func (hc *HelmController) installOrUpdate(app *models.App) error {
	rawVals, err := hc.vals(hc.ValueFiles)
	if err != nil {
		return err
	}

	chartPath, err := hc.GetRemoteChart()
	if err != nil {
		return err
	}
	hc.ChartPath = chartPath

	releaseName := hc.releaseName(app)
	if hc.releaseExists(releaseName) {
		_, err := hc.Helm.UpdateRelease(
			releaseName,
			hc.ChartPath,
			helm.UpdateValueOverrides(rawVals),
			helm.UpgradeWait(hc.Wait),
			helm.UpgradeTimeout(hc.WaitTimeout),
		)
		return err
	}

	_, err = hc.Helm.InstallRelease(
		hc.ChartPath,
		hc.Namespace,
		helm.ReleaseName(releaseName),
		helm.ValueOverrides(rawVals),
		helm.InstallWait(hc.Wait),
		helm.InstallTimeout(hc.WaitTimeout),
	)
	return err
}

func (hc *HelmController) vals(valueFiles valueFiles) ([]byte, error) {
	currentMap := map[string]interface{}{}

	for _, filePath := range valueFiles {
		var bytes []byte
		var err error

		if strings.TrimSpace(filePath) == "-" {
			bytes, err = ioutil.ReadAll(os.Stdin)
		} else {
			bytes, err = readFile(filePath)
		}

		if err != nil {
			return []byte{}, err
		}

		if err := yaml.Unmarshal(bytes, &currentMap); err != nil {
			return []byte{}, fmt.Errorf("failed to parse %s: %s", filePath, err)
		}
	}

	return yaml.Marshal(currentMap)
}

func (hc *HelmController) releaseName(app *models.App) string {
	return fmt.Sprintf("%s-%s", app.Owner, swag.StringValue(app.Name))
}

func (hc *HelmController) releaseExists(releaseName string) bool {
	statuses := []release.Status_Code{
		release.Status_UNKNOWN,
		release.Status_DEPLOYED,
		release.Status_DELETED,
		release.Status_DELETING,
		release.Status_FAILED,
		release.Status_PENDING_INSTALL,
		release.Status_PENDING_UPGRADE,
		release.Status_PENDING_ROLLBACK,
	}

	rel, err := hc.Helm.ListReleases(
		helm.ReleaseListNamespace(hc.Namespace),
		helm.ReleaseListFilter(releaseName),
		helm.ReleaseListStatuses(statuses),
	)

	if err != nil || rel == nil {
		return false
	}

	for _, r := range rel.Releases {
		if r.GetName() == releaseName {
			return true
		}
	}

	return false
}
