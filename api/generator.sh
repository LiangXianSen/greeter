# brew install protobuf (MacOS)
# go get -u google.golang.org/protobuf/cmd/protoc-gen-go
# go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
# go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
# go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2

protoc -I . -I $GOPATH/src/github.com/googleapis/googleapis \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=. --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    --grpc-gateway_opt generate_unbound_methods=true \
    --openapiv2_out . --openapiv2_opt logtostderr=true \
    greeter.proto