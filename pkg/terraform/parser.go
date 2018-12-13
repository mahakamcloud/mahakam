package terraform

import "fmt"

type Parser interface {
	parseTemplate()
}

type AbstractParser struct {
	Parser Parser
}

func (ap *AbstractParser) Start() {
	ap.Parser.parseTemplate()
}

type BackendParser struct {
	*AbstractParser
}

func (bp *BackendParser) parseTemplate() {
	fmt.Println("Parsing Terraform backend.tf template")
}

type DataParser struct {
	*AbstractParser
}

func (dp *DataParser) parseTemplate() {
	fmt.Println("Parsing Terraform data.tf template")
}
