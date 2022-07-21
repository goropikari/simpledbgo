 SHELL = /bin/bash

ROOT_DIR = $(shell pwd)
GOBIN = $(ROOT_DIR)/bin
export PATH := $(ROOT_DIR)/bin:$(PATH)

MOCK_FILE := $(shell find -name "*.go" | xargs grep mockgen | cut -d: -f1)

build:
	go build -o simpledb main.go

.PHONY: tools
tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.45.2
	GOBIN=$(GOBIN) go install github.com/golang/mock/mockgen@v1.6.0
	GOBIN=$(GOBIN) go install github.com/jstemmer/go-junit-report@v1.0.0
	GOBIN=$(GOBIN) go install github.com/jandelgado/gcov2lcov@v1.0.5
	GOBIN=$(GOBIN) go install golang.org/x/tools/cmd/godoc
	GOBIN=$(GOBIN) go install github.com/google/wire/cmd/wire@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
	curl -Lf -o protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v3.20.0/protoc-3.20.0-linux-x86_64.zip
	unzip protoc.zip bin/protoc
	rm -f protoc.zip

.PHONY: docker-build
docker-build:
	docker build -t simpledb .

.PHONY: docker-run
docker-run:
	docker run --rm -p 5432:5432 simpledb

.PHONY: test
test: mockgen protoc wire
	go test -timeout 20s -shuffle=on ./...

.PHONY: test-all
test-all: mockgen protoc wire docker-build
	go test -timeout 20s -shuffle=on ./...
	bash e2e/test.sh

.PHONY: ci-test
ci-test:
	go test -v -cover -shuffle=on ./... 2>&1 | go-junit-report > report.xml

.PHONY: lint
lint:
	./bin/golangci-lint run --fix ./...

mockgen:
	rm -rf testing/mock
	for f in $(MOCK_FILE); do ROOT_DIR=$(ROOT_DIR) go generate $$f; done

.PHONY: protoc
protoc:
	bin/protoc -I=. --go_out=./tx/logrecord ./tx/logrecord/protofile/*.proto

.PHONY: wire
wire:
	bin/wire ./database

.PHONY: coverage
coverage:
	go test -cover ./... -coverprofile=coverage.out
	bin/gcov2lcov -infile=coverage.out -outfile=coverage.lcov
	genhtml coverage.lcov -o docs/coverage

.PHONY: site
site: coverage
	mkdocs build
	mkdocs serve

.PHONY: godoc
godoc:
	bin/godoc -http=:8080 &
	sleep 5
	wget -r -np -N -E -k http://localhost:8080/pkg/github.com/goropikari/simpledbgo
	mv localhost:8080/pkg/github.com/goropikari docs
