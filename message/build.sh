#!/bin/sh
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --proto_path=. *.proto

sed -i "" -e "s/,omitempty//g" ./*.go
goimports -w .