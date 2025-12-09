# Glyphic Project Overview

## Purpose

Glyphic is a quantum-resistant Diceware password generator built in Go 1.25.5. It generates secure, memorable passphrases
with a distinctive Matrix-style "decode" reveal animation using sci-fi glyphs from the Matrix Code NFI font.

## Tech Stack

- **Language**: Go 1.25.5
- **Module Path**: github.com/greysquirr3l/glyphic
- **TUI Framework**: Bubble Tea + Lip Gloss (Charmbracelet ecosystem)
- **Testing**: testify/assert, gomock, fuzz testing
- **Security**: crypto/rand, unix.Mlock, secure memory handling

## Key Principles

1. **Security-First**: Cryptographic randomness (crypto/rand only), secure memory handling (mlock), zero persistence, constant-time operations
2. **Performance**: Capable of generating up to 1 billion passwords efficiently
3. **Modern Go**: Leverage Go 1.25.5 features (iter.Seq, slices, sync.OnceValue, structs package, etc.)
4. **Clean Architecture**: Strict layer separation (domain, application, infrastructure)
5. **Domain-Driven Design**: Rich domain models with explicit business rules
6. **CQRS**: Separate read and write operations where appropriate
7. **Terminal UX**: Beautiful Bubble Tea TUI with Lip Gloss styling

## Project Structure

```text
glyphic/
├── VERSION                           # Semantic version (e.g., 0.1.0)
├── Makefile                          # Build with version injection
├── cmd/glyphic/                      # CLI entry point
├── pkg/version/                      # Version management with Git integration
├── internal/
│   ├── security/                     # Crypto primitives, memory locking
│   ├── wordlist/                     # Wordlist fetching, caching, management
│   ├── generator/                    # Password generation logic
│   ├── font/                         # Matrix Code NFI font handling
│   ├── tui/                          # Bubble Tea UI components
│   ├── config/                       # Configuration management
│   └── clipboard/                    # Clipboard integration
├── fonts/                            # MatrixCodeNfi-YPPj.otf
└── docs/dev/                         # Development documentation
```

## External Dependencies

- golang.org/x/sys (for unix.Mlock)
- github.com/stretchr/testify (testing)
- Charmbracelet packages (TUI - to be added)
