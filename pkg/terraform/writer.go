package terraform

import "fmt"

type Writer interface {
	writeFile()
}

type AbstractWriter struct {
	Writer Writer
}

func (self *AbstractWriter) Start() {
	self.Writer.writeFile()
}

type BackendWriter struct {
	*AbstractWriter
}

func (self *BackendWriter) writeFile() {
	fmt.Println("Writing Terraform backend.tf template")
}

type DataWriter struct {
	*AbstractWriter
}

func (self *DataWriter) writeFile() {
	fmt.Println("Writing Terraform data.tf template")
}
