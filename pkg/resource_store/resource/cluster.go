package resource

import (
	"github.com/mahakamcloud/mahakam/pkg/config"
)

type ClusterPlan string

const (
	ClusterPlanSmall   ClusterPlan = "small"
	ClusterPlanMedium  ClusterPlan = "medium"
	ClusterPlanLarge   ClusterPlan = "large"
	ClusterPlanDefault ClusterPlan = ClusterPlanSmall
)

// ResourceCluster represents stored resource with cluster kind
type ResourceCluster struct {
	BaseResource
	Plan           ClusterPlan
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
		Plan: ClusterPlanDefault,
	}
}

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
