package tfmodule

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
)

// TerraformFile represents a terraform file as TerraformFile object
type TerraformFile struct {
	// TODO: the Source and Destination path would be absolute
	// TODO: validation for struct values
	FileType string `json:"filetype"`
	Source   string `json:"source"`
	DestDir  string `json:"destdir"`
	DestFile string `json:"destfile"`
}

// ParseTerraformFile parses terraform file and return string content
func (tfFile TerraformFile) ParseTerraformFile(data map[string]string) string {
	// TODO: check if source path exists -> raise error
	tfFileTemplate := template.New(tfFile.FileType)

	tfFileTemplate.Parse(tfFile.Source)

	var buf bytes.Buffer
	err := tfFileTemplate.Execute(&buf, data)

	if err != nil {
		log.Errorf("TF Template Execute error : %s", err)
	}

	return buf.String()
}

// WriteTerraformFile writes string content to files
func (tfFile TerraformFile) WriteTerraformFile(data string) {
	// TODO: check if destination path exists -> create if not
	destdir := tfFile.DestDir
	if tfFile.FileType == "cloud-init" {
		destdir = filepath.Join(tfFile.DestDir, "templates")
	}

	if _, err := os.Stat(destdir); os.IsNotExist(err) {
		os.MkdirAll(destdir, os.ModePerm)
	}

	fo, err := os.Create(filepath.Join(destdir, tfFile.DestFile))

	if err != nil {
		log.Errorf("Write error : %s", err)
	}

	defer fo.Close()
	_, err = io.Copy(fo, strings.NewReader(data))

	if err != nil {
		log.Errorf("Copy error : %s", err)
	}

}
