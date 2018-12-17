package writers

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mahakamcloud/mahakam/pkg/terraform/parsers"
)

type BackendWriter struct {
	*AbstractWriter
}

func (bw *BackendWriter) writeFile() {
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		os.MkdirAll(basePath, os.ModePerm)
	}
	destinationPath := filepath.Join(basePath, "mahakam-test-cluster")
	os.MkdirAll(destinationPath, os.ModePerm)

	backendParser := &parsers.AbstractParser{&parsers.BackendParser{}}
	bakcendTf := backendParser.Parse()

	fo, _ := os.Create("/tmp/mahakam/terraform/backend.tf")
	defer fo.Close()
	io.Copy(fo, strings.NewReader(bakcendTf))
}
