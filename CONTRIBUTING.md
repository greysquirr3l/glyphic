# Contributing to Glyphic

Thank you for your interest in contributing to Glyphic! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Code Standards](#code-standards)
- [Testing Requirements](#testing-requirements)
- [Security Guidelines](#security-guidelines)
- [Pull Request Process](#pull-request-process)
- [Release Process](#release-process)

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/1/code_of_conduct/). By participating, you are expected to uphold this code.

## Getting Started

### Prerequisites

- **Go 1.25.5+** (leverages latest features)
- **Make** (optional, for convenience)
- **golangci-lint** (for linting)
- **gosec** (for security scanning)

### Fork and Clone

```bash
# Fork the repository on GitHub
# Clone your fork
git clone https://github.com/YOUR_USERNAME/glyphic.git
cd glyphic

# Add upstream remote
git remote add upstream https://github.com/greysquirr3l/glyphic.git
```

### Build and Test

```bash
# Install dependencies
go mod download

# Build
make build

# Run tests
make test

# Run all checks
make check
```

## Development Workflow

### Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

### Make Your Changes

1. Write clean, idiomatic Go code
2. Follow existing code style and patterns
3. Add tests for new functionality
4. Update documentation as needed

### Commit Your Changes

Use conventional commit messages:

```bash
git commit -m "feat: add new color scheme"
git commit -m "fix: resolve entropy calculation bug"
git commit -m "docs: update installation instructions"
git commit -m "test: add wordlist manager tests"
git commit -m "refactor: improve generator performance"
```

**Commit message prefixes:**
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `test:` - Test additions/changes
- `refactor:` - Code refactoring
- `perf:` - Performance improvements
- `chore:` - Build process or auxiliary tool changes
- `style:` - Code style changes (formatting)

### Keep Your Branch Updated

```bash
git fetch upstream
git rebase upstream/main
```

## Code Standards

### Go Style Guide

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` and `gofmt -s` for formatting
- Follow the project's existing patterns

### Naming Conventions

```go
// Exported identifiers: PascalCase
type PasswordGenerator struct {}
func GeneratePassword() string {}

// Unexported identifiers: camelCase
type internalState struct {}
func parseOptions() {}

// Constants: ALL_CAPS or PascalCase (for exported)
const MaxWordCount = 12
const defaultTimeout = 30
```

### Package Organization

```
internal/          # Private application code
├── domain/       # Pure business logic (no external dependencies)
├── application/  # Use cases, orchestration
├── infrastructure/ # External concerns (HTTP, filesystem)
└── [feature]/    # Feature-specific packages
```

### Error Handling

```go
// Always wrap errors with context
if err != nil {
    return fmt.Errorf("failed to load wordlist: %w", err)
}

// Use custom error types for domain errors
var ErrInvalidWordCount = errors.New("word count must be between 3 and 12")

// Check errors with errors.Is() and errors.As()
if errors.Is(err, ErrInvalidWordCount) {
    // handle specific error
}
```

### Logging

```go
// Use slog throughout
slog.Info("wordlist loaded",
    slog.String("source", source.ID),
    slog.Int("word_count", len(words)),
)

slog.Error("failed to fetch wordlist",
    slog.String("url", url),
    slog.Any("error", err),
)
```

## Testing Requirements

### Test Coverage

- **Aim for >80% coverage** on critical paths
- Security and generator packages should have near 100% coverage
- Use table-driven tests for multiple scenarios

### Writing Tests

```go
func TestGeneratePassword(t *testing.T) {
    tests := []struct {
        name    string
        opts    Options
        want    int // expected word count
        wantErr bool
    }{
        {
            name: "default options",
            opts: DefaultOptions,
            want: 6,
            wantErr: false,
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Generate(tt.opts)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.NoError(t, err)
            assert.Len(t, got, tt.want)
        })
    }
}
```

### Test Commands

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run with race detector
go test -race ./...

# Run specific package tests
go test -v ./internal/generator

# Run benchmarks
make bench

# Run fuzz tests
make fuzz
```

## Security Guidelines

### Critical Rules

1. **NEVER use `math/rand`** - Only `crypto/rand` for randomness
2. **Always zero sensitive data** after use with `SecureZero()`
3. **Use memory locking** (`unix.Mlock()`) for password buffers
4. **No persistence** - Never write passwords to disk or logs
5. **Constant-time operations** where applicable

### Secure Coding Examples

```go
// ✅ CORRECT: Using crypto/rand
idx, err := security.SecureRandomIndex(len(list))
if err != nil {
    return fmt.Errorf("random selection failed: %w", err)
}

// ❌ WRONG: Never use math/rand
idx := rand.Intn(len(list)) // FORBIDDEN!

// ✅ CORRECT: Zero sensitive data
defer security.SecureZero(passwordBuffer)

// ✅ CORRECT: Lock memory
if err := security.LockMemory(buffer); err != nil {
    return err
}
defer security.UnlockMemory(buffer)
```

### Security Checklist

Before submitting security-sensitive code:

- [ ] Uses `crypto/rand` exclusively
- [ ] Sensitive data is zeroed after use
- [ ] Memory locking is applied where appropriate
- [ ] No password data in logs or error messages
- [ ] Constant-time comparisons for secrets
- [ ] Input validation prevents injection attacks
- [ ] HTTPS used for all external requests
- [ ] Checksums verified for downloaded files

## Pull Request Process

### Before Submitting

```bash
# Run all quality checks
make check

# Ensure tests pass
make test

# Verify build succeeds
make build

# Test the binary
./glyphic --version
./glyphic --no-reveal --quiet --words 6
```

### Pull Request Checklist

- [ ] Code follows project style guidelines
- [ ] All tests pass (`make test`)
- [ ] No lint warnings (`make lint`)
- [ ] Security scan clean (`make security`)
- [ ] Test coverage maintained or improved
- [ ] Documentation updated (README, godoc, inline comments)
- [ ] CHANGELOG.md updated (user-facing changes)
- [ ] Commit messages follow convention
- [ ] Branch is up-to-date with main
- [ ] No merge conflicts

### PR Description Template

```markdown
## Description
Brief description of the change

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
Describe how you tested your changes

## Checklist
- [ ] Tests pass locally
- [ ] Linting passes
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
```

### Review Process

1. **Automated checks** must pass (tests, lint, build)
2. **Code review** by at least one maintainer
3. **Security review** for security-sensitive changes
4. **Approval** required before merge
5. **Squash and merge** to keep history clean

## Release Process

### Version Numbering

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes (e.g., 2.0.0)
- **MINOR**: New features, backward compatible (e.g., 1.1.0)
- **PATCH**: Bug fixes, backward compatible (e.g., 1.0.1)

### Creating a Release

```bash
# Update VERSION file
echo "1.1.0" > VERSION

# Update CHANGELOG.md
# Add release notes under new version header

# Commit version bump
git commit -am "chore: bump version to 1.1.0"

# Create and push tag
git tag -a v1.1.0 -m "Release v1.1.0"
git push origin v1.1.0

# Build release binaries
make build-all
```

## Areas for Contribution

### High Priority

- **Clipboard integration** (`--copy` flag with auto-clear)
- **Config file support** (YAML/TOML)
- **Fuzz testing** suite expansion
- **i18n wordlist support** (Spanish, French, German)
- **Tab completion** (bash, zsh, fish)

### Medium Priority

- **Additional color schemes**
- **Custom glyph sets**
- **Animation customization** (reveal patterns)
- **Password strength analyzer**
- **Benchmark CI integration**

### Documentation

- **Tutorial videos**
- **Blog posts** about architecture/security
- **Translation** of README to other languages
- **Example scripts** and use cases
- **Performance tuning guide**

### Infrastructure

- **GitHub Actions** CI/CD pipeline
- **Docker container**
- **Homebrew formula**
- **Snap package**
- **Release automation**

## Questions?

- **Issues**: [GitHub Issues](https://github.com/greysquirr3l/glyphic/issues)
- **Discussions**: [GitHub Discussions](https://github.com/greysquirr3l/glyphic/discussions)
- **Security**: See [SECURITY.md](SECURITY.md)

---

Thank you for contributing to Glyphic! 🔐✨
