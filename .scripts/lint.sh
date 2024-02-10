#!/usr/bin/env bash
set -x
GOLANGCI_LINT_VERSION=v1.43.0

apt update
apt-get install -y jq
cd /$1

go mod tidy
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin $GOLANGCI_LINT_VERSION
golangci-lint cache clean
golangci-lint run -c .golangci.yml --issues-exit-code 0 --out-format code-climate | tee gl-code-quality-report.json | jq -r '.[] | "\(.location.path):\(.location.lines.begin) \(.description)"'
golangci-lint run -c .golangci.yml --issues-exit-code 0 --out-format checkstyle > sonarqube_report.xml
