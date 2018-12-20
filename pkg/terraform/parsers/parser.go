package parsers

import (
	"bytes"
	"text/template"
)

type TerraformParser struct {
	Name        string
	Source      string
	Destination string
	Data        map[string]string
}

type Parser interface {
	ParseTemplate() string
}

func (self *TerraformParser) ParseTemplate() string {
	t := template.New(self.Name)
	t.Parse(self.Source)

	var buf bytes.Buffer
	t.Execute(&buf, &self.Data)
	return buf.String()
}
