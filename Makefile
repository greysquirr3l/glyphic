# Glyphic Makefile
# Quantum-Resistant Diceware Password Generator

# Variables
VERSION := $(shell cat VERSION)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION := $(shell go version | cut -d' ' -f3)

# Build output directory
BUILD_DIR := ./bin
BINARY_NAME := glyphic

# Ldflags for version injection
LDFLAGS := -ldflags "\
	-X github.com/greysquirr3l/glyphic/pkg/version.GitCommit=$(GIT_COMMIT) \
	-X github.com/greysquirr3l/glyphic/pkg/version.BuildTime=$(BUILD_TIME)"

# Default target
.PHONY: all
all: clean build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME) v$(VERSION) ($(GIT_COMMIT))"
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/$(BINARY_NAME)

# Build for multiple platforms
.PHONY: build-all
build-all: build-linux build-darwin build-windows

.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/$(BINARY_NAME)
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/$(BINARY_NAME)

.PHONY: build-darwin
build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/$(BINARY_NAME)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/$(BINARY_NAME)

.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/$(BINARY_NAME)

# Version management targets
.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Git Branch: $(GIT_BRANCH)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Version: $(GO_VERSION)"

# Bump version (requires 'bump' parameter: make bump-version bump=patch|minor|major)
.PHONY: bump-version
bump-version:
	@if [ "$(bump)" = "" ]; then \
		echo "Usage: make bump-version bump=patch|minor|major"; \
		exit 1; \
	fi
	@echo "Current version: $(VERSION)"
	@python3 -c "\
import sys; \
parts = '$(VERSION)'.split('.'); \
major, minor, patch = int(parts[0]), int(parts[1]), int(parts[2]); \
if '$(bump)' == 'major': major += 1; minor = 0; patch = 0; \
elif '$(bump)' == 'minor': minor += 1; patch = 0; \
elif '$(bump)' == 'patch': patch += 1; \
else: print('Invalid bump type'); sys.exit(1); \
print(f'{major}.{minor}.{patch}')" > VERSION
	@echo "New version: $$(cat VERSION)"
	@git add VERSION
	@git commit -m "Bump version to $$(cat VERSION)"
	@git tag "v$$(cat VERSION)"

# Create a release
.PHONY: release
release: test lint build-all
	@echo "Creating release v$(VERSION)"
	@echo "Built binaries:"
	@ls -la $(BUILD_DIR)/

# Development build (faster, no version injection)
.PHONY: dev
dev:
	@echo "Building development version..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/$(BINARY_NAME)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

# Install the binary to GOPATH/bin
.PHONY: install
install:
	go install $(LDFLAGS) ./cmd/$(BINARY_NAME)

# Run tests
.PHONY: test
test:
	go test -v -race ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
.PHONY: bench
bench:
	go test -bench=. -benchmem ./...

# Run fuzzing
.PHONY: fuzz
fuzz:
	go test -fuzz=FuzzGenerate -fuzztime=30s ./internal/generator

# Format code
.PHONY: fmt
fmt:
	go fmt ./...
	gofmt -s -w .

# Lint code
.PHONY: lint
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Run: brew install golangci-lint" && exit 1)
	golangci-lint run

# Tidy dependencies
.PHONY: tidy
tidy:
	go mod tidy
	go mod verify

# Security scan
.PHONY: security
security:
	@which gosec > /dev/null || (echo "gosec not installed. Run: go install github.com/securego/gosec/v2/cmd/gosec@latest" && exit 1)
	gosec -quiet ./...

# Run all checks (test, lint, security)
.PHONY: check
check: fmt tidy lint test security
	@echo "All checks passed!"

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  build-all     - Build for all platforms"
	@echo "  dev           - Build development version (faster)"
	@echo "  install       - Install to GOPATH/bin"
	@echo "  test          - Run tests with race detector"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  bench         - Run benchmarks"
	@echo "  fuzz          - Run fuzz tests"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  tidy          - Tidy dependencies"
	@echo "  security      - Run security scan"
	@echo "  check         - Run all checks (fmt, tidy, lint, test, security)"
	@echo "  version       - Show version information"
	@echo "  bump-version  - Bump version (bump=patch|minor|major)"
	@echo "  release       - Create a release build"
	@echo "  clean         - Clean build artifacts"
	@echo "  help          - Show this help"
