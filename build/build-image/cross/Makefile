.PHONY: build push

IMAGE=maha-cross
TAG=$(shell cat VERSION)


all: push

build:
	docker build --pull -t mahakamcloud/$(IMAGE):$(TAG) .

push: build
	docker push mahakamcloud/$(IMAGE):$(TAG)