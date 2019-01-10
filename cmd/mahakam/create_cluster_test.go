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
	validCreateClusterOpts = &CreateClusterOptions{
		Name:     "fake-cluster-name",
		Owner:    "fake-cluster-owner",
		NumNodes: 2,
	}
	createClusterRes = &clusters.CreateClusterCreated{
		Payload: &models.Cluster{
			ID:       1,
			Name:     swag.String("fake-cluster-name"),
			NumNodes: 2,
			Owner:    "fake-cluster-owner",
			Status:   "pending",
		},
	}
	createClusterErr = &clusters.CreateClusterCreated{}
)

func TestCreateCluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clusterAPI := v1.NewMockClusterAPI(ctrl)

	validCreateClusterOpts.ClusterAPI = clusterAPI

	// create cluster should succeed
	clusterAPI.EXPECT().CreateCluster(gomock.Any()).
		Return(createClusterRes, nil)
	res, err := RunCreateCluster(validCreateClusterOpts)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), res.ID)
	assert.Equal(t, "fake-cluster-name", *res.Name)
	assert.Equal(t, int64(2), res.NumNodes)
	assert.Equal(t, "fake-cluster-owner", res.Owner)
	assert.Equal(t, "pending", res.Status)

	// create cluster should error out
	clusterAPI.EXPECT().CreateCluster(gomock.Any()).
		Return(createClusterErr, fmt.Errorf("fake create cluster error"))
	res, err = RunCreateCluster(validCreateClusterOpts)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "fake create cluster error")
}
