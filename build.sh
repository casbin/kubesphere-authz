#!/bin/bash
mkdir build
echo "===build executable binary for webhook==="
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy

cp -r config build/config
# todo: remove file version of policies and models if casbin CRD adaptor is finished
cp -r example build/example
go build -o build/main cmd/webhook/main.go 

echo "=====build docker images for webhook====="
docker build -t casbin-kubesphere-authz:v1 .
