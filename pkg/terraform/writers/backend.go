package writers

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	r "github.com/mahakamcloud/mahakam/pkg/resource_store/resource"
	"github.com/mahakamcloud/mahakam/pkg/terraform/parsers"
	"github.com/mahakamcloud/mahakam/pkg/terraform/templates"
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

	var data = map[string]string{
		"Bucket": "tf-mahakam",
		"Key":    "gofinance-k8s/control-plane/terraform.tfstate",
		"Region": "ap-southeast-1",
	}
	terraformResource := r.NewResourceTerraform("backend.tf", data)

	backendParser := parsers.TerraformParser{
		"backend",
		templates.Backend,
		terraformResource.GetName(),
		terraformResource.GetAttributes(),
	}
	bakcendTf := backendParser.ParseTemplate()

	fo, _ := os.Create("/tmp/mahakam/terraform/backend.tf")
	defer fo.Close()
	io.Copy(fo, strings.NewReader(bakcendTf))
}
