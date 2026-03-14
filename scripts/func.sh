#!/bin/bash

statan() {
  local path=$1
  local flag=$2
  export GOFLAGS="-buildvcs=false"
  eval "golangci-lint run -v ${flag} --timeout=2m ${path}"
  status=$?
  if [ ${status} -ne 0 ]; then
    echo "error: linter failed, please fix the errors..."
    exit ${status}
  fi
  eval "go vet ${path}"
  status=$?
  if [ ${status} -ne 0 ]; then
    echo "error: vetting failed, please fix the errors..."
    exit ${status}
  fi
  eval "govulncheck -show verbose ${path}"
  status=$?
  if [ ${status} -ne 0 ]; then
    echo "error: vulnerability check failed, please fix the errors..."
    exit ${status}
  fi
}
