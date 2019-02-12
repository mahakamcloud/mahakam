package validation

import (
	"fmt"
	"time"

	"github.com/go-openapi/swag"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/client"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/client/clusters"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/sirupsen/logrus"
)

type ValidationError struct {
	Name    string
	Message string
}

type Validator interface {
	ValidateNWithDelay(clustername string, timeout time.Duration, log logrus.FieldLogger, count int, delay time.Duration) bool
}

type ClusterValidator struct {
	client *client.Mahakam
}

func NewClusterValidator(client *client.Mahakam) Validator {
	return &ClusterValidator{
		client: client,
	}
}

func (v *ClusterValidator) validate(clustername string) (bool, error) {
	req := &models.Cluster{
		Name: swag.String(clustername),
	}
	res, err := v.client.Clusters.ValidateCluster(clusters.NewValidateClusterParams().WithBody(req))
	if err != nil {
		return false, fmt.Errorf("error validating cluster %s: %s", clustername, err)
	}

	if len(res.Payload.NodeFailures) > 0 ||
		len(res.Payload.ComponentFailures) > 0 ||
		len(res.Payload.PodFailures) > 0 {
		return false, nil
	}
	return true, nil
}

func (v *ClusterValidator) ValidateNWithDelay(clustername string, timeout time.Duration, log logrus.FieldLogger,
	count int, delay time.Duration) bool {
	for i := 0; i < count; i++ {
		ready, err := v.validate(clustername)
		if err != nil && ready {
			log.Infof("cluster %s is ready", clustername)
			return true
		}
		time.Sleep(delay)
	}
	log.Errorf("validating cluster %s timeout", clustername)
	return false
}
