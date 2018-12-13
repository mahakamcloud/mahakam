package writers

import "fmt"

type DataWriter struct {
	*AbstractWriter
}

func (dw *DataWriter) writeFile() {
	fmt.Println("Writing Terraform data.tf template")
}
