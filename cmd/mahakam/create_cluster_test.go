package main

import (
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
	validCreateClusterRes = &clusters.CreateClusterCreated{
		Payload: &models.Cluster{
			ID:       1,
			Name:     swag.String("fake-cluster-name"),
			NumNodes: 2,
			Owner:    "fake-cluster-owner",
			Status:   "pending",
		},
	}
)

func TestCreateCluster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clusterAPI := v1.NewMockClusterAPI(ctrl)

	validCreateClusterOpts.ClusterAPI = clusterAPI

	clusterAPI.EXPECT().CreateCluster(gomock.Any()).
		Return(validCreateClusterRes, nil)

	// create cluster should succeed
	res, err := RunCreateCluster(validCreateClusterOpts)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), res.ID)
	assert.Equal(t, "fake-cluster-name", *res.Name)
	assert.Equal(t, int64(2), res.NumNodes)
	assert.Equal(t, "fake-cluster-owner", res.Owner)
	assert.Equal(t, "pending", res.Status)
}
