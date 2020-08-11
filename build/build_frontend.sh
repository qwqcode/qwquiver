#!/bin/sh

set -e

cd ../frontend

yarn build

cd ../

go-bindata -fs -prefix "frontend/dist/" -pkg "bindata" -o "bindata/bindata.go" "./frontend/dist/..."
