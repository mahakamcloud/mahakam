package provisioner

import (
	"fmt"

	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
)

// CreateCluster creates cluster
func CreateCluster(cluster *models.Cluster) error {
	fmt.Println("Creating cluster...")
	return nil
}
