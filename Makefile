BINARY=satgate
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"
LDFLAGS_RELEASE=-ldflags "-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"
GOFLAGS=-trimpath

BUILD_DIR=bin
RELEASE_DIR=$(BUILD_DIR)/release

.PHONY: build release clean test

build:
	@echo "Building $(BINARY) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) .
	@echo "✓ Built $(BUILD_DIR)/$(BINARY)"

release:
	@echo "Building release $(BINARY) $(VERSION) (stripped)..."
	@mkdir -p $(RELEASE_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS_RELEASE) -o $(RELEASE_DIR)/$(BINARY)-linux-amd64 .
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(GOFLAGS) $(LDFLAGS_RELEASE) -o $(RELEASE_DIR)/$(BINARY)-linux-arm64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS_RELEASE) -o $(RELEASE_DIR)/$(BINARY)-darwin-amd64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) $(LDFLAGS_RELEASE) -o $(RELEASE_DIR)/$(BINARY)-darwin-arm64 .
	@echo "Generating checksums..."
	@cd $(RELEASE_DIR) && shasum -a 256 $(BINARY)-* > SHA256SUMS
	@echo "✓ Release binaries in $(RELEASE_DIR)/"
	@ls -lh $(RELEASE_DIR)/

test:
	go test -v -race ./...

clean:
	rm -rf $(BUILD_DIR)
