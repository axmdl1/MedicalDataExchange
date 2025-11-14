build:
	@go build -o bin/crm ./cmd/main.go

run: build
	@./bin/crm --config=./config/local.yaml

test:
	@go test -v ./...
