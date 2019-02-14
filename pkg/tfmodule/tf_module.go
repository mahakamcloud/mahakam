package tfmodule

import (
	"fmt"
	"path/filepath"
)

type Data struct {
	Output []byte
	Error  error
}

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

// GenerateProvisionerFiles writes terraform files to disk
func (tfProvisioner TerraformProvisioner) GenerateProvisionerFiles(data map[string]string) {
	for _, tfFile := range tfProvisioner.Files {
		parsedFile := tfFile.ParseTerraformFile(data)
		tfFile.WriteTerraformFile(parsedFile)
	}
}

// ExecuteProvisioner executes terraform init and terraform apply
func (tfProvisioner *TerraformProvisioner) ExecuteProvisioner() error {
	t := New()

	tfModuleDestDir := tfProvisioner.DestDir
	tfVarsFile := filepath.Join(tfProvisioner.DestDir, "terraform.tfvars")

	// TODO(giri/himani): pass proper thread safe logger instead of fmt.Println
	resInit, err := t.Init(tfModuleDestDir)
	fmt.Println(resInit)
	if err != nil {
		return fmt.Errorf("error initializing terraform: %s", err)
	}

	// TODO(giri/himani): pass proper thread safe logger instead of fmt.Println
	resApply, err := t.ApplyWithTFVars(tfVarsFile, tfModuleDestDir)
	fmt.Println(resApply)
	if err != nil {
		return fmt.Errorf("error applying terraform files: %s", err)
	}
	return nil
}
