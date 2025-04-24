# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=rcc
VERSION?=0.0.0

# Test parameters
TEST_PACKAGES=./...
TEST_FLAGS=-v -race -cover

.PHONY: all test build clean run lint help

all: test build

# Run all tests
test:
	$(GOTEST) $(TEST_PACKAGES) $(TEST_FLAGS)

# Run CLI tests specifically
test-cli:
	$(GOTEST) ./test/cli/... $(TEST_FLAGS)

# Run buildpacks tests specifically
test-buildpacks:
	$(GOTEST) ./test/buildpacks/... $(TEST_FLAGS)

# Run goreleaser tests specifically
test-goreleaser:
	$(GOTEST) ./test/goreleaser/... $(TEST_FLAGS)

# Build the binary
build:
	$(GOBUILD) -o $(BINARY_NAME) -v

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf dist/

# Run the application
run:
	$(GOCMD) run main.go

# Install dependencies
deps:
	$(GOGET) -v -t -d ./...

# Lint the code
lint:
	golangci-lint run

# Build for all platforms
build-all:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-linux-amd64
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BINARY_NAME)-linux-arm64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BINARY_NAME)-darwin-arm64
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-windows-amd64.exe
	GOOS=windows GOARCH=arm64 $(GOBUILD) -o $(BINARY_NAME)-windows-arm64.exe

# Run goreleaser
release:
	goreleaser release --rm-dist

# Show help
help:
	@echo "Available targets:"
	@echo "  all           - Run tests and build"
	@echo "  test          - Run all tests"
	@echo "  test-cli      - Run CLI tests"
	@echo "  test-buildpacks - Run buildpacks tests"
	@echo "  test-goreleaser - Run goreleaser tests"
	@echo "  build         - Build the binary"
	@echo "  clean         - Clean build files"
	@echo "  run           - Run the application"
	@echo "  deps          - Install dependencies"
	@echo "  lint          - Lint the code"
	@echo "  build-all     - Build for all platforms"
	@echo "  release       - Run goreleaser"
	@echo "  help          - Show this help message"