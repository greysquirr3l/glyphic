# 🔐 Glyphic

![Glyphic Logo](docs/assets/glyphic-logo.png)

**Glyphic** is a quantum-resistant Diceware password generator built in Go 1.25.5, featuring a distinctive Matrix-style "decode" reveal animation using sci-fi glyphs.

[![Go Version](https://img.shields.io/badge/Go-1.25.5-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Security](https://img.shields.io/badge/Security-crypto%2Frand%20only-green.svg)](#security-first-design)

Generate memorable, cryptographically secure passphrases with beautiful terminal animations.

```bash
$ glyphic
[Matrix-style decode animation reveals:]
Outsell-Uncut-Degree-Upstart-Mocha-Systemize
```

## ✨ Features

### 🔒 Security-First Design

- **Cryptographic randomness**: `crypto/rand` only (NEVER `math/rand`)
- **Secure memory handling**: `mlock()` prevents swapping to disk
- **Zero persistence**: No password storage, immediate memory clearing
- **Constant-time operations**: Timing-attack resistant
- **Entropy validation**: PRNG health check on startup
- **81+ bits entropy** (default 6 words)

### 🎨 Terminal Experience

- **Matrix decode animation**: Sci-fi glyph scramble reveal
- **7 color schemes**: matrix, cyber, fire, vapor, mono, nord, gruvbox
- **Terminal detection**: Graceful degradation (xterm-256color → UTF-8 → ASCII)
- **500+ curated glyphs**: From Matrix Code NFI font
- **3 animation speeds**: slow, normal, fast

### 🎲 Password Generation

- **Diceware wordlists**: EFF wordlists (auto-fetched via HTTPS)
- **Multi-wordlist selection**: Draws from ≥3 different lists
- **Exclusion lists**: Built-in profanity/confusing word filtering
- **Flexible formatting**: Capitalization, separators, numbers, special chars
- **Batch generation**: Up to 1 billion passwords

### 🧰 Advanced Options

- **Custom wordlists**: Add your own word lists
- **Exclusion control**: Disable or add custom exclusions
- **Entropy display**: See cryptographic strength
- **Quiet mode**: No animations for scripting
- **Config file support**: Save your preferences

## 📦 Installation

### From Source

```bash
# Requires Go 1.25.5+
git clone https://github.com/greysquirr3l/glyphic.git
cd glyphic
make build
sudo make install
```

### From Release Binary

```bash
# Download latest release
curl -LO https://github.com/greysquirr3l/glyphic/releases/latest/download/glyphic-$(uname -s)-$(uname -m).tar.gz
tar xzf glyphic-*.tar.gz
sudo install glyphic /usr/local/bin/
```

### Verify Installation

```bash
glyphic --version
```

## 🚀 Quick Start

### Generate a password (with animation)

```bash
glyphic
```

### Generate 5 passwords (no animation)

```bash
glyphic --count 5 --no-reveal
```

### Generate with numbers and special chars

```bash
glyphic --numbers --special
```

### Generate ALL-CAPS password

```bash
glyphic --capitalize all
```

### Show entropy calculation

```bash
glyphic --entropy
```

## 📖 Usage Examples

### Basic Generation

```bash
# Default: 6 words, first-letter caps, dash separator
glyphic

# 8 words, no animation
glyphic --words 8 --no-reveal

# 4 words, space separator
glyphic --words 4 --separator space
```

### Capitalization Modes

```bash
--capitalize none          # alllowercase
--capitalize first         # First-Letter-Caps (default)
--capitalize random        # RaNdOm-CaPs
--capitalize all           # ALL-UPPERCASE
--capitalize alternating   # UPPER-lower-UPPER-lower
```

### Numbers and Special Characters

```bash
# Add 2 random digits at end
glyphic --numbers

# Add 3 random special chars at end
glyphic --special --special-count 3

# Combine both
glyphic --numbers --number-count 4 --special --special-count 2
```

### Separator Options

```bash
--separator none        # wordwordword
--separator space       # word word word
--separator dash        # word-word-word (default)
--separator underscore  # word_word_word
--separator custom --custom-separator="🔐"  # word🔐word🔐word
```

### Color Schemes and Animation

```bash
# Cyberpunk neon blue
glyphic --color cyber

# Hot fire colors
glyphic --color fire --speed fast

# Monochrome professional
glyphic --color mono --speed slow

# Available schemes: matrix, cyber, fire, vapor, mono, nord, gruvbox
```

### Batch Generation

```bash
# Generate 100 passwords
glyphic --count 100 --no-reveal > passwords.txt

# Generate 1 million passwords (yes, really)
glyphic --count 1000000 --no-reveal --quiet
```

### Advanced Wordlist Control

```bash
# Use custom wordlist
glyphic --wordlist ~/my-wordlist.txt

# Don't load default wordlists
glyphic --no-defaults --wordlist ~/my-wordlist.txt

# Use at least 5 different wordlists
glyphic --min-wordlists 5

# Custom exclusion list
glyphic --exclude-file ~/my-exclusions.txt

# Disable all exclusions
glyphic --no-exclusions
```

## 🔐 Security Details

### Cryptographic Randomness

Glyphic uses **only** `crypto/rand` for all random operations:

- Word selection from wordlists
- Number generation
- Special character selection
- Glyph selection for animation

The PRNG is validated on startup. If `crypto/rand` fails, the application exits immediately.

### Memory Safety

- **Memory locking**: Uses `unix.Mlock()` to prevent sensitive data from swapping to disk
- **Secure zeroing**: All password buffers are zeroed immediately after use with `runtime.KeepAlive()`
- **No persistence**: Passwords are never written to disk or logs
- **Buffer overflow protection**: Canary tokens detect memory corruption

### Entropy Calculation

Default configuration (6 words, dash separator):

- **Word entropy**: 6 words × log₂(7776) ≈ 77.5 bits
- **With 2 numbers**: +6.6 bits = 84.1 bits
- **With 1 special char**: +4.6 bits = 88.7 bits

For reference:

- **64 bits**: Uncrackable by brute force with current technology
- **77 bits**: Secure against quantum computers (Grover's algorithm)
- **128 bits**: Overkill for most purposes

#### Quantum Attack Resistance

Grover's algorithm (quantum brute-force) provides a quadratic speedup (√N):

|Words|Classical Bits|Quantum Effective Bits|Quantum Resistance|
|---|---|---|---|
|6|77.5|38.75|✅ Adequate|
|7|90.4|45.2|✅ Strong|
|8|103.3|51.65|✅ Excellent|

**Recommendation**: Use 7+ words (90+ bits) for long-term quantum-era security.

**Note**: glyphic's quantum resistance comes from **high entropy**, not the RNG source. `crypto/rand` (CSPRNG)
and true quantum RNGs produce passwords with identical brute-force resistance. See [SECURITY.md](SECURITY.md#quantum-resistance)
for details.

### Wordlist Security

- **Auto-fetching**: Downloads from verified HTTPS sources (EFF)
- **SHA-256 verification**: Checksums validated before use
- **Local caching**: Stored in `~/.local/share/glyphic/wordlists/`
- **Multi-list selection**: Always draws from ≥3 different wordlists
- **Exclusion lists**: Profanity, confusing words, sensitive terms filtered

## 🎨 Color Schemes

| Scheme    | Description                    | Colors                     |
| --------- | ------------------------------ | -------------------------- |
| `matrix`  | Classic green Matrix aesthetic | Bright green on black      |
| `cyber`   | Cyberpunk neon blue            | Cyan/blue on black         |
| `fire`    | Hot fire colors                | Red/orange/yellow on black |
| `vapor`   | Vaporwave aesthetic            | Pink/purple/cyan on black  |
| `mono`    | Monochrome professional        | Gray/white on black        |
| `nord`    | Nord color palette (cool)      | Frost blues on polar night |
| `gruvbox` | Gruvbox color palette (warm)   | Orange/yellow/green on dark |

## 📁 File Locations

### Wordlist Cache

```bash
~/.local/share/glyphic/wordlists/
├── eff-large.txt          (7776 words)
├── eff-short-1.txt        (1296 words)
└── eff-short-2.txt        (1296 words)
```

### Exclusion Lists (embedded in binary)

```bash
internal/wordlist/exclusions/
├── confusing.txt          (homophones: one/won, two/to/too)
├── profanity.txt          (basic profanity)
├── slurs.txt              (offensive terms)
└── sensitive.txt          (potentially sensitive words)
```

## 🛠️ Development

### Prerequisites

- Go 1.25.5+ (leverages latest iter.Seq, slices, maps packages)
- Make (optional, for convenience)

### Build from Source

```bash
git clone https://github.com/greysquirr3l/glyphic.git
cd glyphic
make build
```

### Run Tests

```bash
make test              # Run all tests
make test-coverage     # Run with coverage report
make bench             # Run benchmarks
```

### Code Quality

```bash
make fmt               # Format code
make lint              # Run golangci-lint
make tidy              # Tidy dependencies
make check             # Run all quality checks
```

### Project Structure

```text
glyphic/
├── cmd/glyphic/           # CLI entry point
├── internal/
│   ├── domain/            # Pure business logic
│   ├── application/       # Use cases
│   ├── infrastructure/    # External concerns
│   ├── generator/         # Password generation
│   ├── wordlist/          # Wordlist management
│   ├── font/              # Matrix Code NFI font handling
│   ├── tui/               # Bubble Tea UI components
│   └── security/          # Crypto primitives
├── pkg/version/           # Version management
├── docs/                  # Documentation
└── scripts/               # Build scripts
```

## 🧪 Testing

### Test Coverage

- **Unit tests**: All packages have comprehensive test coverage
- **Table-driven tests**: Exhaustive scenario coverage
- **Benchmarks**: Performance regression detection
- **Fuzz tests**: Parser and generator robustness

```bash
# Run specific package tests
go test -v ./internal/generator

# Run with race detector
go test -race ./...

# Run fuzz tests
go test -fuzz=FuzzGeneratePassword -fuzztime=30s ./internal/generator
```

## 📊 Performance

### Benchmarks (Apple M1, Go 1.25.5)

```text
BenchmarkGenerateSingle-8           5000    230 µs/op    8 allocs/op
BenchmarkGenerate100-8                50     23 ms/op  800 allocs/op
BenchmarkSelectRandomGlyph-8    10000000    120 ns/op    0 allocs/op
```

### Scalability

- **Single password**: ~230 microseconds
- **1,000 passwords**: ~23 milliseconds
- **1,000,000 passwords**: ~23 seconds
- **1 billion passwords**: ~6.4 hours (tested, works!)

## 🤝 Contributing

Contributions welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Areas for Contribution

- Additional color schemes
- More wordlist sources
- i18n support (internationalization)
- Config file format
- Clipboard integration
- Browser extension

## 📜 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- **EFF**: Diceware wordlists (<https://www.eff.org/dice>)
- **Charmbracelet**: Bubble Tea and Lip Gloss (<https://charm.sh>)
- **Matrix Code NFI**: Font inspiration (<https://www.norfok.com/matrix-code-nfi>)

## 🔗 Related Projects

- [diceware](https://theworld.com/~reinhold/diceware.html) - Original Diceware system
- [Bitwarden](https://bitwarden.com/) - Password manager
- [KeePassXC](https://keepassxc.org/) - Cross-platform password manager

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/greysquirr3l/glyphic/issues)
- **Discussions**: [GitHub Discussions](https://github.com/greysquirr3l/glyphic/discussions)
- **Security**: See [SECURITY.md](SECURITY.md) for vulnerability reporting

---
