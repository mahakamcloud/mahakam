package tfmodule

var tfModules = make(map[string]TerraformModule)

// Defining terraform file type to terraform name mapping
const (
	Backend     = "backend.tf"
	Main        = "main.tf"
	Data        = "data.tf"
	Vars        = "vars.tf"
	CloudConfig = "cloud_config.tpl"
)

func () defineModule () {
	privateVMModule = TerraformModule(
		"mahakam-spike-11",
		"/tmp/mahakam_cloud/mahakam/payportal/",
		[],
	)
}
