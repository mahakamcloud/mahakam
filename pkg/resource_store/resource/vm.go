package resource

// ResourceTerraformBackend represents stored resource with terraform backend kind
type ResourceVM struct {
	Name                  string
	IPAddress             string
	MacAddress            string
	Host                  string
	ImageSourcePath       string
	DNSDhcpServerUsername string
	UserData              string
}

// NewResourceTerraformBackend creates new resource cluster
func NewResourceVM(name string) *ResourceVM {
	return &ResourceVM{}
}
