 SHELL = /bin/bash

ROOT_DIR = $(shell pwd)
GOBIN = $(ROOT_DIR)/bin
export PATH := $(ROOT_DIR)/bin:$(PATH)

MOCK_FILE := $(shell find -name "*.go" | xargs grep mockgen | cut -d: -f1)

tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.45.0
	GOBIN=$(GOBIN) go install github.com/golang/mock/mockgen@v1.6.0

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	./bin/golangci-lint run --fix ./...

mockgen:
	for f in $(MOCK_FILE); do ROOT_DIR=$(ROOT_DIR) go generate $$f; done
