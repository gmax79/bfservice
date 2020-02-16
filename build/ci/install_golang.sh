#!/bin/bash

GOLANG_VERSION="1.13.7"

curl -OL https://dl.google.com/go/go$GOLANG_VERSION.linux-amd64.tar.gz
tar -C / -xzf go$GOLANG_VERSION.linux-amd64.tar.gz

