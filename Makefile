ROOT_DIR = $(shell pwd)
export PATH := $(ROOT_DIR)/bin:$(PATH)

MOCK_FILE := $(shell find -name "*.go" | xargs grep mockgen | cut -d: -f1)

tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.44.2

test:
	go test ./...

lint:
	./bin/golangci-lint run ./...

mockgen:
	ROOT_DIR=$(ROOT_DIR) go generate $(MOCK_FILE)
