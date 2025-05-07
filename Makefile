# Variables for generating version information
VERSION ?= $(shell date '+%Y-%m-%d')
BUILD_TIME = $(shell date -u '+%Y-%m-%d %H:%M:%S %Z')
COMMIT_HASH = $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
LDFLAGS = -X 'main.version=$(VERSION)' -X 'main.buildTime=$(BUILD_TIME)' -X 'main.commitHash=$(COMMIT_HASH)'

# Check if GOPATH is set
check-gopath:
	@if [ -z "$(GOPATH)" ]; then \
		echo "Error: GOPATH environment variable is not set"; \
		echo "Please set GOPATH before running make install"; \
		exit 1; \
	fi

build:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go build -ldflags "${LDFLAGS}" -o sem

install: build
	mv sem ${GOPATH}/bin

uninstall: check-gopath
	-rm ${GOPATH}/bin/sem

# Test targets
.PHONY: test test-verbose test-coverage test-all test-integration

test:
	go test ./... -v=0

test-verbose:
	go test ./... -v

test-coverage:
	go test ./... -cover

test-all: test test-integration

test-integration:
	./tests/run_tests.sh
