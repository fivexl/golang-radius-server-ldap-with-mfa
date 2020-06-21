#!/usr/bin/env bash

set -ex

TOP="$(git rev-parse --show-toplevel)"
BUILD_DIR="${TOP}/build"

mkdir -p "${BUILD_DIR}"

VERSION=${1}
[ -z "${VERSION}" ] && VERSION=beta-$(date "+%Y_%m_%d_%H_%M_%S")-$(git rev-parse HEAD)

TARGET_GOOSES=("linux" "darwin")
export GOARCH=amd64

for OS in ${TARGET_GOOSES[@]}
do
  echo "Building for ${OS}"
  export GOOS=${OS}
  go build -v -o "build/rserver-${GOOS}-${GOARCH}" -v -ldflags "-w -X main.VERSION=${VERSION}"
done