#!/bin/bash

apt-get update -y
apt-get install unzip

PROTOC_VERSION="3.9.1"
PROTOC_ZIP="protoc-$PROTOC_VERSION-linux-x86_64.zip"
curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOC_VERSION/$PROTOC_ZIP
unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
rm -f $PROTOC_ZIP

GIT_TAG="v1.3.2" # protoc-gen-go version
go get -d -u github.com/golang/protobuf/protoc-gen-go
mkdir -p "$(go env GOPATH)"/src/github.com/golang/protobuf
git -C "$(go env GOPATH)"/src/github.com/golang/protobuf checkout $GIT_TAG
go install github.com/golang/protobuf/protoc-gen-go
