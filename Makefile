.PHONY: gen test run-server run-client tidy

PROTOC ?= $(LOCALAPPDATA)/protoc/bin/protoc
export PATH := $(shell go env GOPATH)/bin:$(LOCALAPPDATA)/protoc/bin:$(PATH)

gen:
	@mkdir -p gen/user/v1
	$(PROTOC) -I api \
		--go_out=gen --go_opt=module=github.com/example/grpc-user-service/gen \
		--go-grpc_out=gen --go-grpc_opt=module=github.com/example/grpc-user-service/gen \
		api/user/v1/user.proto

test:
	go test ./...

run-server:
	go run ./cmd/server -addr :50051 -http :8080

run-client:
	go run ./cmd/client -action demo

tidy:
	go mod tidy
