#!/bin/bash

GOLANG_VERSION="1.13.7"

curl -OL https://dl.google.com/go/go$GOLANG_VERSION.linux-amd64.tar.gz
tar -C /usr/local -xzf go$GOLANG_VERSION.linux-amd64.tar.gz
export PATH=/usr/local/go:$PATH
