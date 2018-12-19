package writers

import (
	"io"
	"os"
	"path/filepath"
	"strings"
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

	// TODO(iqbal/himani): comment below snippet to enable running server
	// since undefined: parsers.AbstractParser
	// backendParser := &parsers.AbstractParser{&parsers.BackendParser{}}
	// bakcendTf := backendParser.Parse()
	bakcendTf := ""

	fo, _ := os.Create("/tmp/mahakam/terraform/backend.tf")
	defer fo.Close()
	io.Copy(fo, strings.NewReader(bakcendTf))
}
