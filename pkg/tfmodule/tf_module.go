package tfmodule

import "fmt"

// TerraformModule which defines the files of a module
type TerraformModule struct {
	Name     string          `json:"name"`
	BasePath string          `json:"basepath"`
	Files    []TerraformFile `json:"files"`
}

// define or update a TerraformModule
func (tfModule TerraformModule) updateModule(filetype string, source string, destination string) {
	tfFile := TerraformFile{
		filetype,
		source,
		destination,
	}
	fmt.Print(tfFile)
	// if success store module in Module
}

func (tfModule TerraformModule) generateModule(data map[string]string) {
	for _, tfFile := range tfModule.Files {
		parsedFile := tfFile.parseTerraformFile(data)
		tfFile.writeTerraformFile(parsedFile)
	}
}

func (tfModule TerraformModule) executeModule() {
	// run cmd for this module
}
