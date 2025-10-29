BINARY_NAME=collector

all: test lint

test:
	@echo "Running tests..."
	@go test -v ./...

lint:
	@echo "Running linter..."
	@golangci-lint run ./...

build:
	@echo "Building binary..."
	@go build -o $(BINARY_NAME) ./cmd/collector
	@echo "$(BINARY_NAME) successfully built!"

.PHONY: test lint
