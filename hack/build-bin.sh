#!/bin/bash
set -eo pipefail

BINARY_NAME=$1
PKG_DIR=$2
BUILD_ENV=$3

if [ -z "$BUILD_ENV" ]; then
  BUILD_ENV=$(go env GOOS)
fi

# get parent directory of this script
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

echo "==> Building for $BUILD_ENV"
GOOS=${BUILD_ENV} go build \
  -o "dist/bin/${BINARY_NAME}-${BUILD_ENV}/${BINARY_NAME}" \
  ./${PKG_DIR}/...

echo "==> Results:"
echo "==>./dist/bin"
cp "dist/bin/${BINARY_NAME}-${BUILD_ENV}/${BINARY_NAME}" dist/bin/
ls -hlR dist/bin/*