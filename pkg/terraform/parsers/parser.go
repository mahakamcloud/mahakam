package parsers

import (
	"bytes"
	"text/template"
)

type TerraformTemplate struct {
	Name        string
	Source      string
	Destination string
	Data        map[string]string
}

type TerraformParser struct {
	parsers []TerraformTemplate
}

type Parser interface {
	ParseTemplate() string
}

func (self *TerraformTemplate) ParseTemplate() string {
	t := template.New(self.Name)
	t.Parse(self.Source)

	var buf bytes.Buffer
	t.Execute(&buf, &self.Data)
	return buf.String()
}
