# Suggested Commands for Glyphic Development

## Build Commands

```bash
make build           # Build the application
make build-all       # Build for all platforms (Linux, macOS, Windows)
make dev             # Quick development build (no version injection)
make install         # Install to GOPATH/bin
make clean           # Clean build artifacts
```

## Testing Commands

```bash
make test            # Run tests with race detector
make test-coverage   # Run tests with coverage report (generates coverage.html)
make bench           # Run benchmarks
make fuzz            # Run fuzz tests (30s timeout)
```

## Code Quality Commands

```bash
make fmt             # Format code (go fmt + gofmt -s)
make lint            # Lint code (requires golangci-lint)
make tidy            # Tidy and verify dependencies
make security        # Run security scan (requires gosec)
make check           # Run all checks (fmt, tidy, lint, test, security)
```

## Version Management

```bash
make version                      # Show version information
make bump-version bump=patch      # Bump patch version
make bump-version bump=minor      # Bump minor version
make bump-version bump=major      # Bump major version
make release                      # Create a release build
```

## Direct Go Commands

```bash
go test -v ./...                  # Run all tests verbosely
go test -race ./...               # Run tests with race detector
go build ./cmd/glyphic            # Build main binary
go mod tidy                       # Tidy dependencies
go vet ./...                      # Run go vet
```

## macOS-Specific Commands

```bash
# Install tools
brew install golangci-lint        # Install linter
go install github.com/securego/gosec/v2/cmd/gosec@latest  # Install security scanner

# System commands
ls -la                            # List files with details
find . -name "*.go"               # Find Go files
grep -r "pattern" .               # Search for pattern
```

## Project Commands

```bash
# Run the application
./bin/glyphic --help              # Show help
./bin/glyphic                     # Generate password with defaults
./bin/glyphic -n 10 -w 8          # Generate 10 passwords with 8 words each
```
