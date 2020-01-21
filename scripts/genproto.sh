#!/bin/bash

cd ../api/grpc
protoc -I. -I/usr/local/include --go_out=plugins=grpc:. abf.proto

