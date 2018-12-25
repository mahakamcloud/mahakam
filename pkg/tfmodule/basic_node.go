package tfmodule

import (
	"github.com/mahakamcloud/mahakam/pkg/tfmodule/templates"
)

func CreateNode(name, destdir string, data map[string]string) error {
	var basicNode = &TerraformProvisioner{
		Name:    name,
		DestDir: destdir,
		Files: []TerraformFile{
			TerraformFile{"backend", templates.Backend, destdir, "backend.tf"},
			TerraformFile{"data", templates.Data, destdir, "data.tf"},
			TerraformFile{"main", templates.MainFile, destdir, "main.tf"},
			TerraformFile{"tfvars", templates.TFVars, destdir, "terraform.tfvars"},
			TerraformFile{"vars", templates.Vars, destdir, "vars.tf"},
			TerraformFile{"foo", templates.CloudInit, destdir + "/templates/", "user_data.tpl"},
		},
	}

	basicNode.GenerateProvisionerFiles(data)
	err := basicNode.ExecuteProvisioner()
	if err != nil {
		return err
	}
	return nil
}
