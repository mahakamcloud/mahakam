package builder

import "github.com/mahakamcloud/mahakam/pkg/api/v1/models"

// Kind represents stored resource kind
type Kind string

const (
	KindCluster          Kind = "cluster"
	KindTerraform        Kind = "terraform"
	KindTerraformBackend Kind = "terraform backend"
	KindTask             Kind = "task"
	KindNetwork          Kind = "network"
	KindNode             Kind = "node"
	KindIPPool           Kind = "ippool"
)

// ResourceBuilder defines actions on ResourceBuilder
type ResourceBuilder interface {
	Build(name, kind, owner, role string) ResourceBuilder
	BuildChildKey(parentKey string, key string) string
	BuildKey(optKeys ...string) (string, error)
	BuildMetadata() ResourceBuilder
	GetLabels() []*models.BaseResourceLabelsItems0
	GetResourceID() string
	Update() ResourceBuilder
	UpdateRevision(index uint64) ResourceBuilder
	Validate() error
}
