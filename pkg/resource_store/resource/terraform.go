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
func NewResourceTerraform(name string, attributes map[string]string) *ResourceTerraform {
	return &ResourceTerraform{
		BaseResource: BaseResource{
			Name:  name,
			Kind:  string(KindTerraformBackend),
			Owner: config.ResourceOwnerGojek,
		},
		Bucket: attributes["Bucket"],
		Key:    attributes["Key"],
		Region: attributes["Region"],
	}
}

func (r ResourceTerraform) GetName() string {
	return r.BaseResource.Name
}

func (r ResourceTerraform) GetAttributes() map[string]string {
	return map[string]string{
		"Bucket": r.Bucket,
		"Key":    r.Key,
		"Region": r.Region,
	}
}
