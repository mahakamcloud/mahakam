package writers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/mahakamcloud/mahakam/pkg/terraform/templates"
)

type Backend struct {
	Bucket string
	Key    string
	Region string
}

type BackendWriter struct {
	*AbstractWriter
}

func (bw *BackendWriter) writeFile() {
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		os.MkdirAll(basePath, os.ModePerm)
	}

	destinationPath := filepath.Join(basePath, "mahakam-test-cluster")
	os.MkdirAll(destinationPath, os.ModePerm)

	// copyTerraformFiles(templatesPath, destinationPath)
	t := template.New("backend.tf")
	t.Parse(templates.Backend)

	terraformBackend := &Backend{
		"tf-mahakam", "mahakam-spike/terraform.tfstate", "ap-southeast-1",
	}

	var buf bytes.Buffer
	t.Execute(&buf, terraformBackend)
	fmt.Println(buf.String())

	err := ioutil.WriteFile("/tmp/mahakam/terraform/backend.tf", buf.Bytes(), 0644)
	if err != nil {
		fmt.Println(err)
	}
}
