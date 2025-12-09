# Code Style and Conventions for Glyphic

## Go Style Guidelines

- Follow standard Go conventions (PascalCase for exports, camelCase for unexported)
- Use meaningful variable names: `userID` not `u`, `isValid` not `ok`
- ALL_CAPS for constants
- Write godoc comments for all exported types and functions
- Keep functions focused and small (single responsibility)

## Modern Go 1.25.5 Features to Use

- Range over integers: `for i := range count { ... }`
- Range over functions: `iter.Seq[T]` and `iter.Seq2[K,V]` for custom iterators
- slices package: `slices.Sort()`, `slices.BinarySearch()`, `slices.Chunk()`, `slices.Compact()`
- maps package: `maps.Clone()`, `maps.Copy()`
- cmp package: `cmp.Or()`, `cmp.Clamp()`
- log/slog: Structured logging throughout
- unique package: `unique.Handle[string]` for interned strings
- structs package: `structs.HostLayout` for secure buffer memory layout
- sync.OnceValue/OnceFunc: Lazy initialization
- weak package: Weak pointers for caches

## Security Requirements

### CRITICAL: Only use crypto/rand for all random operations

```go
// ✅ Correct
import "crypto/rand"
var buf [32]byte
rand.Read(buf[:])

// ❌ NEVER use math/rand
import "math/rand"  // FORBIDDEN
```

### Memory Safety

- Use `unix.Mlock()` to prevent sensitive data from being swapped
- Zero all password buffers immediately after use with `SecureZero()`
- Implement canary tokens for buffer overflow detection
- Use `structs.HostLayout` for predictable memory layout

## Error Handling

- Always wrap errors with context: `fmt.Errorf("failed to load wordlist: %w", err)`
- Use typed errors for domain errors
- Never ignore errors
- Use `errors.Is()` and `errors.As()` for error checking

## Testing Standards

- Use testify/assert for readable assertions
- Table-driven tests for multiple scenarios
- Use gomock for interface mocking
- Fuzz testing for parsers and generators
- testing/synctest for concurrent code
- Aim for >80% coverage on critical paths

## Architecture Patterns

- Clean Architecture: Strict layer separation (domain, application, infrastructure)
- Domain-Driven Design: Rich domain models with explicit business rules
- CQRS: Separate read and write operations where appropriate
- Dependency injection: Use interfaces for testability

## Logging

- Use slog throughout
- Include relevant context in log entries
- Log at appropriate levels (Debug, Info, Warn, Error)

```go
slog.Info("wordlist loaded",
    slog.String("source", source.ID),
    slog.Int("word_count", len(words)),
    slog.String("checksum", checksum),
)
```

## Common Pitfalls to Avoid

1. Never use `math/rand` - only `crypto/rand`
2. Don't persist passwords - generate, display, zero
3. Don't share memory between goroutines - use channels
4. Avoid premature optimization - measure first
5. Don't expose internal types - use interfaces
6. Never ignore context cancellation - respect `ctx.Done()`
7. Don't use `panic` for error handling - return errors
8. Avoid global mutable state - use dependency injection
