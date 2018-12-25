package provisioner

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
)

// CreateCluster creates cluster
func CreateCluster(cluster *models.Cluster) error {
	fmt.Println("Creating cluster...")

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
