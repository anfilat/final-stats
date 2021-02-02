BIN := "./bin/symo"
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%S)
GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.buildDate=$(BUILD_DATE) -X main.gitHash=$(GIT_HASH)

.PHONY: build
build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/symo

.PHONY: run
run: build
	LOG_LEVEL=DEBUG $(BIN) -config ./configs/config.toml

.PHONY: version
version: build
	$(BIN) version

.PHONY: test
test:
	go test -v -count=1 -race -timeout=30s ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.35.0

.PHONY: lint
lint: install-lint-deps
	golangci-lint run ./...

.PHONY: generate
generate:
	go generate ./...
