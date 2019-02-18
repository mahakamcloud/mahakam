# Mahakam [![Build Status](https://circleci.com/gh/mahakamcloud/mahakam.svg?style=shield)](https://circleci.com/gh/mahakamcloud/mahakam) [![codecov](https://codecov.io/gh/mahakamcloud/mahakam/branch/master/graph/badge.svg)](https://codecov.io/gh/mahakamcloud/mahakam)
Cloud application platform on Kubernetes

## Common Development Task

### Unit tests

To run unit tests, execute the following command:

```
$ make test
```

### Dev server

To run dev server, first populate necessary info in `pkg/config/example/config.sample.yaml`. Or, you can create new config yaml file and change the volume mount in `docker-compose.dev.yaml` accordingly.

Run dev server and consul with docker-compose:

```
$ make server-linux
$ make dev
```

### CLI

Building mahakam cli as per your machine, find the build under `dist/bin`

```
$ make cli
```

### Server API

Generate mahakam server api using swagger

```
$ make generate-server
```

### Client API

Generate mahakam client api using swagger

```
$ make generate-client
```

### Using go modules

If you're using go path for src

```
$ export GO111MODULE=on
```

Add new dependency

```
$ go get golang.org/x/text@v0.3.0
```

Add vendoring

```
$ go mod vendor
```

Building with vendor dir

```
$ go build -mod=vendor
```
