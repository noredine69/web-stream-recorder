#!/usr/bin/env bash
set -x
cd /$1
apt update
go install github.com/jstemmer/go-junit-report@latest
go test ./... -coverprofile=coverage.out -v 2>&1 | $GOPATH/bin/go-junit-report -set-exit-code > report.xml
