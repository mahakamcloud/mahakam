package provisioner

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	rs "github.com/mahakamcloud/mahakam/pkg/resource_store/resource"
	"github.com/mahakamcloud/mahakam/pkg/terraform/parsers"
	"github.com/mahakamcloud/mahakam/pkg/terraform/templates"
	tw "github.com/mahakamcloud/mahakam/pkg/terraform/writers"
)

// CreateCluster creates cluster
func CreateCluster(cluster *models.Cluster) error {
	fmt.Println("Creating cluster...")

	var data = map[string]string{
		"Bucket": "tf-mahakam",
		"Key":    "gofinance-k8s/control-plane/terraform.tfstate",
		"Region": "ap-southeast-1",
	}
	terraformResource := rs.NewResourceTerraform("backend.tf", data)

	backendParser := parsers.TerraformParser{
		Name:        "backend",
		Source:      templates.Backend,
		Destination: terraformResource.GetName(),
		Data:        terraformResource.GetAttributes(),
	}
	backendTf := backendParser.ParseTemplate()

	tfWriter := tw.TerraformWriter{
		Data:          backendTf,
		DestDirectory: "/tmp/mahakam/terraform/",
		DestFile:      "backend.tf",
	}
	tfWriter.WriteFile()

	return nil
}

func copyTerraformFiles(src, dst string) error {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
	return nil
}
