package main

import (
	"fmt"
	"testing"

	"github.com/go-openapi/swag"
	"github.com/golang/mock/gomock"

	"github.com/mahakamcloud/mahakam/pkg/api/v1"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/client/clusters"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/stretchr/testify/assert"
)

var (
	validValidateClusterOpts = &ValidateClusterOptions{
		Name:  "fake-cluster-name",
		Owner: "fake-cluster-owner",
	}
	validateClusterRes = &clusters.ValidateClusterCreated{
		Payload: &models.Cluster{
			Name:  swag.String("fake-cluster-name"),
			Owner: "fake-cluster-owner",
		},
	}
	validateClusterErr = &clusters.ValidateClusterCreated{}
)

func TestValidateClusterWithoutFailures(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clusterAPI := v1.NewMockClusterAPI(ctrl)

	validValidateClusterOpts.ClusterAPI = clusterAPI

	// validate cluster should succeed
	clusterAPI.EXPECT().ValidateCluster(gomock.Any()).
		Return(validateClusterRes, nil)
	res, err := RunValidateCluster(validValidateClusterOpts)
	assert.NoError(t, err)
	assert.Equal(t, "fake-cluster-name", *res.Name)
	assert.Equal(t, "fake-cluster-owner", res.Owner)
	assert.Equal(t, 0, len(res.NodeFailures))
	assert.Equal(t, 0, len(res.ComponentFailures))
	assert.Equal(t, 0, len(res.PodFailures))
}

func TestValidateClusterWithFailures(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clusterAPI := v1.NewMockClusterAPI(ctrl)

	validValidateClusterOpts.ClusterAPI = clusterAPI

	validateClusterRes.Payload.NodeFailures = []string{"fake-node-failures"}
	validateClusterRes.Payload.ComponentFailures = []string{"fake-component-failures"}
	validateClusterRes.Payload.PodFailures = []string{"fake-pod-failures"}

	// validate cluster should succeed
	clusterAPI.EXPECT().ValidateCluster(gomock.Any()).
		Return(validateClusterRes, nil)
	res, err := RunValidateCluster(validValidateClusterOpts)
	assert.NoError(t, err)
	assert.Equal(t, "fake-cluster-name", *res.Name)
	assert.Equal(t, "fake-cluster-owner", res.Owner)
	assert.Equal(t, 1, len(res.NodeFailures))
	assert.Equal(t, 1, len(res.ComponentFailures))
	assert.Equal(t, 1, len(res.PodFailures))
}

func TestValidateClusterWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clusterAPI := v1.NewMockClusterAPI(ctrl)

	validValidateClusterOpts.ClusterAPI = clusterAPI

	// validate cluster should error out
	clusterAPI.EXPECT().ValidateCluster(gomock.Any()).
		Return(validateClusterErr, fmt.Errorf("fake create cluster error"))
	res, err := RunValidateCluster(validValidateClusterOpts)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "fake create cluster error")
}
