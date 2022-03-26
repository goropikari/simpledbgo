 SHELL = /bin/bash

ROOT_DIR = $(shell pwd)
GOBIN = $(ROOT_DIR)/bin
export PATH := $(ROOT_DIR)/bin:$(PATH)

MOCK_FILE := $(shell find -name "*.go" | xargs grep mockgen | cut -d: -f1)

.PHONY: tools
tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.45.0
	GOBIN=$(GOBIN) go install github.com/golang/mock/mockgen@v1.6.0
	GOBIN=$(GOBIN) go install github.com/jstemmer/go-junit-report@v1.0.0


.PHONY: test
test:
	go test ./...

ci-test:
	go test -v -cover ./... 2>&1 | go-junit-report > report.xml

.PHONY: lint
lint:
	./bin/golangci-lint run --fix ./...

mockgen:
	for f in $(MOCK_FILE); do ROOT_DIR=$(ROOT_DIR) go generate $$f; done

.PHONY: coverage
coverage:
	ci/coverage
