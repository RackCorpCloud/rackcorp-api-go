
GO_FILES = $(shell find . -name '*.go')

default: build lint test

build: $(GO_FILES) go.mod go.sum
	go build ./...
.PHONY: build

test:
	go test ./...
.PHONY: test

integration-test:
	. ./.env && INTEGRATION_TEST=1 go test ./...
.PHONY: integration-test

lint:
	golangci-lint run
.PHONY: lint

install-lint:
	# https://golangci-lint.run/welcome/install/#binaries
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.4.0
.PHONY: install-lint
