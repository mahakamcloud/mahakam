// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/swag"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations"
	"github.com/mahakamcloud/mahakam/pkg/api/v1/restapi/operations/clusters"
	"github.com/mahakamcloud/mahakam/pkg/config"
	"github.com/mahakamcloud/mahakam/pkg/handlers"
)

var opts struct {
	ConfigFilePath string `short:"c" long:"config" description:"Mahakam server config file"`
}

//go:generate swagger generate server --target .. --name Mahakam --spec ../../../../swagger/mahakam.yaml --client-package mahakam

func configureFlags(api *operations.MahakamAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		swag.CommandLineOptionsGroup{
			Options: &opts,
		},
	}
}

func configureAPI(api *operations.MahakamAPI) http.Handler {

	mahakamConfig, err := config.LoadConfig(opts.ConfigFilePath)
	if err != nil {
		log.Fatalf("Error loading configuration file for mahakam server: %s\n", err)
	}
	h := handlers.New(mahakamConfig.KVStoreConfig)

	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.ClustersCreateClusterHandler = &handlers.CreateCluster{Handlers: *h}

	api.ClustersGetClustersHandler = clusters.GetClustersHandlerFunc(func(params clusters.GetClustersParams) middleware.Responder {
		return middleware.NotImplemented("operation clusters.GetClusters has not yet been implemented")
	})

	api.ClustersDescribeClustersHandler = &handlers.DescribeCluster{Handlers: *h}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
