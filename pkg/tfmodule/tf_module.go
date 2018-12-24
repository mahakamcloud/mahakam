package tfmodule

import (
	"fmt"
	"path/filepath"
)

// TerraformProvisioner which defines the files of a module
type TerraformProvisioner struct {
	Name    string          `json:"name"`
	DestDir string          `json:"destdir"`
	Files   []TerraformFile `json:"files"`
}

// define or update a UpdateProvisionerFile
func (tfProvisioner *TerraformProvisioner) UpdateProvisionerFile(filetype string, source string, destfile string) {
	destdir := filepath.Join(tfProvisioner.DestDir, tfProvisioner.Name)
	if filetype == "cloud-init" {
		destdir = filepath.Join(tfProvisioner.DestDir, tfProvisioner.Name, "templates")
	}

	tfFile := TerraformFile{
		filetype,
		source,
		destdir,
		destfile,
	}
	tfProvisioner.Files = append(tfProvisioner.Files, tfFile)
}

func (tfProvisioner TerraformProvisioner) GenerateProvisionerFiles(data map[string]string) {
	for _, tfFile := range tfProvisioner.Files {
		parsedFile := tfFile.ParseTerraformFile(data)
		tfFile.WriteTerraformFile(parsedFile)
	}
}

func (tfProvisioner *TerraformProvisioner) ExecuteProvisioner() {
	t := New()

	tfModuleDestDir := filepath.Join(tfProvisioner.DestDir, tfProvisioner.Name)
	tfVarsFile := filepath.Join(tfProvisioner.DestDir, tfProvisioner.Name, "terraform.tfvars")
	planOptions := "-var-file=" + tfVarsFile
	applyOptions := "-var-file=" + tfVarsFile + " -auto-approve"

	initOutput, initErr := t.Init(tfModuleDestDir)
	fmt.Println(initOutput)
	fmt.Println(initErr)

	planOutput, planErr := t.Plan(planOptions, tfModuleDestDir)
	fmt.Println(planOutput)
	fmt.Println(planErr)

	applyOutput, applyErr := t.Apply(applyOptions, tfModuleDestDir)
	fmt.Println(applyOutput)
	fmt.Println(applyErr)
}
