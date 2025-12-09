# Version Management

This document explains how to manage versioning in this Go project using the `VERSION` file and Git integration.

## Overview

The project uses a combination of:

- `VERSION` file for semantic versioning
- Git commit hash for build identification
- Automatic version injection during build time

## Go Implementation

### Reading Version from File

```go
package main

import (
    "embed"
    "fmt"
    "os/exec"
    "strings"
    "time"
)

//go:embed VERSION
var versionFile embed.FS

// BuildInfo contains version and build information
type BuildInfo struct {
    Version   string
    GitCommit string
    BuildTime string
    GoVersion string
}

// GetVersion reads the version from the embedded VERSION file
func GetVersion() (string, error) {
    data, err := versionFile.ReadFile("VERSION")
    if err != nil {
        return "", fmt.Errorf("failed to read VERSION file: %w", err)
    }
    return strings.TrimSpace(string(data)), nil
}

// GetGitCommit gets the current git commit hash
func GetGitCommit() (string, error) {
    cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("failed to get git commit: %w", err)
    }
    return strings.TrimSpace(string(output)), nil
}

// GetBuildInfo returns complete build information
func GetBuildInfo() BuildInfo {
    version, _ := GetVersion()
    gitCommit, _ := GetGitCommit()
    
    return BuildInfo{
        Version:   version,
        GitCommit: gitCommit,
        BuildTime: time.Now().UTC().Format(time.RFC3339),
        GoVersion: runtime.Version(),
    }
}

// Example usage
func main() {
    buildInfo := GetBuildInfo()
    fmt.Printf("Version: %s\n", buildInfo.Version)
    fmt.Printf("Git Commit: %s\n", buildInfo.GitCommit)
    fmt.Printf("Build Time: %s\n", buildInfo.BuildTime)
    fmt.Printf("Go Version: %s\n", buildInfo.GoVersion)
}
```

### Using ldflags for Build-Time Injection

```go
package main

import (
    "fmt"
    "runtime"
)

// These variables will be set at build time using -ldflags
var (
    Version   = "dev"
    GitCommit = "unknown"
    BuildTime = "unknown"
    GoVersion = runtime.Version()
)

// BuildInfo contains version and build information
type BuildInfo struct {
    Version   string `json:"version"`
    GitCommit string `json:"git_commit"`
    BuildTime string `json:"build_time"`
    GoVersion string `json:"go_version"`
}

// GetBuildInfo returns the build information
func GetBuildInfo() BuildInfo {
    return BuildInfo{
        Version:   Version,
        GitCommit: GitCommit,
        BuildTime: BuildTime,
        GoVersion: GoVersion,
    }
}

// PrintVersion prints version information
func PrintVersion() {
    info := GetBuildInfo()
    fmt.Printf("Version:    %s\n", info.Version)
    fmt.Printf("Git Commit: %s\n", info.GitCommit)
    fmt.Printf("Build Time: %s\n", info.BuildTime)
    fmt.Printf("Go Version: %s\n", info.GoVersion)
}

func main() {
    PrintVersion()
}
```

## Makefile Implementation

```makefile
# VERSION.md Makefile Examples

# Variables
VERSION := $(shell cat VERSION)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION := $(shell go version | cut -d' ' -f3)

# Build output directory
BUILD_DIR := ./bin
BINARY_NAME := your-app

# Ldflags for version injection
LDFLAGS := -ldflags "\
    -X main.Version=$(VERSION) \
    -X main.GitCommit=$(GIT_COMMIT) \
    -X main.BuildTime=$(BUILD_TIME) \
    -X main.GoVersion=$(GO_VERSION)"

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
release: bump-version build-all
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
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	golangci-lint run

# Tidy dependencies
.PHONY: tidy
tidy:
	go mod tidy

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  build-all     - Build for all platforms"
	@echo "  dev           - Build development version (faster)"
	@echo "  install       - Install to GOPATH/bin"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  tidy          - Tidy dependencies"
	@echo "  version       - Show version information"
	@echo "  bump-version  - Bump version (bump=patch|minor|major)"
	@echo "  release       - Create a release"
	@echo "  clean         - Clean build artifacts"
	@echo "  help          - Show this help"
```

## Advanced Version Management

### Automatic Version from Git Tags

```go
// GetVersionFromGit gets version from git tags
func GetVersionFromGit() (string, error) {
    // Try to get version from git tag
    cmd := exec.Command("git", "describe", "--tags", "--exact-match", "HEAD")
    if output, err := cmd.Output(); err == nil {
        return strings.TrimPrefix(strings.TrimSpace(string(output)), "v"), nil
    }
    
    // Fallback to VERSION file
    return GetVersion()
}

// GetSemanticVersion returns a semantic version with git info
func GetSemanticVersion() (string, error) {
    version, err := GetVersion()
    if err != nil {
        return "", err
    }
    
    // Check if we're on a tag
    cmd := exec.Command("git", "describe", "--tags", "--exact-match", "HEAD")
    if _, err := cmd.Output(); err == nil {
        return version, nil // Clean version on tag
    }
    
    // Add commit info for non-tag builds
    gitCommit, _ := GetGitCommit()
    return fmt.Sprintf("%s-dev+%s", version, gitCommit), nil
}
```

### Makefile with Git Tag Integration

```makefile
# Get version from git tag or VERSION file
VERSION := $(shell git describe --tags --exact-match HEAD 2>/dev/null | sed 's/^v//' || cat VERSION)
IS_TAG := $(shell git describe --tags --exact-match HEAD 2>/dev/null && echo "true" || echo "false")

# Add -dev suffix if not on a tag
ifeq ($(IS_TAG),false)
    VERSION := $(VERSION)-dev
endif
```

## Usage Examples

### In CLI Applications

```go
package main

import (
    "flag"
    "fmt"
    "os"
)

var (
    Version   = "dev"
    GitCommit = "unknown"
    BuildTime = "unknown"
)

func main() {
    var showVersion = flag.Bool("version", false, "Show version information")
    flag.Parse()
    
    if *showVersion {
        fmt.Printf("%s version %s (commit %s, built %s)\n", 
            os.Args[0], Version, GitCommit, BuildTime)
        os.Exit(0)
    }
    
    // Your application logic here
}
```

### In Web Applications

```go
package main

import (
    "encoding/json"
    "net/http"
)

var (
    Version   = "dev"
    GitCommit = "unknown"
    BuildTime = "unknown"
)

func versionHandler(w http.ResponseWriter, r *http.Request) {
    info := map[string]string{
        "version":    Version,
        "git_commit": GitCommit,
        "build_time": BuildTime,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(info)
}

func main() {
    http.HandleFunc("/version", versionHandler)
    http.ListenAndServe(":8080", nil)
}
```

## Best Practices

1. **Always embed VERSION file** using `//go:embed` for consistent versioning
2. **Use semantic versioning** (MAJOR.MINOR.PATCH)
3. **Inject build info at compile time** using `-ldflags`
4. **Tag releases** in git for clean version tracking
5. **Include git commit** for development builds
6. **Automate version bumping** in your CI/CD pipeline
7. **Expose version info** in CLI and web applications

## Integration with CI/CD

```yaml
# GitHub Actions example
name: Build and Release
on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      
      - name: Build
        run: make build-all
        
      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: bin/*
```

This approach provides a robust version management system that works both in
development and production environments.
