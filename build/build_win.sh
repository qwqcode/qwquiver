#!/bin/sh

set -e

cd ../

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o "./build/qwquiver.exe" ./cli