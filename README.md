# Mahakam [![Build Status](https://circleci.com/gh/mahakamcloud/mahakam.svg?style=shield)](https://circleci.com/gh/mahakamcloud/mahakam)
Cloud application platform on Kubernetes

## Common Development Task
Running unit test
```
$ make test
```

Run consul server:
```
$ docker run -d --name=dev-consul -e CONSUL_BIND_INTERFACE=eth0 -p 8500:8500 -p 8600:8600 consul
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
