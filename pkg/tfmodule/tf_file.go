package tfmodule

import (
	"bytes"
	"html/template"
	"io"
	"os"
	"strings"
)

// Represents a terraform file as TerraformFile object
type TerraformFile struct {
	// TODO: the Source and Destination path would be absolute
	// TODO: validation for struct values
	FileType    string `json:"filetype"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

func (tfFile TerraformFile) ParseTerraformFile(data map[string]string) string {
	// TODO: check if source path exists -> raise error
	tfFileTemplate := template.New(tfFile.FileType)
	tfFileTemplate.Parse(tfFile.Source)

	var buf bytes.Buffer
	tfFileTemplate.Execute(&buf, data)
	return buf.String()
}

func (tfFile TerraformFile) WriteTerraformFile(data string) {
	// TODO: check if destination path exists -> create if not

	fo, _ := os.Create(tfFile.Destination)
	defer fo.Close()
	io.Copy(fo, strings.NewReader(data))
}
