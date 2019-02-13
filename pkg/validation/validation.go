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
	ValidateNWithDelay(owner, clustername string, timeout time.Duration, log logrus.FieldLogger, count int, delay time.Duration) bool
}

type ClusterValidator struct {
	client *client.Mahakam
}

func NewClusterValidator(client *client.Mahakam) Validator {
	return &ClusterValidator{
		client: client,
	}
}

func (v *ClusterValidator) validate(owner, clustername string, log logrus.FieldLogger) (bool, error) {
	req := &models.Cluster{
		Owner: owner,
		Name:  swag.String(clustername),
	}
	res, err := v.client.Clusters.ValidateCluster(clusters.NewValidateClusterParams().WithBody(req))
	if err != nil {
		return false, fmt.Errorf("error validating cluster %s: %s", clustername, err)
	}
	log.Debugf("cluster validation result: %v", res)

	if len(res.Payload.NodeFailures) > 0 ||
		len(res.Payload.ComponentFailures) > 0 ||
		len(res.Payload.PodFailures) > 0 {
		return false, nil
	}
	return true, nil
}

func (v *ClusterValidator) ValidateNWithDelay(owner, clustername string, timeout time.Duration, log logrus.FieldLogger,
	count int, delay time.Duration) bool {
	for i := 0; i < count; i++ {
		ready, err := v.validate(owner, clustername, log)
		if err == nil && ready {
			log.Infof("cluster %s is ready", clustername)
			return true
		}
		log.Warnf("cluster %s is not ready, retrying after %f seconds", clustername, delay.Seconds())
		time.Sleep(delay)
	}
	log.Errorf("validating cluster %s timeout", clustername)
	return false
}
