ROOT_DIR = $(shell pwd)
export PATH := $(ROOT_DIR)/bin:$(PATH)

MOCK_FILE := $(shell find -name "*.go" | xargs grep mockgen | cut -d: -f1)

test:
	go test ./...

mockgen:
	ROOT_DIR=$(ROOT_DIR) go generate $(MOCK_FILE)
