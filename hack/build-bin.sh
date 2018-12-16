#!/bin/bash
set -eo pipefail

BINARY_NAME=$1
PKG_DIR=$2

# get parent directory of this script
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

echo "==> Building for $(go env GOOS)-$(go env GOARCH)"
go build \
  -o "dist/bin/${BINARY_NAME}-$(go env GOOS)-$(go env GOARCH)/${BINARY_NAME}" \
  ./${PKG_DIR}/...

echo "==> Results:"
echo "==>./dist/bin"
cp "dist/bin/${BINARY_NAME}-$(go env GOOS)-$(go env GOARCH)/${BINARY_NAME}" dist/bin/
ls -hlR dist/bin/*