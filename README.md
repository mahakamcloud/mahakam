# Mahakam [![Build Status](https://circleci.com/gh/mahakamcloud/mahakam.svg?style=shield)](https://circleci.com/gh/mahakamcloud/mahakam)
Cloud application platform on Kubernetes

## Common Development Task
Running unit test
```
$ make test
```

Run dev store with consul backend:
```
$ make dev-store
```

Run dev server:
```
$ make dev-server
```

Building mahakam cli as per your machine, find the build under `dist/bin`
```
$ make mahakam-cli
```

Generate mahakam server api using swagger
```
$ make generate-server
```

Generate mahakam client api using swagger
```
$ make generate-client
```
