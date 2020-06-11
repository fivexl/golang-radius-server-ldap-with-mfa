#!/usr/bin/env bash

set -ex

export GO111MODULE=on

BUILD_DIR="$PWD/build"
GOPACKAGE_NAME="github.com/fivexl/ldap-radius-with-mfa"

mkdir -p "$BUILD_DIR"
cd "$BUILD_DIR"

VERSION=${1}
[ -z "${VERSION}" ] && VERSION=beta

TARGET_GOOSES=("linux" "darwin")

for OS in ${TARGET_GOOSES[@]}
do
  echo "Building for ${OS}"
  export GOOS=${OS}
  export GOARCH=amd64
  export CGO_ENABLED=0
  go build -o "server-${GOOS}-${GOARCH}" -v -ldflags "-w -X main.VERSION=${VERSION}" "$GOPACKAGE_NAME"
done