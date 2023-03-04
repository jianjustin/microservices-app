#!/bin/bash

protoc \
    -I/usr/local/include -I. -I"${GOPATH}"/src \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
    proto/add.proto