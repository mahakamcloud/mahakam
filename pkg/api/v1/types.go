package v1

import (
	"github.com/mahakamcloud/mahakam/pkg/api/v1/client/apps"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/client/clusters"
)

type ClusterAPI interface {
	CreateCluster(*clusters.CreateClusterParams) (*clusters.CreateClusterCreated, error)
	DescribeClusters(*clusters.DescribeClustersParams) (*clusters.DescribeClustersOK, error)
	GetClusters(*clusters.GetClustersParams) (*clusters.GetClustersOK, error)
}

type AppAPI interface {
	CreateApp(*apps.CreateAppParams) (*apps.CreateAppCreated, error)
	GetApps(*apps.GetAppsParams) (*apps.GetAppsOK, error)
	UploadAppValues(*apps.UploadAppValuesParams) (*apps.UploadAppValuesCreated, error)
}
