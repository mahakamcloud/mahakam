GIT_VERSION = $(shell git describe --tags --dirty)
VERSION ?= $(GIT_VERSION)

GO ?= go
GOVERSIONS ?= go1.9 go1.10 go1.11
OS := $(shell uname)
SHELL := /bin/bash

BASE = $(GOPATH)/src/github.com/mahakamcloud/mahakam

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: mahakam
mahakam: fmt vet ## Build mahakam cli binary
	./hack/build-bin.sh mahakam cmd/mahakam

.PHONY: fmt
fmt: ## Run go fmt against code
	go fmt ./pkg/... ./cmd/...

.PHONY: vet
vet: ## Run go vet against code
	go vet ./pkg/... ./cmd/...

.PHONY: test
test: ## run tests
	@echo running tests...
	$(GO) test -v $(shell go list -v ./... | grep -v /vendor/ | grep -v integration )

.PHONY: generate-server
generate-server: ## Generate swagger server
	cd $(BASE)/pkg/api/v1 && swagger generate server \
		-A mahakam -f $(BASE)/swagger/mahakam.yaml --exclude-main --skip-validation

.PHONY: generate-client
generate-client: ## Generate swagger client
	cd $(BASE)/pkg/api/v1 && swagger generate client \
		-A mahakam -f $(BASE)/swagger/mahakam.yaml -c client --skip-validation
