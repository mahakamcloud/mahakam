package parsers

import (
	"bytes"
	"text/template"

	"github.com/mahakamcloud/mahakam/pkg/terraform/templates"
)

type Backend struct {
	Bucket string
	Key    string
	Region string
}

type BackendParser struct {
	*AbstractParser
}

func (self *BackendParser) parseTemplate() string {
	t := template.New("backend.tf")
	t.Parse(templates.Backend)

	terraformBackend := &Backend{
		"tf-mahakam", "mahakam-spike/terraform.tfstate", "ap-southeast-1",
	}

	var buf bytes.Buffer
	t.Execute(&buf, terraformBackend)
	return buf.String()
}
