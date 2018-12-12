# Mahakam [![Build Status](https://circleci.com/gh/MahakamCloud/mahakam.svg?style=svg)](https://circleci.com/gh/MahakamCloud/mahakam)
Cloud application platform on Kubernetes

## Common Development Task
Running unit test
```
$ make test
```

Building mahakam cli as per your machine, find the build under `dist/bin`
```
$ make mahakam
```

Generate mahakam server api using swagger
```
$ make generate-server
```

Generate mahakam client api using swagger
```
$ make generate-client
```
