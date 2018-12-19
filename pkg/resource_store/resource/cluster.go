package resource

import "github.com/mahakamcloud/mahakam/pkg/config"

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
	Plan     ClusterPlan
	NumNodes int
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
