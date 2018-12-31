package handlers

import (
	"io/ioutil"
	"os"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/models"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/apps"
	"github.com/mahakamcloud/mahakam/pkg/config"
	log "github.com/sirupsen/logrus"
)

var (
	defaultBufReader = 1024 * 10
)

type UploadAppValues struct {
	Handlers
}

func (h *UploadAppValues) Handle(params apps.UploadAppValuesParams) middleware.Responder {
	chartValuesFile := getChartValuesFile(
		swag.StringValue(params.Owner),
		swag.StringValue(params.ClusterName),
		swag.StringValue(params.AppName),
	)

	if err := os.MkdirAll(config.HelmDefaultChartValuesDirectory, 0700); err != nil {
		log.Errorf("error creating directory for chart values: %v\n", err)
		return apps.NewUploadAppValuesDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: "cannot create default directory for helm chart values",
		})
	}

	buf := make([]byte, defaultBufReader)
	size, err := params.Values.Read(buf)

	err = ioutil.WriteFile(config.HelmDefaultChartValuesDirectory+chartValuesFile, buf[:size], 0644)
	if err != nil {
		log.Errorf("error writing helm chart values file: %v\n", err)
		return apps.NewUploadAppValuesDefault(405).WithPayload(&models.Error{
			Code:    405,
			Message: "cannot upload helm chart values",
		})
	}

	return apps.NewUploadAppValuesCreated()
}
