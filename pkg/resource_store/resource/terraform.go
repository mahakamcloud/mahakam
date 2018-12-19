package resource

import "github.com/mahakamcloud/mahakam/pkg/config"

// ResourceTerraformBackend represents stored resource with terraform backend kind
type ResourceTerraform struct {
	BaseResource
	Bucket            string
	Key               string
	Region            string
	TerraformVersion  string
	LibvirtModulePath string
	LibvirtVersion    string
}

// NewResourceTerraform creates new resource cluster
func NewResourceTerraform(name string) *ResourceTerraform {
	return &ResourceTerraform{
		BaseResource: BaseResource{
			Name:  name,
			Kind:  string(KindTerraformBackend),
			Owner: config.ResourceOwnerGojek,
		},
	}
}
