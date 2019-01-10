package v1

import (
	"github.com/mahakamcloud/mahakam/pkg/api/v1/client/clusters"
)

type ClusterAPI interface {
	CreateCluster(*clusters.CreateClusterParams) (*clusters.CreateClusterCreated, error)
	DescribeClusters(*clusters.DescribeClustersParams) (*clusters.DescribeClustersOK, error)
	GetClusters(*clusters.GetClustersParams) (*clusters.GetClustersOK, error)
}
