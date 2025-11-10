PROTO_PATH       ?= pkg/proto
OUTPUT_PATH_GO   ?= pkg/gen/go
OUTPUT_PATH_PY   ?= pkg/gen/py
OUTPUT_PATH_CPP  ?= pkg/gen/cpp

PROTO_FILES ?= \
    user/user.proto \

proto:
	@mkdir -p $(OUTPUT_PATH_GO)
	@protoc -I $(PROTO_PATH) $(addprefix $(PROTO_PATH)/, $(PROTO_FILES)) \
		--go_out=$(OUTPUT_PATH_GO) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUTPUT_PATH_GO) --go-grpc_opt=paths=source_relative

build:
	@go build -o bin/core ./cmd/core/main.go

run: build
	@./bin/core --config=./config/local.yaml

test:
	@go test -v ./...
