package writers

import "fmt"

type DataWriter struct {
	*AbstractWriter
}

func (self *DataWriter) writeFile() {
	fmt.Println("Writing Terraform data.tf template")
}
