#!/bin/bash
set -eo pipefail

cd $(git rev-parse --show-toplevel)

echo "==> Getting current commit, version, and build date"
CURRENT_COMMIT=$(git rev-parse --short HEAD)
CURRENT_VERSION=${VERSION-"$(git describe --tags --always)"}

echo GIT_COMMIT=${CURRENT_COMMIT}
echo GIT_DESCRIBE=${CURRENT_VERSION}

BUILD_DATE=$(date -I'seconds')
echo BUILD_DATE=${BUILD_DATE}
