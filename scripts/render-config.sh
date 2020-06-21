
#!/usr/bin/env bash

set -ex

TOP="$(git rev-parse --show-toplevel)"
BUILD_DIR="${TOP}/build"

mkdir -p "${BUILD_DIR}"

secrethub inject --force --in-file "${TOP}/scripts/config.gcfg.tpl" --out-file "${BUILD_DIR}/config.gcfg"