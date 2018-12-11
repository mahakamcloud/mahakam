#!/bin/bash
set -eo pipefail

BINARY_NAME=$1
PKG_DIR=$2

# get parent directory of this script
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

echo "==> Building for $(go env GOOS)-$(go env GOARCH)"
XC_OS=$(go env GOOS)
XC_ARCH=$(go env GOARCH)

LDFLAGS="-s -w"

echo "==> Building ..."
gox \
  -verbose \
  -os="${XC_OS}" \
  -arch="${XC_ARCH}" \
  -osarch="!darwin/arm" \
  -ldflags "${LDFLAGS}" \
  -output "dist/bin/${BINARY_NAME}-{{.OS}}-{{.Arch}}/${BINARY_NAME}" \
  ./${PKG_DIR}/...

echo "==> Results:"
echo "==>./dist/bin"
cp "dist/bin/${BINARY_NAME}-$(go env GOOS)-$(go env GOARCH)/${BINARY_NAME}" dist/bin/
ls -hlR dist/bin/*