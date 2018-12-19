package resource

import "github.com/mahakamcloud/mahakam/pkg/config"

// ResourceTerraformBackend represents stored resource with terraform backend kind
type ResourceTerraformBackend struct {
	BaseResource
	Bucket string
	Key    string
	Region string
}

// NewResourceTerraformBackend creates new resource cluster
func NewResourceTerraformBackend(name string) *ResourceTerraformBackend {
	return &ResourceTerraformBackend{
		BaseResource: BaseResource{
			Name:  name,
			Kind:  string(KindTerraformBackend),
			Owner: config.ResourceOwnerGojek,
		},
	}
}
