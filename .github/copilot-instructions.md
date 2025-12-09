# GitHub Copilot Instructions for Glyphic

This file provides context and guidelines for GitHub Copilot when assisting with the glyphic project.

## Project Overview

**Glyphic** is a quantum-resistant Diceware password generator built in Go 1.25.5. It generates secure, memorable passphrases with a distinctive Matrix-style "decode" reveal animation using sci-fi glyphs from the Matrix Code NFI font.

## Core Principles

1. **Security-First**: Cryptographic randomness (crypto/rand only), secure memory handling (mlock), zero persistence, constant-time operations
2. **Performance**: Capable of generating up to 1 billion passwords efficiently
3. **Modern Go**: Leverage Go 1.25.5 features (iter.Seq, slices, sync.OnceValue, structs package, etc.)
4. **Clean Architecture**: Strict layer separation (domain, application, infrastructure)
5. **Domain-Driven Design**: Rich domain models with explicit business rules
6. **CQRS**: Separate read and write operations where appropriate
7. **Terminal UX**: Beautiful Bubble Tea TUI with Lip Gloss styling

## Architecture Layers

```
cmd/glyphic/          # CLI entry point
pkg/version/          # Version management (semantic versioning with Git integration)
internal/
  ├── domain/         # Pure business logic (no dependencies)
  ├── application/    # Use cases, orchestration
  ├── infrastructure/ # External concerns (DB, filesystem, HTTP)
  ├── generator/      # Password generation logic
  ├── wordlist/       # Wordlist fetching, caching, management
  ├── font/           # Matrix Code NFI font handling
  ├── tui/            # Bubble Tea UI components
  └── security/       # Crypto primitives, memory locking
```

## Go 1.25.5 Features to Use

- **Range over integers**: `for i := range count { ... }`
- **Range over functions**: `iter.Seq[T]` and `iter.Seq2[K,V]` for custom iterators
- **slices package**: `slices.Sort()`, `slices.BinarySearch()`, `slices.Chunk()`, `slices.Compact()`
- **maps package**: `maps.Clone()`, `maps.Copy()`
- **cmp package**: `cmp.Or()`, `cmp.Clamp()`
- **log/slog**: Structured logging throughout
- **unique package**: `unique.Handle[string]` for interned strings (wordlist optimization)
- **structs package**: `structs.HostLayout` for secure buffer memory layout
- **sync.OnceValue/OnceFunc**: Lazy initialization
- **weak package**: Weak pointers for caches

## Security Requirements

### Memory Safety
- Use `unix.Mlock()` to prevent sensitive data from being swapped
- Zero all password buffers immediately after use
- Implement canary tokens for buffer overflow detection
- Use `structs.HostLayout` for predictable memory layout

### Randomness
- **ONLY** use `crypto/rand` for all random operations
- Validate PRNG on startup
- Never fall back to `math/rand`
- Implement entropy monitoring

### Code Patterns
```go
// Secure random selection
func secureRandomIndex(max int) (int, error) {
    if max <= 0 {
        return 0, errors.New("max must be positive")
    }
    var buf [8]byte
    if _, err := rand.Read(buf[:]); err != nil {
        return 0, fmt.Errorf("crypto/rand failed: %w", err)
    }
    n := binary.BigEndian.Uint64(buf[:])
    return int(n % uint64(max)), nil
}

// Secure zeroing
func SecureZero(data []byte) {
    for i := range data {
        data[i] = 0
    }
    runtime.KeepAlive(data)
}

// Memory locking
func LockMemory(data []byte) error {
    return unix.Mlock(data)
}

func UnlockMemory(data []byte) error {
    return unix.Munlock(data)
}
```

## Wordlist Management

- Auto-fetch from verified HTTPS sources on first run
- SHA-256 checksum verification for all wordlists
- Cache in `~/.local/share/glyphic/wordlists/`
- Always select words from **at least 3 different wordlists**
- Support custom wordlists via CLI flags
- Implement exclusion lists (profanity, confusing words, etc.)

## Testing Standards

- Use `testify/assert` for readable assertions
- Table-driven tests for multiple scenarios
- Use `gomock` for interface mocking
- Fuzz testing for parsers and generators
- `testing/synctest` for concurrent code (Go 1.24+)
- Aim for >80% coverage on critical paths

### Test Structure
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
        // More cases...
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

## TUI Implementation (Bubble Tea + Lip Gloss)

### Decode Animation
- Use Matrix Code NFI font glyphs for scramble effect
- Implement CLI-safe glyph validation
- Terminal detection (xterm-256color, kitty, alacritty = full Unicode; vt/linux = basic UTF-8; dumb/vt100 = ASCII fallback)
- Support multiple color schemes (matrix, cyber, fire, vapor, mono)
- Configurable speeds (slow, normal, fast)

### Bubble Tea Pattern
```go
type model struct {
    password    []byte
    revealed    []bool
    currentStep int
    opts        RevealOptions
}

func (m model) Init() tea.Cmd {
    return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "q" {
            return m, tea.Quit
        }
    case tickMsg:
        // Update reveal state
        return m, tick()
    }
    return m, nil
}

func (m model) View() string {
    return lipgloss.NewStyle().
        Foreground(lipgloss.Color("#00FF00")).
        Render(m.renderPassword())
}
```

## CLI Flags

Key flags to implement:
- `-n, --count`: Number of passwords (1 to 1B)
- `-w, --words`: Words per password (3-12)
- `-c, --capitalize`: Capitalization strategy
- `--numbers`: Add random digits
- `--special`: Add special characters
- `-C, --copy`: Copy to clipboard with auto-clear
- `--no-reveal`: Skip animation
- `--config`: Load settings from file
- `--verify`: Show entropy breakdown
- `-v, --version`: Show version info (uses pkg/version)
- `--version-json`: Version info as JSON

## Version Management

Use the `pkg/version` package for all version-related operations:

```go
import "github.com/youruser/glyphic/pkg/version"

func printVersion() {
    info := version.GetBuildInfo()
    fmt.Printf("glyphic %s\n", info.Version)
    fmt.Printf("Git Commit: %s\n", info.GitCommit)
    fmt.Printf("Build Time: %s\n", info.BuildTime)
    fmt.Printf("Go Version: %s\n", info.GoVersion)
}
```

Build with version injection:
```makefile
LDFLAGS := -ldflags "\
    -X github.com/youruser/glyphic/pkg/version.Version=$(VERSION) \
    -X github.com/youruser/glyphic/pkg/version.GitCommit=$(GIT_COMMIT) \
    -X github.com/youruser/glyphic/pkg/version.BuildTime=$(BUILD_TIME) \
    -X github.com/youruser/glyphic/pkg/version.GoVersion=$(GO_VERSION)"
```

## Code Style

### Naming Conventions
- PascalCase for exported identifiers
- camelCase for unexported identifiers
- ALL_CAPS for constants
- Use descriptive names: `userID` not `u`, `isValid` not `ok`

### Error Handling
- Always wrap errors with context: `fmt.Errorf("failed to load wordlist: %w", err)`
- Use typed errors for domain errors
- Never ignore errors
- Use `errors.Is()` and `errors.As()` for error checking

### Logging
- Use `slog` throughout
- Include relevant context in log entries
- Log at appropriate levels (Debug, Info, Warn, Error)
```go
slog.Info("wordlist loaded",
    slog.String("source", source.ID),
    slog.Int("word_count", len(words)),
    slog.String("checksum", checksum),
)
```

## Documentation

- Write godoc comments for all exported types and functions
- Include examples in comments where helpful
- Keep comments up-to-date with code changes
- Document "why" not "what" in inline comments

## Resources

### Project Documentation
- `INITIAL_PROMPT.md` - Complete project specification
- `CHARM_HELPER.md` - Charmbracelet ecosystem reference
- `docs/dev/VERSION.md` - Version management guide
- `docs/dev/software_principals.md` - Core engineering principles
- `docs/dev/testing_go.md` - Testing best practices
- `docs/dev/go_concurrency.md` - Concurrency patterns
- `docs/dev/security_first_repository_design.md` - Security patterns

### External References
- [Go 1.25 Release Notes](https://go.dev/doc/go1.25)
- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Lip Gloss Documentation](https://github.com/charmbracelet/lipgloss)
- [EFF Wordlists](https://www.eff.org/deeplinks/2016/07/new-wordlists-random-passphrases)

## Common Pitfalls to Avoid

1. **Never use `math/rand`** - only `crypto/rand`
2. **Don't persist passwords** - generate, display, zero
3. **Don't share memory between goroutines** - use channels
4. **Avoid premature optimization** - measure first
5. **Don't expose internal types** - use interfaces
6. **Never ignore context cancellation** - respect `ctx.Done()`
7. **Don't use `panic` for error handling** - return errors
8. **Avoid global mutable state** - use dependency injection

## When to Ask for Clarification

- Security-related decisions
- Breaking API changes
- Major architectural changes
- Performance trade-offs
- External dependency additions

## Quick Reference Commands

```bash
# Run tests
go test -v ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Build with version injection
make build

# Run linting
golangci-lint run

# Format code
go fmt ./...

# Run fuzzing
go test -fuzz=FuzzGeneratePassword -fuzztime=30s
```

---

Remember: This is a security-critical application. When in doubt, favor security and clarity over cleverness and performance.
