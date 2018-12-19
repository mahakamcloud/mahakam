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

	terraformResource := r.NewResourceTerraform("backend.tf")
	terraformResource.Bucket = "tf-mahakam"
	terraformResource.Key = "gofinance-k8s/control-plane/terraform.tfstate"
	terraformResource.Region = "ap-southeast-1"

	var data = map[string]string{
		"Bucket": terraformResource.Bucket,
		"Key":    terraformResource.Key,
		"Region": terraformResource.Region,
	}

	backendParser := parsers.TerraformTemplate{
		"backend",
		templates.Backend,
		"backend.tf",
		data,
	}
	bakcendTf := backendParser.ParseTemplate()

	fo, _ := os.Create("/tmp/mahakam/terraform/backend.tf")
	defer fo.Close()
	io.Copy(fo, strings.NewReader(bakcendTf))
}
