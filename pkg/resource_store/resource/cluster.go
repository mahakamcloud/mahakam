package resource

import (
	"github.com/mahakamcloud/mahakam/pkg/config"
)

// Cluster represents stored resource with cluster kind
type Cluster struct {
	BaseResource
	NodeSize       string
	NumNodes       int
	NetworkName    string
	KubeconfigPath string
}

// NewCluster creates new resource cluster
func NewCluster(name string) *Cluster {
	return &Cluster{
		BaseResource: BaseResource{
			Name:  name,
			Kind:  string(KindCluster),
			Owner: config.ResourceOwnerGojek,
		},
	}
}

// ClusterList represents a group of Clusters
type ClusterList struct {
	Items []*Cluster
}

// Resource returns a empty Cluster
func (l *ClusterList) Resource() Resource {
	return &Cluster{}
}

// WithItems returns list of Cluster
func (l *ClusterList) WithItems(items []Resource) {
	for _, i := range items {
		cluster := i.(*Cluster)
		l.Items = append(l.Items, cluster)
	}
}
