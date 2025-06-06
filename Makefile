# Variables for generating version information
BUILD_TIME = $(shell date -u '+%Y-%m-%d %H:%M:%S %Z')
COMMIT_HASH = $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
GIT_TAG = $(shell git describe --tags --abbrev=0 2>/dev/null || echo "unknown")
VERSION ?= $(GIT_TAG)
LDFLAGS = -X 'main.version=$(VERSION)' -X 'main.date=$(BUILD_TIME)' -X 'main.commit=$(COMMIT_HASH)'

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
	mv sem ${GOPATH}/bin/sem

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
