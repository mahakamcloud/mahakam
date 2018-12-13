package terraform

import "fmt"

type Parser interface {
	parseTemplate()
}

type AbstractParser struct {
	Parser Parser
}

func (self *AbstractParser) Start() {
	self.Parser.parseTemplate()
}

type BackendParser struct {
	*AbstractParser
}

func (self *BackendParser) parseTemplate() {
	fmt.Println("Parsing Terraform backend.tf template")
}

type DataParser struct {
	*AbstractParser
}

func (self *DataParser) parseTemplate() {
	fmt.Println("Parsing Terraform data.tf template")
}
