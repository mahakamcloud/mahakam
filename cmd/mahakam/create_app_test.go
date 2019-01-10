package main

import (
	"fmt"
	"testing"

	"github.com/go-openapi/swag"
	"github.com/golang/mock/gomock"
	"github.com/mahakamcloud/mahakam/pkg/api/v1"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/client/apps"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/stretchr/testify/assert"
)

var (
	validCreateAppOpts = &CreateAppOptions{
		Name:        "fake-app-name",
		Owner:       "fake-app-owner",
		ClusterName: "fake-cluster-name",
		ChartURL:    "fake-chart-url",
		ChartValues: "",
	}
	createAppRes = &apps.CreateAppCreated{
		Payload: &models.App{
			ID:          1,
			Name:        swag.String("fake-app-name"),
			ClusterName: "fake-cluster-name",
			ChartURL:    "fake-chart-url",
			ChartValues: "",
			ServiceFQDN: "fake-service-fqdn",
			Owner:       "fake-app-owner",
			Status:      "pending",
		},
	}
	createAppErr = &apps.CreateAppCreated{}
)

func TestCreateApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appAPI := v1.NewMockAppAPI(ctrl)

	validCreateAppOpts.AppAPI = appAPI

	// create app should succeed
	appAPI.EXPECT().CreateApp(gomock.Any()).
		Return(createAppRes, nil)
	res, err := RunCreateApp(validCreateAppOpts)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), res.ID)
	assert.Equal(t, "fake-app-name", *res.Name)
	assert.Equal(t, "fake-cluster-name", res.ClusterName)
	assert.Equal(t, "fake-chart-url", res.ChartURL)
	assert.Equal(t, "", res.ChartValues)
	assert.Equal(t, "fake-service-fqdn", res.ServiceFQDN)
	assert.Equal(t, "fake-app-owner", res.Owner)
	assert.Equal(t, "pending", res.Status)

	// create app should error out
	appAPI.EXPECT().CreateApp(gomock.Any()).
		Return(createAppErr, fmt.Errorf("fake create app error"))
	res, err = RunCreateApp(validCreateAppOpts)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "fake create app error")
}
