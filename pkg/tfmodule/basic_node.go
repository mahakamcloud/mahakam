package tfmodule

import (
	"github.com/mahakamcloud/mahakam/pkg/tfmodule/templates"
)

func CreateNode(data map[string]string) {
	var name = "mahakam-spike-01"
	var destdir = "/tmp/mahakam/terraform/mahakam-cluster-01"

	var basicNode = TerraformProvisioner{
		Name:    name,
		DestDir: destdir,
		Files: []TerraformFile{
			TerraformFile{"backend", templates.Backend, destdir, "backend.tf"},
			TerraformFile{"data", templates.UserData, destdir, "data.tf"},
			TerraformFile{"main", templates.MainFile, destdir, "main.tf"},
			TerraformFile{"tfvars", templates.TFVars, destdir, "terraform.tfvars"},
			TerraformFile{"vars", templates.Vars, destdir, "vars.tf"},
			TerraformFile{"cloud-init", templates.CloudInit, destdir, "cloud-init.tpl"},
		},
	}

	basicNode.GenerateProvisionerFiles(data)
	basicNode.ExecuteProvisioner()
}
