package helmcontroller

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mahakamcloud/mahakam/pkg/utils"
	"k8s.io/helm/pkg/downloader"
	"k8s.io/helm/pkg/getter"
	"k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/helmpath"
)

func newChartDownloader() downloader.ChartDownloader {
	helmHome := helmpath.Home(os.Getenv("HOME") + "/.helm")
	dl := downloader.ChartDownloader{
		HelmHome: helmHome,
		Out:      os.Stdout,
		Getters: getter.All(environment.EnvSettings{
			Home: helmHome,
		}),
	}
	return dl
}

func (hc *HelmController) GetRemoteChart() (string, error) {
	dl := newChartDownloader()

	chartCacheDir := filepath.Join("/tmp", hc.ReleaseName, Hash(hc.ChartPath))
	if err := os.MkdirAll(chartCacheDir, 0700); err != nil {
		return "", fmt.Errorf("cannot create work directory '%s'", chartCacheDir)
	}

	// Assume helm homdir is preconfigured, the cachedir must include
	// repository.yaml otherwise tiller will complain, so copy from existing one
	_, err := utils.CopyFile(string(dl.HelmHome)+"/repository/repositories.yaml", chartCacheDir+"/repositories.yaml")
	if err != nil {
		return "", err
	}

	// TODO(giri): Always download latest version for now,
	// must store our managed version somewhere and retrieve
	filepath, _, err := dl.DownloadTo(hc.ChartPath, "", chartCacheDir)
	if err != nil {
		return "", fmt.Errorf("failed to download '%s': %s", hc.ChartPath, err)
	}

	return filepath, nil
}

// Hash generates base64 encoded string
func Hash(s string) string {
	return base64.URLEncoding.EncodeToString([]byte(s))
}
