package writers

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

type TerraformWriter struct {
	Data          string
	DestDirectory string
	DestFile      string
}

// WriteFile func to write the data to destination address
func (tfWriter TerraformWriter) WriteFile() {
	if _, err := os.Stat(tfWriter.DestDirectory); os.IsNotExist(err) {
		os.MkdirAll(tfWriter.DestDirectory, os.ModePerm)
	}

	fo, _ := os.Create(filepath.Join(tfWriter.DestDirectory, tfWriter.DestFile))
	defer fo.Close()
	io.Copy(fo, strings.NewReader(tfWriter.Data))
}
