#!/usr/bin/env bash

set -ex

VERSION=${1}
[ -z "${VERSION}" ] && VERSION=$(git describe --tags)

TOP="$(git rev-parse --show-toplevel)"
REPO_NAME="$(basename ${TOP})"
BUILD_DIR="${TOP}/build"
RELEASE_DIR="${TOP}/release/${REPO_NAME}/${VERSION}"

mkdir -p "${BUILD_DIR}"
mkdir -p "${RELEASE_DIR}"

TARGET_GOOSES=("linux" "darwin" "windows")
TARGET_GOARCHES=("amd64" "arm")

for OS in "${TARGET_GOOSES[@]}"
do
    for ARCH in "${TARGET_GOARCHES[@]}"
    do
        if [ "${ARCH}" == "arm" ] && [ "${OS}" != "linux" ]; then
            continue
        fi
        echo "Building for ${OS} and arch ${ARCH}"
        export GOOS="${OS}"
        export GOARCH="${ARCH}"
        go build -o "build/rserver-${GOOS}-${GOARCH}" -ldflags "-s -w -X main.VERSION=${VERSION}"
        cp "build/rserver-${GOOS}-${GOARCH}" build/rserver
        zip -j "${RELEASE_DIR}/rserver_${VERSION}_${OS}_${ARCH}.zip" build/rserver
        rm -rf build/rserver
    done
done

cd "${RELEASE_DIR}"
for FILE in *.zip
do
    sha256sum "${FILE}" > "${FILE%.zip}_sha256"
done
