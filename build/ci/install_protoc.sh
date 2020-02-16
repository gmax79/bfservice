#!/bin/bash

apt-get update -y
apt-get install unzip

PROTOC_VERSION="3.9.1"
PROTOC_ZIP="protoc-$PROTOC_VERSION-linux-x86_64.zip"
curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOC_VERSION/$PROTOC_ZIP
unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
unzip -o $PROTOC_ZIP -d /usr/local 'include/*'

GOGEN_VERSION="1.3.2" # protoc-gen-go version
GOGEN_ZIP="protobuf-$GOGEN_VERSION.zip"
UPACKED="protobuf-$GOGEN_VERSION"
curl -o $GOGEN_ZIP -OL https://github.com/golang/protobuf/archive/v$GOGEN_VERSION.zip 
unzip -o $GOGEN_ZIP
cd $UPACKED
make install
