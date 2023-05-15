# go-microservices

生成proto文件

```shell
protoc --proto_path=../common/proto \
    --go_out=./proto --go_opt=paths=source_relative \
    --go-grpc_out=./proto --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false \
    --grpc-gateway_out=./proto --grpc-gateway_opt=paths=source_relative --plugin=protoc-gen-grpc-gateway=$GOPATH/bin/protoc-gen-grpc-gateway \
    calculate.proto
```