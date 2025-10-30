BINARY_NAME=collector
# Exclude specific folders from tests.
TESTFOLDER := $(shell go list ./... | grep -v '/cmd')

all: test lint

test:
	@echo "Running tests..."
	@go test  -v -race -covermode=atomic -coverprofile=coverage.txt $(TESTFOLDER)

lint:
	@echo "Running linter..."
	@golangci-lint run ./...

build:
	@echo "Building binary..."
	@go build -o $(BINARY_NAME) ./cmd/collector
	@echo "$(BINARY_NAME) successfully built!"

.PHONY: test lint
