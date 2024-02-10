#!/usr/bin/env bash
set -x
cd /$1

apt update
go mod tidy
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build  -o $2
mkdir build
mv $2 build/
