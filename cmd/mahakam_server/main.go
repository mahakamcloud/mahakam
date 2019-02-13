package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	errors "github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/apps"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/networks"
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/handlers"
	"github.com/mahakamcloud/mahakam/pkg/provisioner"
	utils "github.com/mahakamcloud/mahakam/pkg/utils"
)

type mahakamServerOpts struct {
	configFilePath string
}

func main() {
	log := logrus.New().WithField("app", "mahakam")
	log.Println("initializing...")
	defer log.Println("exiting...")

	opts := registerFlags()
	pflag.Parse()

	conf, err := config.LoadConfig(opts.configFilePath)
	if err != nil {
		log.Fatalf("error loading config file for mahakam server: %s\n", err)
	}

	pingCheck := utils.NewPingCheck()

	storageBackendCheck := config.NewCheckStorageBackendConnection(conf.KVStoreConfig, log, pingCheck)

	err = storageBackendCheck.ValidateAvailability()
	if err != nil {
		log.Fatalf("Error connecting to backend server: %s\n", err)
	}

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewMahakamAPI(swaggerSpec)
	api.Logger = log.Printf

	server := restapi.NewServer(api)
	defer server.Shutdown()

	app := handlers.New(
		conf,
		provisioner.NewTerraformProvisioner(conf.TerraformConfig),
		log,
	)
	registerHandlers(app, api, opts)

	configureServer(server, conf.MahakamServerConfig.Address, conf.MahakamServerConfig.Port)
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

func registerFlags() *mahakamServerOpts {
	opts := &mahakamServerOpts{}
	fs := pflag.NewFlagSet("mahakam-server-opts", pflag.ExitOnError)
	pflag.StringVarP(&opts.configFilePath, "config", "c", "", "mahakam server config file")
	pflag.CommandLine.AddFlagSet(fs)
	return opts
}

func registerHandlers(app *handlers.Handlers, api *operations.MahakamAPI, opts *mahakamServerOpts) {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.ClustersCreateClusterHandler = handlers.NewCreateClusterHandler(*app)

	api.ClustersGetClustersHandler = handlers.NewGetClusterHandler(*app)

	api.ClustersDescribeClustersHandler = &handlers.DescribeCluster{Handlers: *app}

	api.ClustersValidateClusterHandler = handlers.NewValidateClusterHandler(*app)

	api.NetworksCreateNetworkHandler = handlers.NewCreateNetworkHandler(*app)

	api.NetworksGetNetworksHandler = networks.GetNetworksHandlerFunc(func(params networks.GetNetworksParams) middleware.Responder {
		return middleware.NotImplemented("operation apps.GetNetworks has not yet been implemented")
	})

	api.NetworksCreateIPPoolHandler = handlers.NewCreateIPPoolHandler(*app)

	api.NetworksAllocateOrReleaseFromIPPoolHandler = handlers.NewAllocateOrReleaseFromIPPool(*app)

	api.AppsCreateAppHandler = handlers.NewCreateAppHandler(*app)

	api.AppsGetAppsHandler = apps.GetAppsHandlerFunc(func(params apps.GetAppsParams) middleware.Responder {
		return middleware.NotImplemented("operation apps.GetApps has not yet been implemented")
	})

	api.AppsUploadAppValuesHandler = handlers.NewUploadAppValuesHandler(*app)

	api.ServerShutdown = func() {}
}

func configureServer(s *restapi.Server, host string, port int) {
	s.Host = host
	s.Port = port
}
