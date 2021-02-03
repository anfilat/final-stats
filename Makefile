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

LOAD_BIN := "./bin/load"

.PHONY: build-pprof
build-pprof:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" -tags pprof ./cmd/symo

.PHONY: build-load
build-load:
	go build -v -o $(LOAD_BIN) ./cmd/load

.PHONY: test-load
test-load: build-load build-pprof
	$(LOAD_BIN)
