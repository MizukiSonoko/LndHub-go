
build:
	go build -o lndhub-go main.go

proto:
	protoc --proto_path $(GOPATH)/src/github.com/google/protobuf/src \
    --proto_path $(GOPATH)/src \
    --proto_path protobuf/ \
    --go_out=plugins=grpc:protobuf api.proto
