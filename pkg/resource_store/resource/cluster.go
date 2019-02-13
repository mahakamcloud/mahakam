package resource

import (
	"strconv"

	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/utils"
)

const (
	// ClusterSizeExtraSmall is represents cluster with extra small size nodes
	ClusterSizeExtraSmall = "EXTRASMALL"
	// ClusterSizeSmall is represents cluster with small size nodes
	ClusterSizeSmall = "SMALL"
	// ClusterSizeMedium is represents cluster with medium size nodes
	ClusterSizeMedium = "MEDIUM"
	// ClusterSizeLarge is represents cluster with large size nodes
	ClusterSizeLarge = "LARGE"
	// ClusterSizeExtraLarge is represents cluster with extra large size nodes
	ClusterSizeExtraLarge = "EXTRALARGE"
	// ClusterSizeDefault is represents cluster with default size nodes
	ClusterSizeDefault = "DEFAULT"
)

// clusterNodeSizes represents available node size configurations.
var clusterNodeSizes = map[string]map[string]string{
	ClusterSizeExtraSmall: map[string]string{"cpu": "2", "ram": "4GB"},
	ClusterSizeSmall:      map[string]string{"cpu": "4", "ram": "16GB"},
	ClusterSizeMedium:     map[string]string{"cpu": "8", "ram": "32GB"},
	ClusterSizeLarge:      map[string]string{"cpu": "16", "ram": "64GB"},
	ClusterSizeExtraLarge: map[string]string{"cpu": "32", "ram": "128GB"},
	ClusterSizeDefault:    map[string]string{"cpu": "2", "ram": "4GB"},
}

// GetClusterNodeCPUs returns number of CPUs for a cluster node
func GetClusterNodeCPUs(size string) (int32, error) {
	size = ClusterSizeDefault // TODO: remove
	cpuInString := clusterNodeSizes[size]["cpu"]
	cpuInInt, err := strconv.ParseInt(cpuInString, 10, 32)

	if err != nil {
		return 0, err
	}

	return int32(cpuInInt), nil
}

// GetClusterNodeMemory returns memory for a cluster node in bytes from default GB representation
func GetClusterNodeMemory(size string) (int64, error) {
	size = ClusterSizeDefault
	memoryInGB := clusterNodeSizes[size]["ram"]
	memoryInBytes, err := utils.ToBytes(memoryInGB)
	if err != nil {
		return 0, err
	}

	return int64(memoryInBytes), nil
}

// ResourceCluster represents stored resource with cluster kind
type ResourceCluster struct {
	BaseResource
	NodeSize       string
	NumNodes       int
	NetworkName    string
	KubeconfigPath string
}

// NewResourceCluster creates new resource cluster
func NewResourceCluster(name string) *ResourceCluster {
	return &ResourceCluster{
		BaseResource: BaseResource{
			Name:  name,
			Kind:  string(KindCluster),
			Owner: config.ResourceOwnerGojek,
		},
	}
}

// ResourceClusterList represents a group of Clusters
type ResourceClusterList struct {
	Items []*ResourceCluster
}

// Resource returns a empty ResourceCluster
func (l *ResourceClusterList) Resource() Resource {
	return &ResourceCluster{}
}

// WithItems returns list of ResourceCluster
func (l *ResourceClusterList) WithItems(items []Resource) {
	for _, i := range items {
		cluster := i.(*ResourceCluster)
		l.Items = append(l.Items, cluster)
	}
}
