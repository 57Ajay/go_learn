# Variables
BINARY_NAME := app
BUILD_DIR := build
SRC := $(shell find . -name '*.go' -type f)

# Commands
GO := go
GOFMT := gofmt
GOLINT := golangci-lint

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) ./...

# Run the application
.PHONY: run
run:
	$(GO) run main.go

# Clean build artifacts
.PHONY: clean
clean:
	@rm -rf $(BUILD_DIR)

# Format the code
.PHONY: fmt
fmt:
	$(GOFMT) -s -w .

# Run tests
.PHONY: test
test:
	$(GO) test ./... -v

# Run linting
.PHONY: lint
lint:
	$(GOLINT) run ./...

# Install dependencies
.PHONY: deps
deps:
	$(GO) mod tidy
	$(GO) mod download

# Debug build
.PHONY: debug
debug:
	@mkdir -p $(BUILD_DIR)
	$(GO) build -gcflags="all=-N -l" -o $(BUILD_DIR)/$(BINARY_NAME) ./...

# Release build
.PHONY: release
release:
	@mkdir -p $(BUILD_DIR)
	$(GO) build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) ./...

