#!/usr/bin/env bash

set -ex

TOP="$(git rev-parse --show-toplevel)"

docker run --rm -v "${TOP}:/app" -w /app golangci/golangci-lint:v1.27.0 golangci-lint run -v