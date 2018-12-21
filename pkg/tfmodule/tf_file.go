package tfmodule

import (
	"bytes"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Represents a terraform file as TerraformFile object
type TerraformFile struct {
	// TODO: the Source and Destination path would be absolute
	// TODO: validation for struct values
	FileType string `json:"filetype"`
	Source   string `json:"source"`
	DestDir  string `json:"destdir"`
	DestFile string `json:"destfile"`
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

	if _, err := os.Stat(tfFile.DestDir); os.IsNotExist(err) {
		os.MkdirAll(tfFile.DestDir, os.ModePerm)
	}

	fo, _ := os.Create(filepath.Join(tfFile.DestDir, tfFile.DestFile))
	defer fo.Close()
	io.Copy(fo, strings.NewReader(data))
}
