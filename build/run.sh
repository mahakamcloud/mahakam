#!/bin/bash

set -euo pipefail

: ${DOCKER_REGISTRY:="mahakamcloud"}

PACKAGE=${1}

cd $(git rev-parse --show-toplevel)

# get current commit
CURRENT_COMMIT=$(git rev-parse --short HEAD)
CURRENT_VERSION=${VERSION-"$(git describe --tags --always)"}

echo GIT_COMMIT=${CURRENT_COMMIT}
echo GIT_DESCRIBE=${CURRENT_VERSION}

# get build date
BUILD_DATE=$(date +%s)
echo BUILD_DATE=${BUILD_DATE}

# if no tag, use build date
TAG=""
if [ $CURRENT_COMMIT = $CURRENT_VERSION ]; then
  TAG=$BUILD_DATE
else
  TAG=$CURRENT_VERSION
fi

docker build --pull -t ${DOCKER_REGISTRY}/${PACKAGE}:${TAG} -f ./build/${PACKAGE}/Dockerfile .
docker push ${DOCKER_REGISTRY}/${PACKAGE}:${TAG}