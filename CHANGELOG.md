# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive test suite for `internal/security` package (74.3% coverage)
- Comprehensive test suite for `internal/tui` package (74.3% coverage)
- CONTRIBUTING.md with detailed contribution guidelines
- Security package tests covering all cryptographic primitives
- TUI package tests covering Bubble Tea animation logic
- Benchmarks for security and TUI performance

### Fixed
- `SecureRandomFloat()` now correctly returns values in [0.0, 1.0) range
- Improved test coverage from 46.4% to 66.4% overall

### Changed
- Enhanced test coverage across all critical security paths
- Improved documentation for contributors

## [0.1.0] - 2025-12-09

### Added
- Initial release of glyphic password generator
- Quantum-resistant Diceware password generation
- Matrix Code NFI font integration with 500+ glyphs
- Terminal detection and graceful degradation (Dumb/Basic/Full modes)
- 7 color schemes (Matrix, Cyber, Fire, Vapor, Mono, Nord, Gruvbox)
- Bubble Tea reveal animation with scramble effect
- 5 capitalization modes (none/first/random/all/alternating)
- 5 separator modes (none/space/dash/underscore/custom)
- Optional numbers and special character insertion
- Wordlist management with EFF wordlist support
- SHA-256 checksum verification for wordlist integrity
- Exclusion lists (confusing words, profanity, slurs, sensitive terms)
- Entropy calculation and display
- Batch password generation (1 to 1 billion)
- Memory locking with unix.Mlock for security
- Secure zeroing of sensitive data
- PRNG validation on startup
- Version management with Git integration
- Comprehensive CLI with 40+ flags
- Full test suite (46.4% coverage initially)
- Security scanning with gosec
- MIT License

### Security
- Uses crypto/rand exclusively for all randomness
- Memory locking prevents swapping sensitive data to disk
- Secure zeroing of password buffers with runtime.KeepAlive
- Constant-time comparisons for sensitive data
- No password persistence or logging
- HTTPS-only for wordlist downloads

[Unreleased]: https://github.com/greysquirr3l/glyphic/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/greysquirr3l/glyphic/releases/tag/v0.1.0
