#!/bin/sh

set -e

cd ../

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o "./build/qwquiver" ./cli
