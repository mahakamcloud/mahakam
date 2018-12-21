package tfmodule

// TerraformProvisioner which defines the files of a module
type TerraformProvisioner struct {
	Name    string          `json:"name"`
	DestDir string          `json:"destdir"`
	Files   []TerraformFile `json:"files"`
}

// define or update a UpdateProvisionerFile
func (tfProvisioner *TerraformProvisioner) UpdateProvisionerFile(filetype string, source string, destfile string) {
	tfFile := TerraformFile{
		filetype,
		source,
		tfProvisioner.DestDir + tfProvisioner.Name + "/",
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

func (tfProvisioner TerraformProvisioner) executeModule() {
	// run cmd for this module
}
