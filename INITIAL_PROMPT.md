# Project: Quantum-Resistant Diceware Password Generator

## Codename: `glyphic`

---

## Overview

Build a secure, high-performance Go command-line utility that generates XKCD/Diceware-style passphrases with
quantum-resistant entropy levels. The tool features a distinctive "decode" text reveal animation using sci-fi glyphs
(from the included `lmym.ttf` font) that creates a cinematic terminal experience—characters appear scrambled with
alien-looking symbols before resolving to the actual password.

**Key Principles:**

- Security-first: Cryptographic randomness, secure memory handling, zero persistence
- Performance: Capable of generating up to 1 billion passwords efficiently
- Style: Terminal eye-candy with the scramble→reveal text effect
- Flexibility: Extensive customization options for password composition
- Modern Go: Leverage Go 1.25.5 features throughout
- Self-contained: Auto-fetches wordlists from verified sources, no system dependencies
- Versioned: Semantic versioning with build-time Git integration via pkg/version

**Project Structure:**

```text
glyphic/
├── VERSION                           # Semantic version (e.g., 1.0.0)
├── Makefile                          # Build with version injection
├── cmd/
│   └── glyphic/
│       └── main.go                   # Entry point with version handling
├── pkg/
│   └── version/
│       ├── version.go                # Version package with BuildInfo
│       └── version_test.go           # Version tests
├── internal/
│   ├── generator/                    # Password generation
│   ├── wordlist/                     # Wordlist management
│   ├── font/                         # Matrix Code NFI font handling
│   ├── tui/                          # Decode animation
│   └── security/                     # Crypto primitives
├── fonts/
│   └── MatrixCodeNfi-YPPj.otf        # Embedded font
└── docs/
    └── dev/
        └── VERSION.md                # Version management documentation
```

---

## Go 1.25.5 Features to Utilize

This project should take full advantage of modern Go capabilities:

### Language Features

```go
// Range over integers (Go 1.22+)
for i := range wordCount {
    words[i] = selectRandomWord()
}

// Range over functions / iterators (Go 1.23+)
func (wl *Wordlist) All() iter.Seq[string] {
    return func(yield func(string) bool) {
        for _, word := range wl.words {
            if !yield(word) {
                return
            }
        }
    }
}

// Use in range
for word := range wordlist.All() {
    process(word)
}

// Indexed iterator
func (wl *Wordlist) Enumerate() iter.Seq2[int, string] {
    return func(yield func(int, string) bool) {
        for i, word := range wl.words {
            if !yield(i, word) {
                return
            }
        }
    }
}
```

### Standard Library Improvements

```go
// Enhanced slices package (Go 1.21+)
import "slices"

// Sort wordlist
slices.Sort(wordlist)

// Binary search for validation
_, found := slices.BinarySearch(wordlist, word)

// Chunk processing for batch generation
for chunk := range slices.Chunk(passwords, batchSize) {
    processChunk(chunk)
}

// maps package for config handling
import "maps"

cfg := maps.Clone(defaultConfig)
maps.Copy(cfg, userConfig)

// cmp package for comparisons
import "cmp"

entropy := cmp.Or(userEntropy, defaultEntropy)
wordCount := cmp.Clamp(requested, minWords, maxWords)

// log/slog for structured logging (Go 1.21+)
import "log/slog"

slog.Info("wordlist fetched",
    slog.String("source", source.Name),
    slog.Int("word_count", len(words)),
    slog.String("sha256", checksum),
)

// crypto/rand improvements - use Read directly into any slice
var buf [32]byte
crypto/rand.Read(buf[:])

// unique package (Go 1.23+) for interned strings (wordlist optimization)
import "unique"

type Wordlist struct {
    handles []unique.Handle[string]
}

// structs package (Go 1.24+) for host layout control
import "structs"

type SecureBuffer struct {
    _      structs.HostLayout
    data   []byte
    locked bool
}
```

### Testing Improvements

```go
// testing/synctest for concurrent tests (Go 1.24+)
import "testing/synctest"

func TestConcurrentGeneration(t *testing.T) {
    synctest.Run(func() {
        // Deterministic concurrent test
    })
}

// Enhanced fuzzing
func FuzzPassphraseGeneration(f *testing.F) {
    f.Add(6, true, 1, 1) // wordCount, capitalize, numbers, specials
    f.Fuzz(func(t *testing.T, words int, cap bool, nums, specs int) {
        words = cmp.Clamp(words, 4, 10)
        nums = cmp.Clamp(nums, 0, 5)
        specs = cmp.Clamp(specs, 0, 5)
        
        result, err := Generate(Options{
            WordCount: words,
            Capitalize: cap,
            Numbers: nums,
            Specials: specs,
        })
        if err != nil {
            t.Skip() // Invalid combination
        }
        if len(result) == 0 {
            t.Error("empty result")
        }
    })
}
```

### Performance Features

```go
// sync.OnceValue/OnceValues for lazy initialization (Go 1.21+)
var loadDefaultExclusions = sync.OnceValue(func() []string {
    data, _ := exclusionsFS.ReadFile("exclusions/default.txt")
    return strings.Split(string(data), "\n")
})

// sync.OnceFunc for one-time operations
var initGlyphs = sync.OnceFunc(func() {
    ScrambleGlyphs = parseFont(lmymFont)
})

// context.WithoutCancel for background operations (Go 1.21+)
ctx := context.WithoutCancel(parentCtx)

// weak package for caches (Go 1.24+)
import "weak"

type WordCache struct {
    entries map[string]weak.Pointer[[]byte]
}
```

---

## Wordlist Management System

### Philosophy

Glyphic automatically fetches wordlists from **verified, authoritative sources** on first run and caches them locally.
Passwords are built by **randomly selecting words from at least 3 different wordlists**, increasing unpredictability
beyond single-source approaches.

### Verified Wordlist Sources

All wordlists are fetched via HTTPS with SHA-256 verification:

```go
package wordlist

import (
    "crypto/sha256"
    "encoding/hex"
)

type WordlistSource struct {
    ID          string   // Unique identifier
    Name        string   // Human-readable name
    URL         string   // HTTPS download URL
    SHA256      string   // Expected checksum (empty = fetch and trust on first run, warn user)
    MinWords    int      // Minimum expected word count
    License     string   // License information
    Description string   // What this list contains
    Enabled     bool     // Whether to use by default
    Category    string   // "general", "technical", "nature", "phonetic", etc.
}

// Verified wordlist sources - authoritative, stable URLs
var DefaultSources = []WordlistSource{
    {
        ID:          "eff-large",
        Name:        "EFF Large Wordlist",
        URL:         "https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt",
        SHA256:      "b68d7f5a4d26c9f0a5f08c9e5e2d6a3b9c1e4f7a8d2b5c8e1f4a7b0d3c6e9f2a5", // Replace with actual
        MinWords:    7776,
        License:     "CC-BY-3.0",
        Description: "EFF's curated list optimized for memorability and typing",
        Enabled:     true,
        Category:    "general",
    },
    {
        ID:          "eff-short-1",
        Name:        "EFF Short Wordlist 1",
        URL:         "https://www.eff.org/files/2016/09/08/eff_short_wordlist_1.txt",
        SHA256:      "a1b2c3d4e5f6...", // Replace with actual
        MinWords:    1296,
        License:     "CC-BY-3.0",
        Description: "Shorter words, 4-5 characters average",
        Enabled:     true,
        Category:    "general",
    },
    {
        ID:          "eff-short-2-prefix",
        Name:        "EFF Short Wordlist 2 (Unique Prefixes)",
        URL:         "https://www.eff.org/files/2016/09/08/eff_short_wordlist_2_0.txt",
        SHA256:      "d4e5f6a7b8c9...", // Replace with actual
        MinWords:    1296,
        License:     "CC-BY-3.0",
        Description: "Each word has unique 3-character prefix for autocomplete",
        Enabled:     true,
        Category:    "general",
    },
    {
        ID:          "securedrop",
        Name:        "SecureDrop Wordlist",
        URL:         "https://raw.githubusercontent.com/freedomofpress/securedrop/develop/securedrop/wordlists/en.txt",
        SHA256:      "e5f6a7b8c9d0...", // Replace with actual
        MinWords:    6800,
        License:     "AGPL-3.0",
        Description: "Freedom of the Press Foundation's SecureDrop passphrase list",
        Enabled:     true,
        Category:    "general",
    },
    {
        ID:          "bip39-english",
        Name:        "BIP-39 English",
        URL:         "https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/english.txt",
        SHA256:      "2f5ede7c9e1a6...", // Replace with actual - this one is critical
        MinWords:    2048,
        License:     "MIT",
        Description: "Bitcoin Improvement Proposal 39 mnemonic wordlist",
        Enabled:     true,
        Category:    "crypto",
    },
    {
        ID:          "orchard-street-medium",
        Name:        "Orchard Street Medium",
        URL:         "https://raw.githubusercontent.com/sts10/orchard-street-wordlists/main/lists/orchard-street-medium.txt",
        SHA256:      "f6a7b8c9d0e1...", // Replace with actual
        MinWords:    8192,
        License:     "MIT",
        Description: "Curated for passphrase generation, uniquely decodable",
        Enabled:     true,
        Category:    "general",
    },
    {
        ID:          "orchard-street-long",
        Name:        "Orchard Street Long",
        URL:         "https://raw.githubusercontent.com/sts10/orchard-street-wordlists/main/lists/orchard-street-long.txt",
        SHA256:      "a7b8c9d0e1f2...", // Replace with actual
        MinWords:    17576,
        License:     "MIT",
        Description: "Larger list for maximum entropy per word",
        Enabled:     false, // Opt-in for advanced users
        Category:    "general",
    },
    {
        ID:          "monero",
        Name:        "Monero Wordlist",
        URL:         "https://raw.githubusercontent.com/monero-project/monero/master/src/mnemonics/english.h",
        SHA256:      "b8c9d0e1f2a3...", // Replace with actual
        MinWords:    1626,
        License:     "BSD-3-Clause",
        Description: "Monero cryptocurrency mnemonic seed words",
        Enabled:     false, // Needs parser for .h format
        Category:    "crypto",
    },
    {
        ID:          "pgp-even",
        Name:        "PGP Even Words",
        URL:         "https://raw.githubusercontent.com/artob/pgp-wordlist/master/even.txt",
        SHA256:      "c9d0e1f2a3b4...", // Replace with actual
        MinWords:    256,
        License:     "Public Domain",
        Description: "PGP word list (even bytes) - phonetically distinct",
        Enabled:     true,
        Category:    "phonetic",
    },
    {
        ID:          "pgp-odd",
        Name:        "PGP Odd Words",
        URL:         "https://raw.githubusercontent.com/artob/pgp-wordlist/master/odd.txt",
        SHA256:      "d0e1f2a3b4c5...", // Replace with actual
        MinWords:    256,
        License:     "Public Domain",
        Description: "PGP word list (odd bytes) - phonetically distinct",
        Enabled:     true,
        Category:    "phonetic",
    },
}
```

### Wordlist Manager

```go
package wordlist

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "fmt"
    "io"
    "log/slog"
    "net/http"
    "os"
    "path/filepath"
    "slices"
    "strings"
    "sync"
    "time"
)

var (
    ErrChecksumMismatch = errors.New("wordlist checksum mismatch")
    ErrTooFewWords      = errors.New("wordlist has fewer words than expected")
    ErrFetchFailed      = errors.New("failed to fetch wordlist")
    ErrInsufficientLists = errors.New("insufficient wordlists available")
)

type Manager struct {
    cacheDir    string
    sources     []WordlistSource
    exclusions  *ExclusionList
    httpClient  *http.Client
    mu          sync.RWMutex
    loaded      map[string]*Wordlist
}

func NewManager(cacheDir string) (*Manager, error) {
    // Default: ~/.local/share/glyphic/wordlists/
    if cacheDir == "" {
        home, err := os.UserHomeDir()
        if err != nil {
            return nil, err
        }
        cacheDir = filepath.Join(home, ".local", "share", "glyphic", "wordlists")
    }
    
    if err := os.MkdirAll(cacheDir, 0700); err != nil {
        return nil, err
    }
    
    return &Manager{
        cacheDir: cacheDir,
        sources:  slices.Clone(DefaultSources),
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
            Transport: &http.Transport{
                MaxIdleConns:       10,
                IdleConnTimeout:    90 * time.Second,
                DisableCompression: false,
            },
        },
        loaded: make(map[string]*Wordlist),
    }, nil
}

// EnsureWordlists fetches any missing wordlists (called on first run)
func (m *Manager) EnsureWordlists(ctx context.Context) error {
    var errs []error
    
    for _, source := range m.sources {
        if !source.Enabled {
            continue
        }
        
        cachePath := m.cachePath(source.ID)
        
        if m.isValidCache(cachePath, source) {
            slog.Debug("wordlist cached", slog.String("id", source.ID))
            continue
        }
        
        slog.Info("fetching wordlist",
            slog.String("id", source.ID),
            slog.String("url", source.URL),
        )
        
        if err := m.fetchAndCache(ctx, source); err != nil {
            errs = append(errs, fmt.Errorf("%s: %w", source.ID, err))
            slog.Error("fetch failed",
                slog.String("id", source.ID),
                slog.String("error", err.Error()),
            )
        }
    }
    
    // Require at least 3 wordlists
    available := m.AvailableCount()
    if available < 3 {
        return fmt.Errorf("%w: need 3, have %d: %w",
            ErrInsufficientLists, available, errors.Join(errs...))
    }
    
    if len(errs) > 0 {
        slog.Warn("some wordlists unavailable",
            slog.Int("failed", len(errs)),
            slog.Int("available", available),
        )
    }
    
    return nil
}

func (m *Manager) fetchAndCache(ctx context.Context, source WordlistSource) error {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, source.URL, nil)
    if err != nil {
        return err
    }
    req.Header.Set("User-Agent", "glyphic/1.0 (password generator; +https://github.com/youruser/glyphic)")
    
    resp, err := m.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("%w: %v", ErrFetchFailed, err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("%w: HTTP %d", ErrFetchFailed, resp.StatusCode)
    }
    
    // Read with 10MB limit
    data, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
    if err != nil {
        return err
    }
    
    // Verify checksum
    hash := sha256.Sum256(data)
    actualHash := hex.EncodeToString(hash[:])
    
    if source.SHA256 != "" {
        if actualHash != source.SHA256 {
            return fmt.Errorf("%w: expected %s, got %s",
                ErrChecksumMismatch, source.SHA256, actualHash)
        }
    } else {
        slog.Warn("no checksum for wordlist - trust on first use",
            slog.String("id", source.ID),
            slog.String("sha256", actualHash),
        )
    }
    
    // Validate word count
    words := parseWordlist(data, source.ID)
    if len(words) < source.MinWords {
        return fmt.Errorf("%w: expected %d+, got %d",
            ErrTooFewWords, source.MinWords, len(words))
    }
    
    // Cache the file
    cachePath := m.cachePath(source.ID)
    if err := os.WriteFile(cachePath, data, 0600); err != nil {
        return err
    }
    
    // Write metadata
    meta := fmt.Sprintf("sha256:%s\nwords:%d\nfetched:%s\nurl:%s\n",
        actualHash, len(words), time.Now().UTC().Format(time.RFC3339), source.URL)
    if err := os.WriteFile(cachePath+".meta", []byte(meta), 0600); err != nil {
        return err
    }
    
    slog.Info("wordlist cached",
        slog.String("id", source.ID),
        slog.Int("words", len(words)),
    )
    
    return nil
}

func (m *Manager) cachePath(id string) string {
    return filepath.Join(m.cacheDir, id+".txt")
}

func (m *Manager) isValidCache(path string, source WordlistSource) bool {
    info, err := os.Stat(path)
    if err != nil || info.Size() == 0 {
        return false
    }
    
    if source.SHA256 != "" {
        data, err := os.ReadFile(path)
        if err != nil {
            return false
        }
        hash := sha256.Sum256(data)
        if hex.EncodeToString(hash[:]) != source.SHA256 {
            return false
        }
    }
    
    return true
}

func (m *Manager) AvailableCount() int {
    count := 0
    for _, src := range m.sources {
        if src.Enabled && m.isValidCache(m.cachePath(src.ID), src) {
            count++
        }
    }
    return count
}

// LoadAll loads all available wordlists into memory
func (m *Manager) LoadAll() error {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    for _, source := range m.sources {
        if !source.Enabled {
            continue
        }
        
        data, err := os.ReadFile(m.cachePath(source.ID))
        if err != nil {
            continue
        }
        
        words := parseWordlist(data, source.ID)
        m.loaded[source.ID] = NewWordlist(source.ID, words, source.Category)
    }
    
    if len(m.loaded) < 3 {
        return fmt.Errorf("%w: need 3, have %d", ErrInsufficientLists, len(m.loaded))
    }
    
    return nil
}

// SelectRandomLists returns n randomly selected wordlists
func (m *Manager) SelectRandomLists(n int) ([]*Wordlist, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    if len(m.loaded) < n {
        return nil, fmt.Errorf("requested %d lists, only %d available", n, len(m.loaded))
    }
    
    available := make([]*Wordlist, 0, len(m.loaded))
    for _, wl := range m.loaded {
        available = append(available, wl)
    }
    
    // Fisher-Yates shuffle with crypto/rand
    for i := len(available) - 1; i > 0; i-- {
        j, _ := SecureRandInt(i + 1)
        available[i], available[j] = available[j], available[i]
    }
    
    return available[:n], nil
}

// AddUserWordlist adds a user-provided wordlist
func (m *Manager) AddUserWordlist(path, id string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return err
    }
    
    words := parseWordlist(data, id)
    if len(words) < 100 {
        return fmt.Errorf("wordlist too small: %d words (minimum 100)", len(words))
    }
    
    m.mu.Lock()
    m.loaded[id] = NewWordlist(id, words, "user")
    m.mu.Unlock()
    
    slog.Info("loaded user wordlist", slog.String("id", id), slog.Int("words", len(words)))
    return nil
}

// ListSources returns information about all configured sources
func (m *Manager) ListSources() []WordlistSource {
    return slices.Clone(m.sources)
}

// EnableSource enables or disables a wordlist source
func (m *Manager) EnableSource(id string, enabled bool) bool {
    for i := range m.sources {
        if m.sources[i].ID == id {
            m.sources[i].Enabled = enabled
            return true
        }
    }
    return false
}

func parseWordlist(data []byte, sourceID string) []string {
    lines := strings.Split(string(data), "\n")
    words := make([]string, 0, len(lines))
    
    for _, line := range lines {
        line = strings.TrimSpace(line)
        
        // Skip empty/comments
        if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
            continue
        }
        
        // Handle diceware format: "11111\tword" or "11111 word"
        if fields := strings.Fields(line); len(fields) >= 2 {
            if isAllDigits(fields[0]) {
                line = fields[1]
            }
        }
        
        // Handle C header format (.h files)
        if strings.Contains(line, `"`) {
            line = strings.Trim(line, `",; `)
            if idx := strings.Index(line, `"`); idx >= 0 {
                line = line[idx+1:]
                if end := strings.Index(line, `"`); end >= 0 {
                    line = line[:end]
                }
            }
        }
        
        word := strings.ToLower(strings.TrimSpace(line))
        if isValidWord(word) {
            words = append(words, word)
        }
    }
    
    slices.Sort(words)
    return slices.Compact(words)
}

func isAllDigits(s string) bool {
    if len(s) == 0 {
        return false
    }
    for _, r := range s {
        if r < '0' || r > '9' {
            return false
        }
    }
    return true
}

func isValidWord(s string) bool {
    if len(s) < 2 || len(s) > 12 {
        return false
    }
    for _, r := range s {
        if r < 'a' || r > 'z' {
            return false
        }
    }
    return true
}
```

---

## Exclusion List System

### Default Exclusions

Glyphic includes sensible default exclusions to prevent generating embarrassing or offensive passwords:

```go
package wordlist

import (
    "bufio"
    "embed"
    "os"
    "slices"
    "strings"
    "sync"
)

//go:embed exclusions/*.txt
var embeddedExclusions embed.FS

var defaultExclusions = sync.OnceValue(func() []string {
    var words []string
    
    entries, _ := embeddedExclusions.ReadDir("exclusions")
    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }
        data, err := embeddedExclusions.ReadFile("exclusions/" + entry.Name())
        if err != nil {
            continue
        }
        for _, line := range strings.Split(string(data), "\n") {
            line = strings.TrimSpace(strings.ToLower(line))
            if line != "" && !strings.HasPrefix(line, "#") {
                words = append(words, line)
            }
        }
    }
    
    slices.Sort(words)
    return slices.Compact(words)
})

type ExclusionList struct {
    words []string // Sorted for binary search
    mu    sync.RWMutex
}

func NewExclusionList(useDefaults bool) *ExclusionList {
    el := &ExclusionList{}
    if useDefaults {
        el.words = slices.Clone(defaultExclusions())
    }
    return el
}

func (e *ExclusionList) LoadFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()
    
    var newWords []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        word := strings.TrimSpace(strings.ToLower(scanner.Text()))
        if word != "" && !strings.HasPrefix(word, "#") {
            newWords = append(newWords, word)
        }
    }
    
    e.mu.Lock()
    e.words = append(e.words, newWords...)
    slices.Sort(e.words)
    e.words = slices.Compact(e.words)
    e.mu.Unlock()
    
    return scanner.Err()
}

func (e *ExclusionList) Add(words ...string) {
    e.mu.Lock()
    for _, w := range words {
        e.words = append(e.words, strings.ToLower(w))
    }
    slices.Sort(e.words)
    e.words = slices.Compact(e.words)
    e.mu.Unlock()
}

func (e *ExclusionList) Contains(word string) bool {
    e.mu.RLock()
    defer e.mu.RUnlock()
    _, found := slices.BinarySearch(e.words, strings.ToLower(word))
    return found
}

func (e *ExclusionList) Filter(words []string) []string {
    e.mu.RLock()
    defer e.mu.RUnlock()
    
    return slices.DeleteFunc(slices.Clone(words), func(w string) bool {
        _, found := slices.BinarySearch(e.words, strings.ToLower(w))
        return found
    })
}

func (e *ExclusionList) Count() int {
    e.mu.RLock()
    defer e.mu.RUnlock()
    return len(e.words)
}

func (e *ExclusionList) Disable() {
    e.mu.Lock()
    e.words = nil
    e.mu.Unlock()
}
```

### Embedded Exclusion Files

```text
internal/wordlist/exclusions/
├── profanity.txt       # Common profanity
├── slurs.txt           # Ethnic, racial, gender slurs
├── sensitive.txt       # Potentially offensive terms
└── confusing.txt       # Words confused with numbers/letters
```

Example `exclusions/confusing.txt`:

```text
# Words easily confused when spoken or typed
one
won
two
to
too
for
four
fore
ate
eight
oh
owe
sea
see
bee
be
eye
aye
i
you
ewe
why
y
are
our
hour
```

---

## Multi-Wordlist Password Generation

### Strategy

Each password draws words from **at least 3 randomly-selected wordlists**. This:

1. Eliminates single-source predictability
2. Mixes word styles and lengths
3. Makes wordlist inference attacks significantly harder

```go
package generator

import (
    "cmp"
    "fmt"
    "iter"
    "math"
    "slices"
    
    "github.com/youruser/glyphic/internal/security"
    "github.com/youruser/glyphic/internal/wordlist"
)

type Options struct {
    WordCount      int
    Separator      string
    Capitalize     string   // "none", "first", "each", "random", "all"
    Numbers        int
    Specials       int
    SpecialSet     string
    Misspell       bool
    MisspellRate   float64
    MinEntropy     float64
    MinWordlists   int      // Minimum wordlists to draw from (default: 3)
    PreferredLists []string // Preferred wordlist IDs (optional)
}

var DefaultOptions = Options{
    WordCount:    6,
    Separator:    "-",
    Capitalize:   "none",
    Numbers:      0,
    Specials:     0,
    SpecialSet:   "!@#$%^&*-_+=",
    Misspell:     false,
    MisspellRate: 0.3,
    MinEntropy:   0,
    MinWordlists: 3,
}

type Generator struct {
    manager    *wordlist.Manager
    exclusions *wordlist.ExclusionList
    opts       Options
}

func New(manager *wordlist.Manager, exclusions *wordlist.ExclusionList, opts Options) *Generator {
    // Apply defaults
    opts.WordCount = cmp.Or(opts.WordCount, DefaultOptions.WordCount)
    opts.Separator = cmp.Or(opts.Separator, DefaultOptions.Separator)
    opts.SpecialSet = cmp.Or(opts.SpecialSet, DefaultOptions.SpecialSet)
    opts.MinWordlists = cmp.Or(opts.MinWordlists, DefaultOptions.MinWordlists)
    opts.MisspellRate = cmp.Or(opts.MisspellRate, DefaultOptions.MisspellRate)
    
    // Clamp ranges
    opts.WordCount = cmp.Clamp(opts.WordCount, 4, 10)
    opts.MinWordlists = cmp.Clamp(opts.MinWordlists, 3, 10)
    opts.Numbers = cmp.Clamp(opts.Numbers, 0, 10)
    opts.Specials = cmp.Clamp(opts.Specials, 0, 10)
    
    return &Generator{
        manager:    manager,
        exclusions: exclusions,
        opts:       opts,
    }
}

// Generate creates a single passphrase using multiple wordlists
// Caller must call security.SecureZero on returned bytes
func (g *Generator) Generate() ([]byte, []string, error) {
    // Select random wordlists
    lists, err := g.manager.SelectRandomLists(g.opts.MinWordlists)
    if err != nil {
        return nil, nil, err
    }
    
    listIDs := make([]string, len(lists))
    for i, l := range lists {
        listIDs[i] = l.ID()
    }
    
    // Pre-allocate buffer
    buf := make([]byte, 0, g.opts.WordCount*8+g.opts.Numbers+g.opts.Specials)
    
    // Generate words, each from a randomly-selected list
    for i := range g.opts.WordCount {
        if i > 0 && g.opts.Separator != "" {
            buf = append(buf, g.opts.Separator...)
        }
        
        // Pick random list for this word
        listIdx, _ := security.SecureRandInt(len(lists))
        selectedList := lists[listIdx]
        
        // Select word with exclusion filtering
        var word string
        for attempts := 0; attempts < 100; attempts++ {
            word, err = selectedList.SelectRandom()
            if err != nil {
                return nil, nil, err
            }
            if g.exclusions == nil || !g.exclusions.Contains(word) {
                break
            }
            word = ""
        }
        
        if word == "" {
            return nil, nil, fmt.Errorf("could not find non-excluded word after 100 attempts")
        }
        
        wordBytes := []byte(word)
        
        // Apply modifications
        if g.opts.Misspell {
            wordBytes = g.applyMisspelling(wordBytes)
        }
        wordBytes = g.applyCapitalization(wordBytes, i)
        
        buf = append(buf, wordBytes...)
    }
    
    // Insert random numbers
    for range g.opts.Numbers {
        digit, _ := security.SecureRandInt(10)
        pos, _ := security.SecureRandInt(len(buf) + 1)
        buf = slices.Insert(buf, pos, byte('0'+digit))
    }
    
    // Insert special characters
    if g.opts.Specials > 0 {
        specials := []byte(g.opts.SpecialSet)
        for range g.opts.Specials {
            idx, _ := security.SecureRandInt(len(specials))
            pos, _ := security.SecureRandInt(len(buf) + 1)
            buf = slices.Insert(buf, pos, specials[idx])
        }
    }
    
    return buf, listIDs, nil
}

// GenerateN returns an iterator yielding n passwords
func (g *Generator) GenerateN(n int) iter.Seq2[int, []byte] {
    return func(yield func(int, []byte) bool) {
        for i := range n {
            pwd, _, err := g.Generate()
            if err != nil {
                return
            }
            if !yield(i, pwd) {
                security.SecureZero(pwd)
                return
            }
        }
    }
}

func (g *Generator) applyCapitalization(word []byte, index int) []byte {
    if len(word) == 0 {
        return word
    }
    
    switch g.opts.Capitalize {
    case "first":
        if index == 0 && word[0] >= 'a' && word[0] <= 'z' {
            word[0] -= 32
        }
    case "each":
        if word[0] >= 'a' && word[0] <= 'z' {
            word[0] -= 32
        }
    case "random":
        if flip, _ := security.SecureRandInt(2); flip == 1 {
            if word[0] >= 'a' && word[0] <= 'z' {
                word[0] -= 32
            }
        }
    case "all":
        for i := range word {
            if word[i] >= 'a' && word[i] <= 'z' {
                word[i] -= 32
            }
        }
    }
    
    return word
}

func (g *Generator) applyMisspelling(word []byte) []byte {
    rate, _ := security.SecureRandFloat()
    if rate >= g.opts.MisspellRate {
        return word
    }
    return ApplyMisspelling(word, DefaultMisspellOptions)
}

// Stats includes generation statistics and version information
type Stats struct {
    Version         string
    Count           int
    AvgLength       float64
    MinLength       int
    MaxLength       int
    AvgEntropy      float64
    CharDistribution map[rune]int
}

// CalculateEntropy computes entropy considering multi-list selection
func (g *Generator) CalculateEntropy() EntropyBreakdown {
    lists, _ := g.manager.SelectRandomLists(g.opts.MinWordlists)
    
    // Conservative estimate using smallest list
    minSize := math.MaxInt
    totalWords := 0
    for _, wl := range lists {
        size := wl.Size()
        totalWords += size
        if size < minSize {
            minSize = size
        }
    }
    
    base := float64(g.opts.WordCount) * math.Log2(float64(minSize))
    listSelection := float64(g.opts.WordCount) * math.Log2(float64(len(lists)))
    
    var capEnt float64
    if g.opts.Capitalize == "random" {
        capEnt = float64(g.opts.WordCount)
    }
    
    numEnt := float64(g.opts.Numbers) * math.Log2(10)
    specEnt := float64(g.opts.Specials) * math.Log2(float64(len(g.opts.SpecialSet)))
    
    avgLen := g.opts.WordCount * 5
    posEnt := float64(g.opts.Numbers+g.opts.Specials) * math.Log2(float64(avgLen))
    
    total := base + listSelection + capEnt + numEnt + specEnt + posEnt
    
    return EntropyBreakdown{
        WordlistSize:     minSize,
        TotalPoolSize:    totalWords,
        NumLists:         len(lists),
        WordCount:        g.opts.WordCount,
        BaseEntropy:      base,
        ListSelectionEnt: listSelection,
        CapEntropy:       capEnt,
        NumberEntropy:    numEnt,
        SpecialEntropy:   specEnt,
        PositionEntropy:  posEnt,
        TotalEntropy:     total,
        QuantumEffective: total / 2,
        Strength:         strengthRating(total),
    }
}
```

---

## Version Management

The project uses the `pkg/version` package for consistent versioning across builds with Git integration.

### Version Package Integration

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
    
    "github.com/youruser/glyphic/pkg/version"
)

func main() {
    // Version flag handling
    if showVersion {
        printVersion(false)
        os.Exit(0)
    }
    
    if showVersionJSON {
        printVersion(true)
        os.Exit(0)
    }
    
    // Application logic...
}

func printVersion(asJSON bool) {
    info := version.GetBuildInfo()
    
    if asJSON {
        encoder := json.NewEncoder(os.Stdout)
        encoder.SetIndent("", "  ")
        if err := encoder.Encode(info); err != nil {
            fmt.Fprintf(os.Stderr, "Error encoding version: %v\n", err)
            os.Exit(1)
        }
        return
    }
    
    // Human-readable format
    fmt.Printf("glyphic %s\n", info.Version)
    fmt.Printf("Git Commit: %s\n", info.GitCommit)
    fmt.Printf("Build Time: %s\n", info.BuildTime)
    fmt.Printf("Go Version: %s\n", info.GoVersion)
}
```

### Build Integration

The VERSION file at the project root contains the semantic version. Build-time information is injected via ldflags:

```makefile
# Makefile excerpt
VERSION := $(shell cat VERSION)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION := $(shell go version | cut -d' ' -f3)

LDFLAGS := -ldflags "\
    -X github.com/youruser/glyphic/pkg/version.Version=$(VERSION) \
    -X github.com/youruser/glyphic/pkg/version.GitCommit=$(GIT_COMMIT) \
    -X github.com/youruser/glyphic/pkg/version.BuildTime=$(BUILD_TIME) \
    -X github.com/youruser/glyphic/pkg/version.GoVersion=$(GO_VERSION)"

build:
	@echo "Building glyphic v$(VERSION) ($(GIT_COMMIT))"
	go build $(LDFLAGS) -o bin/glyphic ./cmd/glyphic
```

### Version in Logs

Include version information in structured logs for debugging:

```go
import (
    "log/slog"
    "github.com/youruser/glyphic/pkg/version"
)

func initLogging() {
    info := version.GetBuildInfo()
    
    slog.Info("starting glyphic",
        slog.String("version", info.Version),
        slog.String("commit", info.GitCommit),
        slog.String("build_time", info.BuildTime),
        slog.String("go_version", info.GoVersion),
    )
}
```

---

## CLI Interface

```bash
Usage: glyphic [flags]

Generation Options:
  -n, --count int          Number of passwords (1 to 1B) [default: 1]
  -w, --words int          Words per passphrase (4-10) [default: 6]
  -s, --separator string   Word separator [default: "-"]
  --min-lists int          Minimum wordlists to use (3-10) [default: 3]
  --prefer-lists strings   Preferred wordlist IDs (comma-separated)
  --max-length int         Maximum password length (for compatibility)
  --min-length int         Minimum password length
  --prefix string          Add prefix to each password
  --suffix string          Add suffix to each password
  --template string        Output template (e.g., "Password: {password}")

Wordlist Management:
  --wordlist string        Add custom wordlist file (can repeat)
  --list-sources           Show available wordlist sources and status
  --fetch                  Force re-fetch all wordlists
  --cache-dir string       Cache directory [default: ~/.local/share/glyphic/wordlists]
  --verify-checksums       Re-verify all cached wordlist checksums
  --enable-source strings  Enable specific wordlist sources
  --disable-source strings Disable specific wordlist sources
  --language string        Select wordlist language [default: "en"]

Exclusion Options:
  --exclude-file string    Additional exclusion file (can repeat)
  --exclude strings        Words to exclude (comma-separated)
  --no-default-exclusions  Disable built-in exclusion lists
  --show-exclusions        Show exclusion count and categories

Modification Options:
  -c, --capitalize string  "none", "first", "each", "random", "all" [default: "none"]
  --numbers int            Random digits to insert [default: 0]
  --specials int           Special characters to insert [default: 0]
  --special-set string     Character set for specials [default: "!@#$%^&*-_+="]
  --misspell               Apply random misspellings
  --misspell-rate float    Misspelling probability [default: 0.3]

Output Options:
  -o, --output string      Output file [default: stdout]
  -f, --format string      "plain", "json", "csv" [default: "plain"]
  -q, --quiet              No visual effects or status
  --show-entropy           Display entropy per password
  --show-sources           Show source wordlist for each word
  --no-newline             Don't print trailing newline
  --json-schema            Output JSON schema for programmatic use
  --stats                  Show password statistics (avg length, char distribution)

Visual Effects:
  --no-reveal              Disable decode animation
  --reveal-count int       Passwords to animate [default: 5 if TTY]
  --reveal-speed string    "slow", "normal", "fast" [default: "normal"]
  --color-scheme string    "matrix", "cyber", "fire", "vapor", "mono"

Clipboard:
  -C, --copy               Copy password to clipboard and auto-clear
  --clipboard-timeout int  Seconds before clipboard auto-clears [default: 30]
  --no-clipboard-history   Prevent clipboard managers from saving

Security:
  --min-entropy int        Minimum required entropy (bits)
  --verify                 Show detailed entropy breakdown

Configuration:
  --config string          Load settings from config file
  --save-config            Save current flags to config file

Utility:
  -v, --version            Show version information (version, commit, build time)
  --version-json           Output version information as JSON
  -h, --help               Show this help message
  --benchmark              Run generation benchmark

Examples:
  glyphic                                    # Fetch wordlists (first run), generate password
  glyphic --list-sources                     # Show wordlist sources
  glyphic -n 10 -w 8 --min-lists 5           # 10 passwords, 8 words, 5+ wordlists
  glyphic --wordlist ~/custom.txt            # Include custom wordlist
  glyphic --exclude "password,secret"        # Exclude specific words
  glyphic --no-default-exclusions            # Allow all words
  glyphic -n 1000000 -q -o bulk.txt          # Bulk generation
  glyphic -C --clipboard-timeout 45          # Copy to clipboard, clear after 45s
  glyphic --verify-checksums                 # Verify integrity of cached wordlists
  glyphic --config ~/.glyphic.yaml           # Load configuration from file
  glyphic --max-length 32                    # Generate password under 32 characters
  glyphic --stats -n 1000                    # Generate 1000 and show statistics
```

---

## Included Assets

### Font Files and Glyphs

The project includes the **Matrix Code NFI** font (`MatrixCodeNfi-YPPj.otf`) for the decode animation, along with a curated set of additional sci-fi glyphs for enhanced visual effects.

```go
package font

import (
    "embed"
    "iter"
    "slices"
    "sync"
    
    "golang.org/x/image/font/opentype"
    "golang.org/x/image/font/sfnt"
    
    "github.com/youruser/glyphic/internal/security"
)

//go:embed MatrixCodeNfi-YPPj.otf
var fontData []byte

var parsedFont = sync.OnceValue(func() *sfnt.Font {
    f, err := opentype.Parse(fontData)
    if err != nil {
        panic("embedded font corrupted: " + err.Error())
    }
    return f
})

// Additional curated glyphs for enhanced scramble effect
// These are CLI-safe Unicode code points from various writing systems featured in the Matrix Code NFI font
// and other sci-fi-looking scripts. All glyphs are validated for terminal compatibility.
var additionalGlyphs = []rune{
    // Basic geometric shapes and symbols (always safe)
    '▲', '▼', '◆', '◇', '●', '○', '■', '□', '▪', '▫',
    
    // Box drawing characters (universal terminal support)
    '─', '━', '│', '┃', '┌', '┐', '└', '┘', '├', '┤', '┬', '┴', '┼',
    '═', '║', '╔', '╗', '╚', '╝', '╠', '╣', '╦', '╩', '╬',
    
    // Mathematical operators (technical aesthetic)
    '∀', '∃', '∈', '∉', '∋', '∏', '∑', '−', '∓', '∗', '∘', '∙', '√',
    '∝', '∞', '∠', '∧', '∨', '∩', '∪', '∫', '∴', '∵', '∼', '≃', '≅',
    '≈', '≠', '≡', '≤', '≥', '⊂', '⊃', '⊆', '⊇', '⊕', '⊗', '⊥',
    
    // Arrows and pointers
    '←', '↑', '→', '↓', '↔', '↕', '⇐', '⇑', '⇒', '⇓', '⇔', '⇕',
    
    // Technical symbols
    '⌀', '⌂', '⌈', '⌉', '⌊', '⌋', '⎕', '⏏', '␀', '␣',
    
    // Greek letters (subset for technical/sci-fi feel)
    'Α', 'Β', 'Γ', 'Δ', 'Θ', 'Λ', 'Ξ', 'Π', 'Σ', 'Φ', 'Ψ', 'Ω',
    'α', 'β', 'γ', 'δ', 'ε', 'θ', 'λ', 'μ', 'π', 'σ', 'φ', 'ψ', 'ω',
    
    // Block elements
    '░', '▒', '▓', '█', '▄', '▌', '▐', '▀',
    
    // Miscellaneous technical
    '⊙', '⊚', '⊛', '⊜', '⊝', '⊞', '⊟', '⊠', '⊡',
    '⌘', '⌥', '⎇', '⏎', '⏻', '⏼', '⏽',
    
    // CJK brackets and symbols (from Matrix aesthetic)
    '〈', '〉', '《', '》', '「', '」', '『', '』', '【', '】',
    
    // Buginese script (U+1A00-U+1A1F) - CLI-safe subset
    // Decorative consonants and vowels with alien appearance
    '\u1A00', '\u1A01', '\u1A02', '\u1A03', '\u1A04', '\u1A05', '\u1A06', '\u1A07',
    '\u1A08', '\u1A09', '\u1A0A', '\u1A0B', '\u1A0C', '\u1A0D', '\u1A0E', '\u1A0F',
    '\u1A10', '\u1A11', '\u1A12', '\u1A13', '\u1A14', '\u1A15', '\u1A16',
    
    // Kayah Li script (U+A900-U+A92F) - sci-fi appearance
    '\uA900', '\uA901', '\uA902', '\uA903', '\uA904', '\uA905', '\uA906', '\uA907',
    '\uA908', '\uA909', '\uA90A', '\uA90B', '\uA90C', '\uA90D', '\uA90E', '\uA90F',
    '\uA910', '\uA911', '\uA912', '\uA913', '\uA914', '\uA915', '\uA916', '\uA917',
    '\uA918', '\uA919', '\uA91A', '\uA91B', '\uA91C', '\uA91D', '\uA91E', '\uA91F',
    
    // Phags-pa script (U+A840-U+A877) - Mongolian-derived, Matrix-like
    '\uA840', '\uA841', '\uA842', '\uA843', '\uA844', '\uA845', '\uA846', '\uA847',
    '\uA848', '\uA849', '\uA84A', '\uA84B', '\uA84C', '\uA84D', '\uA84E', '\uA84F',
    '\uA850', '\uA851', '\uA852', '\uA853', '\uA854', '\uA855', '\uA856', '\uA857',
    '\uA858', '\uA859', '\uA85A', '\uA85B', '\uA85C', '\uA85D', '\uA85E', '\uA85F',
    '\uA860', '\uA861', '\uA862', '\uA863', '\uA864', '\uA865', '\uA866', '\uA867',
    '\uA868', '\uA869', '\uA86A', '\uA86B', '\uA86C', '\uA86D', '\uA86E', '\uA86F',
    '\uA870', '\uA871', '\uA872', '\uA873', '\uA874', '\uA875', '\uA876', '\uA877',
    
    // Tibetan script (U+0F00-U+0F47) - architectural, technical appearance
    '\u0F00', '\u0F01', '\u0F02', '\u0F03', '\u0F04', '\u0F05', '\u0F06', '\u0F07',
    '\u0F08', '\u0F09', '\u0F0A', '\u0F0B', '\u0F0C', '\u0F0D', '\u0F0E', '\u0F0F',
    '\u0F10', '\u0F11', '\u0F12', '\u0F13', '\u0F14', '\u0F15', '\u0F16', '\u0F17',
    '\u0F18', '\u0F19', '\u0F1A', '\u0F1B', '\u0F1C', '\u0F1D', '\u0F1E', '\u0F1F',
    '\u0F20', '\u0F21', '\u0F22', '\u0F23', '\u0F24', '\u0F25', '\u0F26', '\u0F27',
    '\u0F28', '\u0F29', '\u0F2A', '\u0F2B', '\u0F2C', '\u0F2D', '\u0F2E', '\u0F2F',
    '\u0F30', '\u0F31', '\u0F32', '\u0F33', '\u0F34', '\u0F35', '\u0F36', '\u0F37',
    '\u0F38', '\u0F39', '\u0F3A', '\u0F3B', '\u0F3C', '\u0F3D', '\u0F3E', '\u0F3F',
    '\u0F40', '\u0F41', '\u0F42', '\u0F43', '\u0F44', '\u0F45', '\u0F46', '\u0F47',
    
    // Kannada script (U+0C80-U+0CDD) - select geometric characters
    '\u0C82', '\u0C83', '\u0C85', '\u0C86', '\u0C87', '\u0C88', '\u0C89', '\u0C8A',
    '\u0C8B', '\u0C8C', '\u0C8E', '\u0C8F', '\u0C90', '\u0C92', '\u0C93', '\u0C94',
    '\u0C95', '\u0C96', '\u0C97', '\u0C98', '\u0C99', '\u0C9A', '\u0C9B', '\u0C9C',
    '\u0C9D', '\u0C9E', '\u0C9F', '\u0CA0', '\u0CA1', '\u0CA2', '\u0CA3', '\u0CA4',
    '\u0CA5', '\u0CA6', '\u0CA7', '\u0CA8', '\u0CAA', '\u0CAB', '\u0CAC', '\u0CAD',
    '\u0CAE', '\u0CAF', '\u0CB0', '\u0CB1', '\u0CB2', '\u0CB3', '\u0CB5', '\u0CB6',
    '\u0CB7', '\u0CB8', '\u0CB9', '\u0CBC', '\u0CBD', '\u0CBE', '\u0CBF', '\u0CC0',
    '\u0CC1', '\u0CC2', '\u0CC3', '\u0CC4', '\u0CC6', '\u0CC7', '\u0CC8', '\u0CCA',
    '\u0CCB', '\u0CCC', '\u0CCD', '\u0CD5', '\u0CD6', '\u0CDD',
    
    // Mongolian script (U+1800-U+1877) - vertical-style characters
    '\u1800', '\u1801', '\u1802', '\u1803', '\u1804', '\u1805', '\u1806', '\u1807',
    '\u1808', '\u1809', '\u180A', '\u180B', '\u180C', '\u180D', '\u180E',
    '\u1820', '\u1821', '\u1822', '\u1823', '\u1824', '\u1825', '\u1826', '\u1827',
    '\u1828', '\u1829', '\u182A', '\u182B', '\u182C', '\u182D', '\u182E', '\u182F',
    '\u1830', '\u1831', '\u1832', '\u1833', '\u1834', '\u1835', '\u1836', '\u1837',
    '\u1838', '\u1839', '\u183A', '\u183B', '\u183C', '\u183D', '\u183E', '\u183F',
    '\u1840', '\u1841', '\u1842', '\u1843', '\u1844', '\u1845', '\u1846', '\u1847',
    '\u1848', '\u1849', '\u184A', '\u184B', '\u184C', '\u184D', '\u184E', '\u184F',
    '\u1850', '\u1851', '\u1852', '\u1853', '\u1854', '\u1855', '\u1856', '\u1857',
    '\u1858', '\u1859', '\u185A', '\u185B', '\u185C', '\u185D', '\u185E', '\u185F',
    '\u1860', '\u1861', '\u1862', '\u1863', '\u1864', '\u1865', '\u1866', '\u1867',
    '\u1868', '\u1869', '\u186A', '\u186B', '\u186C', '\u186D', '\u186E', '\u186F',
    '\u1870', '\u1871', '\u1872', '\u1873', '\u1874', '\u1875', '\u1876', '\u1877',
    
    // Saurashtra script (U+A880-U+A8DF) - modified/geometric subset
    '\uA880', '\uA881', '\uA882', '\uA883', '\uA884', '\uA885', '\uA886', '\uA887',
    '\uA888', '\uA889', '\uA88A', '\uA88B', '\uA88C', '\uA88D', '\uA88E', '\uA88F',
    '\uA890', '\uA891', '\uA892', '\uA893', '\uA894', '\uA895', '\uA896', '\uA897',
    '\uA898', '\uA899', '\uA89A', '\uA89B', '\uA89C', '\uA89D', '\uA89E', '\uA89F',
    '\uA8A0', '\uA8A1', '\uA8A2', '\uA8A3', '\uA8A4', '\uA8A5', '\uA8A6', '\uA8A7',
    '\uA8A8', '\uA8A9', '\uA8AA', '\uA8AB', '\uA8AC', '\uA8AD', '\uA8AE', '\uA8AF',
    '\uA8B0', '\uA8B1', '\uA8B2', '\uA8B3', '\uA8B4', '\uA8B5', '\uA8B6', '\uA8B7',
    '\uA8B8', '\uA8B9', '\uA8BA', '\uA8BB', '\uA8BC', '\uA8BD', '\uA8BE', '\uA8BF',
    '\uA8C0', '\uA8C1', '\uA8C2', '\uA8C3', '\uA8C4', '\uA8C5', '\uA8CE', '\uA8CF',
    '\uA8D0', '\uA8D1', '\uA8D2', '\uA8D3', '\uA8D4', '\uA8D5', '\uA8D6', '\uA8D7',
    '\uA8D8', '\uA8D9', '\uA8DA', '\uA8DB', '\uA8DC', '\uA8DD', '\uA8DE', '\uA8DF',
    
    // Kanji (CJK Unified Ideographs) - select technical/Matrix-style characters
    // Geometric and technical-looking kanji
    '一', '二', '三', '四', '五', '六', '七', '八', '九', '十',
    '口', '回', '囗', '囚', '四', '因', '団', '困', '囲', '図',
    '工', '王', '主', '玉', '丁', '干', '平', '年', '幸', '半',
    '册', '巾', '市', '布', '帆', '希', '帝', '師', '席', '帯',
    '命', '和', '品', '員', '問', '善', '器', '嘉', '墨', '増',
    '士', '壮', '声', '売', '変', '夢', '奇', '契', '奮', '威',
    '字', '宇', '守', '安', '宗', '官', '定', '宝', '実', '宣',
    '室', '宮', '害', '宴', '宵', '密', '富', '寒', '寛', '寝',
}

// ScrambleGlyphs returns an iterator over all available glyphs
// Combines font glyphs with additional curated glyphs
func ScrambleGlyphs() iter.Seq[rune] {
    return func(yield func(rune) bool) {
        font := parsedFont()
        var buf sfnt.Buffer
        
        seen := make(map[rune]bool)
        
        // First, yield all additional curated glyphs
        for _, r := range additionalGlyphs {
            if !seen[r] {
                seen[r] = true
                if !yield(r) {
                    return
                }
            }
        }
        
        // Then scan font for additional printable glyphs
        for r := rune(0x20); r < 0x10000; r++ {
            if seen[r] {
                continue
            }
            idx, err := font.GlyphIndex(&buf, r)
            if err == nil && idx != 0 && isPrintable(r) {
                seen[r] = true
                if !yield(r) {
                    return
                }
            }
        }
    }
}

var glyphSlice = sync.OnceValue(func() []rune {
    return slices.Collect(ScrambleGlyphs())
})

// RandomGlyph returns a cryptographically random glyph from the available set
func RandomGlyph() rune {
    glyphs := glyphSlice()
    if len(glyphs) == 0 {
        return '?' // Fallback
    }
    idx, _ := security.SecureRandInt(len(glyphs))
    return glyphs[idx]
}

// RandomGlyphN returns n random glyphs
func RandomGlyphN(n int) []rune {
    if n <= 0 {
        return nil
    }
    result := make([]rune, n)
    for i := range n {
        result[i] = RandomGlyph()
    }
    return result
}

// GlyphCategories provides categorized access to glyphs for themed effects
type GlyphSet struct {
    Geometric  []rune
    Technical  []rune
    Math       []rune
    Greek      []rune
    Blocks     []rune
    Arrows     []rune
}

var glyphCategories = sync.OnceValue(func() *GlyphSet {
    return &GlyphSet{
        Geometric: []rune{'▲', '▼', '◆', '◇', '●', '○', '■', '□', '▪', '▫'},
        Technical: []rune{'⌀', '⌂', '⌈', '⌉', '⌊', '⌋', '⎕', '⏏', '␀', '␣', '⌘', '⌥', '⎇'},
        Math:      []rune{'∀', '∃', '∈', '∋', '∏', '∑', '√', '∝', '∞', '∠', '∧', '∨', '∩', '∪', '∫'},
        Greek:     []rune{'Α', 'Β', 'Γ', 'Δ', 'Θ', 'Λ', 'Ξ', 'Π', 'Σ', 'Φ', 'Ψ', 'Ω'},
        Blocks:    []rune{'░', '▒', '▓', '█', '▄', '▌', '▐', '▀'},
        Arrows:    []rune{'←', '↑', '→', '↓', '↔', '↕', '⇐', '⇑', '⇒', '⇓', '⇔', '⇕'},
    }
})

// RandomGlyphFromCategory returns a random glyph from a specific category
func RandomGlyphFromCategory(category string) rune {
    cats := glyphCategories()
    var set []rune
    
    switch category {
    case "geometric":
        set = cats.Geometric
    case "technical":
        set = cats.Technical
    case "math":
        set = cats.Math
    case "greek":
        set = cats.Greek
    case "blocks":
        set = cats.Blocks
    case "arrows":
        set = cats.Arrows
    default:
        return RandomGlyph()
    }
    
    if len(set) == 0 {
        return RandomGlyph()
    }
    
    idx, _ := security.SecureRandInt(len(set))
    return set[idx]
}

func isPrintable(r rune) bool {
    // Filter to single-width, non-combining characters
    // Exclude control characters and common problematic ranges
    if r < 0x20 || r == 0x7F {
        return false
    }
    // Exclude surrogate pairs
    if r >= 0xD800 && r <= 0xDFFF {
        return false
    }
    // Exclude private use area unless specifically in our curated list
    if r >= 0xE000 && r <= 0xF8FF {
        return false
    }
    return true
}

// isCLISafe validates that a glyph will render properly in terminal environments
// This ensures consistent display across different terminal emulators
func isCLISafe(r rune) bool {
    // Control characters
    if r < 0x20 || r == 0x7F || (r >= 0x80 && r <= 0x9F) {
        return false
    }
    
    // Combining characters that might cause rendering issues
    if (r >= 0x0300 && r <= 0x036F) || // Combining Diacritical Marks
       (r >= 0x1AB0 && r <= 0x1AFF) || // Combining Diacritical Marks Extended
       (r >= 0x1DC0 && r <= 0x1DFF) || // Combining Diacritical Marks Supplement
       (r >= 0x20D0 && r <= 0x20FF) || // Combining Diacritical Marks for Symbols
       (r >= 0xFE20 && r <= 0xFE2F) {  // Combining Half Marks
        return false
    }
    
    // Zero-width and invisible characters
    if r == 0x200B || r == 0x200C || r == 0x200D || r == 0xFEFF {
        return false
    }
    
    // Variation selectors
    if r >= 0xFE00 && r <= 0xFE0F {
        return false
    }
    
    // Surrogates
    if r >= 0xD800 && r <= 0xDFFF {
        return false
    }
    
    // Private use area (unless explicitly curated)
    if r >= 0xE000 && r <= 0xF8FF {
        return false
    }
    
    // Non-characters
    if (r >= 0xFDD0 && r <= 0xFDEF) ||
       (r&0xFFFE) == 0xFFFE {
        return false
    }
    
    // Characters known to cause problems in most terminals
    // (bidirectional overrides, etc.)
    if (r >= 0x202A && r <= 0x202E) || // Bidirectional formatting
       (r >= 0x2066 && r <= 0x2069) {  // Bidirectional isolates
        return false
    }
    
    return true
}

// ValidateGlyphsForTerminal filters glyphs that are safe for CLI display
func ValidateGlyphsForTerminal(glyphs []rune) []rune {
    safe := make([]rune, 0, len(glyphs))
    for _, r := range glyphs {
        if isCLISafe(r) && isPrintable(r) {
            safe = append(safe, r)
        }
    }
    return safe
}

// TestGlyphRendering checks if a glyph renders properly in the current terminal
// Returns true if the glyph is likely to display correctly
func TestGlyphRendering(r rune) bool {
    // Basic safety checks
    if !isCLISafe(r) {
        return false
    }
    
    // Check for common terminal capabilities
    term := os.Getenv("TERM")
    
    // Very limited terminals - only ASCII
    if term == "dumb" || term == "cons25" {
        return r <= 0x7E
    }
    
    // Basic UTF-8 support but limited glyphs
    if strings.HasPrefix(term, "vt") || term == "linux" {
        // Stick to Latin-1 supplement and basic multilingual plane basics
        return r <= 0x00FF || 
               (r >= 0x2500 && r <= 0x257F) || // Box drawing
               (r >= 0x2580 && r <= 0x259F)    // Block elements
    }
    
    // Modern terminals with good Unicode support
    if strings.Contains(term, "xterm") || 
       strings.Contains(term, "screen") ||
       strings.Contains(term, "tmux") ||
       term == "alacritty" ||
       term == "kitty" ||
       strings.Contains(term, "256color") {
        return true // These generally support all our glyphs
    }
    
    // Default: be conservative
    return r <= 0x00FF
}

// GlyphCount returns the total number of available glyphs
func GlyphCount() int {
    return len(glyphSlice())
}

// SafeGlyphCount returns the number of CLI-safe glyphs
func SafeGlyphCount() int {
    count := 0
    for _, r := range glyphSlice() {
        if isCLISafe(r) {
            count++
        }
    }
    return count
}
```

---

## Decode Reveal Animation

The signature feature of glyphic is the "decode" text reveal animation that uses the Matrix Code NFI font and additional sci-fi glyphs to create a cinematic scramble-to-reveal effect.

### Animation Strategy

```go
package tui

import (
    "fmt"
    "time"
    
    "github.com/youruser/glyphic/internal/font"
    "github.com/youruser/glyphic/internal/security"
)

type RevealSpeed int

const (
    SpeedSlow RevealSpeed = iota
    SpeedNormal
    SpeedFast
)

func (rs RevealSpeed) Duration() time.Duration {
    switch rs {
    case SpeedSlow:
        return 150 * time.Millisecond
    case SpeedFast:
        return 30 * time.Millisecond
    default: // Normal
        return 60 * time.Millisecond
    }
}

type RevealOptions struct {
    Speed       RevealSpeed
    Iterations  int    // Number of scramble iterations per character
    ColorScheme string // "matrix", "cyber", "fire", "vapor", "mono"
    UseCategories bool // Use categorized glyphs for themed effect
}

var DefaultRevealOptions = RevealOptions{
    Speed:      SpeedNormal,
    Iterations: 5,
    ColorScheme: "matrix",
    UseCategories: false,
}

// RevealPassword animates the password reveal with scrambled glyphs
func RevealPassword(password []byte, opts RevealOptions) {
    passwordRunes := []rune(string(password))
    display := make([]rune, len(passwordRunes))
    
    // Initialize with random glyphs
    for i := range display {
        display[i] = font.RandomGlyph()
    }
    
    duration := opts.Speed.Duration()
    
    // Reveal each character with scramble iterations
    for charIdx := range passwordRunes {
        for iteration := 0; iteration < opts.Iterations; iteration++ {
            // Update all unrevealed positions with new random glyphs
            for i := charIdx; i < len(display); i++ {
                if opts.UseCategories {
                    // Cycle through categories for varied effect
                    categories := []string{"geometric", "technical", "math", "greek"}
                    catIdx := (iteration + i) % len(categories)
                    display[i] = font.RandomGlyphFromCategory(categories[catIdx])
                } else {
                    display[i] = font.RandomGlyph()
                }
            }
            
            // Print current state
            fmt.Printf("\r%s%s%s", 
                colorize(string(passwordRunes[:charIdx]), opts.ColorScheme, true),
                colorize(string(display[charIdx:]), opts.ColorScheme, false),
                "\033[K") // Clear to end of line
            
            time.Sleep(duration)
        }
        
        // Reveal the actual character
        display[charIdx] = passwordRunes[charIdx]
        fmt.Printf("\r%s%s\033[K",
            colorize(string(passwordRunes[:charIdx+1]), opts.ColorScheme, true),
            colorize(string(display[charIdx+1:]), opts.ColorScheme, false))
        time.Sleep(duration / 2)
    }
    
    // Final reveal with slight pause
    time.Sleep(200 * time.Millisecond)
    fmt.Printf("\r%s\n", colorize(string(passwordRunes), opts.ColorScheme, true))
}

// colorize applies ANSI color codes based on scheme
func colorize(text string, scheme string, revealed bool) string {
    if scheme == "mono" {
        if revealed {
            return "\033[1m" + text + "\033[0m" // Bold
        }
        return "\033[2m" + text + "\033[0m" // Dim
    }
    
    var color string
    switch scheme {
    case "matrix":
        if revealed {
            color = "\033[1;32m" // Bright green
        } else {
            color = "\033[32m" // Green
        }
    case "cyber":
        if revealed {
            color = "\033[1;36m" // Bright cyan
        } else {
            color = "\033[36m" // Cyan
        }
    case "fire":
        if revealed {
            color = "\033[1;33m" // Bright yellow
        } else {
            color = "\033[31m" // Red
        }
    case "vapor":
        if revealed {
            color = "\033[1;35m" // Bright magenta
        } else {
            color = "\033[35m" // Magenta
        }
    default:
        if revealed {
            color = "\033[1;32m"
        } else {
            color = "\033[32m"
        }
    }
    
    return color + text + "\033[0m"
}

// RevealMultiple reveals multiple passwords with staggered start times
func RevealMultiple(passwords [][]byte, opts RevealOptions, maxConcurrent int) {
    maxConcurrent = cmp.Clamp(maxConcurrent, 1, 10)
    
    for i, pwd := range passwords {
        if i >= maxConcurrent {
            break // Don't animate too many
        }
        
        fmt.Printf("\n")
        RevealPassword(pwd, opts)
        
        // Brief pause between passwords
        if i < len(passwords)-1 && i < maxConcurrent-1 {
            time.Sleep(300 * time.Millisecond)
        }
    }
}
```

### Visual Effects Examples

The combination of the Matrix Code NFI font glyphs and the curated Unicode symbols creates several visual effects:

**Matrix Style (Default):**

```text
ꓕᚦ⟐ᛁ◆ᚹ⎊ꓤ-⟒◢ᚱ⍙▼᛭-ꓷ⎋ᚨ◇⟓ᚾ-ᛁ◆ᚹ⎊-ꓤ⟒◢ᚱ-⍙▼᛭ꓷ
∀∃∈∋∏∑-√∝∞∠∧∨-∩∪∫∴∵-∼≃≅-≈≠≡-≤≥⊂
velvet-alpine-crystal-harbor-nebula-theorem
```

**Geometric/Technical:**

```text
▲▼◆◇●○■□▪▫-⌀⌂⌈⌉⌊⌋⎕⏏-░▒▓█▄▌▐▀-◆◇●○■
velvet-alpine-crystal-harbor-nebula-theorem
```

**Greek/Mathematical:**

```text
ΑΒΓΔΘΛΞΠΣΦΨΩαβγδεθλμπσφψω-∀∃∈∋∏∑√∞∧∨∩∪
velvet-alpine-crystal-harbor-nebula-theorem
```

### Glyph Script Selection

The Matrix-inspired reveal effect uses glyphs from diverse writing systems, creating an authentic "falling code" aesthetic.
These scripts were chosen for their visual properties and CLI safety:

**Primary Scripts:**

1. **Buginese** (U+1A00-U+1A1F) - Decorative consonants with alien appearance
2. **Kayah Li** (U+A900-U+A92F) - Angular, technical-looking characters
3. **Phags-pa** (U+A840-U+A877) - Mongolian-derived script with Matrix-like vertical strokes
4. **Tibetan** (U+0F00-U+0F47) - Architectural, technical appearance with stacked components
5. **Kannada** (U+0C80-U+0CDD) - Geometric South Indian script with circular elements
6. **Mongolian** (U+1800-U+1877) - Vertical script with flowing connectors
7. **Saurashtra** (U+A880-U+A8DF) - Modified geometric characters
8. **CJK Kanji** - Technical and geometric ideographs (工, 王, 回, 図, etc.)

**Script Characteristics:**

- **Angular & Geometric**: Phags-pa, Kayah Li provide sharp, technical aesthetics
- **Circular & Flowing**: Kannada, Tibetan add organic complexity
- **Vertical Elements**: Mongolian, some Tibetan characters create unique patterns
- **Block-like**: Kanji ideographs provide solid, recognizable shapes
- **Decorative**: Buginese, Saurashtra offer ornamental complexity

### CLI Safety and Terminal Compatibility

All glyphs are validated for terminal display safety:

```go
// CLI Safety Validation
func init() {
    // Filter glyphs at startup
    additionalGlyphs = ValidateGlyphsForTerminal(additionalGlyphs)
    
    // Log glyph statistics
    slog.Info("glyph system initialized",
        slog.Int("total_glyphs", len(additionalGlyphs)),
        slog.Int("cli_safe", SafeGlyphCount()),
        slog.String("term", os.Getenv("TERM")),
    )
}

// Safety checks performed:
// - No control characters (0x00-0x1F, 0x7F-0x9F)
// - No combining diacriticals that might overlay incorrectly
// - No zero-width or invisible characters
// - No bidirectional formatting overrides
// - No surrogate pairs or non-characters
// - No private use area unless explicitly curated
```

**Terminal Support Levels:**

1. **Full Unicode** (xterm-256color, kitty, alacritty): All scripts render correctly
2. **Good UTF-8** (most modern terminals): Scripts display with minor font variations
3. **Basic UTF-8** (older xterm, screen): Box drawing + Latin works; exotic scripts may show as boxes
4. **ASCII-only** (dumb, vt100): Fallback to basic geometric symbols only

**Automatic Fallback Strategy:**

```go
// Detect terminal capabilities and adjust glyph set
func AdjustGlyphsForTerminal() {
    term := os.Getenv("TERM")
    
    switch {
    case term == "dumb" || !utf8.ValidString(string(additionalGlyphs[0])):
        // ASCII-only fallback
        additionalGlyphs = asciiSafeGlyphs
        slog.Warn("limited terminal detected, using ASCII-safe glyphs")
        
    case term == "linux" || strings.HasPrefix(term, "vt"):
        // Basic UTF-8 - stick to box drawing and basic shapes
        additionalGlyphs = basicUTF8Glyphs
        slog.Info("basic terminal detected, using simplified glyph set")
        
    default:
        // Full Unicode support - use all glyphs
        slog.Info("modern terminal detected, using full glyph set")
    }
}

var asciiSafeGlyphs = []rune{
    '+', '-', '*', '/', '\\', '|', '_', '=', '<', '>',
    '[', ']', '{', '}', '(', ')', '#', '@', '%', '&',
}

var basicUTF8Glyphs = []rune{
    // Box drawing + block elements only
    '─', '│', '┌', '┐', '└', '┘', '├', '┤', '┬', '┴', '┼',
    '░', '▒', '▓', '█', '▀', '▄', '▌', '▐',
    '▲', '▼', '◆', '●', '■',
}
```

### Performance Considerations

- Glyphs are pre-loaded at startup using `sync.OnceValue`
- CLI safety validation runs once during initialization
- Terminal capability detection is cached
- Random glyph selection uses `crypto/rand` for unpredictability
- Categorized glyphs enable thematic reveal patterns
- Animation can be disabled with `--no-reveal` for bulk generation
- Only the first N passwords are animated (configurable with `--reveal-count`)
- Fallback to ASCII-safe glyphs if terminal doesn't support UTF-8

**Font Rendering Notes:**

The Matrix Code NFI font (`MatrixCodeNfi-YPPj.otf`) includes custom glyphs for many of these scripts, but the Unicode fallbacks
ensure the animation works even without the font installed in the terminal. The diverse script selection creates visual variety
while maintaining terminal compatibility.

---

## Security Implementation

```go
package security

import (
    "crypto/rand"
    "crypto/subtle"
    "encoding/binary"
    "errors"
    "log/slog"
    "math"
    "os"
    "runtime"
    "strconv"
    "strings"
    "sync"
    "unsafe"
    
    "golang.org/x/sys/unix"
)

// Verify crypto/rand is properly seeded at startup
func init() {
    var testBuf [32]byte
    if n, err := rand.Read(testBuf[:]); err != nil || n != 32 {
        panic("crypto/rand is not properly initialized")
    }
    
    // Verify non-zero output (extremely unlikely to be all zeros)
    allZero := true
    for _, b := range testBuf {
        if b != 0 {
            allZero = false
            break
        }
    }
    if allZero {
        panic("crypto/rand returned all zeros - possible PRNG failure")
    }
}

// SecureBuffer with canary tokens for memory corruption detection
type SecureBuffer struct {
    canaryHead uint64
    data       []byte
    canaryTail uint64
}

func NewSecureBuffer(size int) *SecureBuffer {
    canary, _ := SecureRandUint64()
    sb := &SecureBuffer{
        canaryHead: canary,
        data:       make([]byte, size),
        canaryTail: canary,
    }
    
    // Try to lock memory (best effort)
    LockMemory(sb.data)
    
    runtime.SetFinalizer(sb, (*SecureBuffer).checkAndDestroy)
    return sb
}

func (sb *SecureBuffer) Bytes() []byte { return sb.data }

func (sb *SecureBuffer) checkAndDestroy() {
    if sb.canaryHead != sb.canaryTail {
        panic("memory corruption detected: canary mismatch")
    }
    UnlockMemory(sb.data)
    SecureZero(sb.data)
}

type SecureBytes struct {
    data []byte
    once sync.Once
}

func NewSecureBytes(size int) *SecureBytes {
    sb := &SecureBytes{data: make([]byte, size)}
    LockMemory(sb.data)
    runtime.SetFinalizer(sb, (*SecureBytes).Destroy)
    return sb
}

func (sb *SecureBytes) Bytes() []byte { return sb.data }

func (sb *SecureBytes) Destroy() {
    sb.once.Do(func() {
        UnlockMemory(sb.data)
        SecureZero(sb.data)
        runtime.SetFinalizer(sb, nil)
    })
}

// LockMemory prevents password memory from being swapped to disk
func LockMemory(b []byte) error {
    if len(b) == 0 {
        return nil
    }
    // Try to lock memory (may fail without privileges)
    if err := unix.Mlock(b); err != nil {
        // Log warning but don't fail - not all systems support this
        slog.Warn("failed to lock memory", slog.String("error", err.Error()))
    }
    return nil
}

func UnlockMemory(b []byte) {
    if len(b) > 0 {
        unix.Munlock(b) // Best effort
    }
}

func SecureZero(b []byte) {
    if len(b) == 0 {
        return
    }
    ptr := unsafe.Pointer(unsafe.SliceData(b))
    for i := range b {
        *(*byte)(unsafe.Add(ptr, i)) = 0
    }
    runtime.KeepAlive(b)
    clear(b)
}

func SecureRandInt(max int) (int, error) {
    if max <= 0 {
        return 0, errors.New("max must be positive")
    }
    maxUint := uint64(max)
    threshold := (math.MaxUint64 / maxUint) * maxUint
    
    var buf [8]byte
    // Limit attempts to prevent infinite loops
    for attempts := 0; attempts < 100; attempts++ {
        if _, err := rand.Read(buf[:]); err != nil {
            return 0, err
        }
        n := binary.BigEndian.Uint64(buf[:])
        if n < threshold {
            return int(n % maxUint), nil
        }
    }
    // Fail if we can't get unbiased random in 100 tries
    return 0, errors.New("failed to generate unbiased random number")
}

func SecureRandUint64() (uint64, error) {
    var buf [8]byte
    if _, err := rand.Read(buf[:]); err != nil {
        return 0, err
    }
    return binary.BigEndian.Uint64(buf[:]), nil
}

func SecureRandFloat() (float64, error) {
    var buf [8]byte
    if _, err := rand.Read(buf[:]); err != nil {
        return 0, err
    }
    return float64(binary.BigEndian.Uint64(buf[:])>>11) / (1 << 53), nil
}

// ConstantTimeEqual provides constant-time string comparison
func ConstantTimeEqual(a, b string) bool {
    return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// CheckEntropyAvailable warns if system entropy is low (Linux only)
func CheckEntropyAvailable() (int, error) {
    data, err := os.ReadFile("/proc/sys/kernel/random/entropy_avail")
    if err != nil {
        return 0, err // Not Linux or unsupported
    }
    
    avail, err := strconv.Atoi(strings.TrimSpace(string(data)))
    if err != nil {
        return 0, err
    }
    
    if avail < 256 {
        slog.Warn("low system entropy", slog.Int("available", avail))
    }
    
    return avail, nil
}
```

---

## Clipboard Integration

```go
package clipboard

import (
    "context"
    "log/slog"
    "time"
    
    "golang.design/x/clipboard"
    "github.com/youruser/glyphic/internal/security"
)

// CopySecure copies password to clipboard with auto-clear
func CopySecure(password []byte, timeout time.Duration, preventHistory bool) error {
    if err := clipboard.Init(); err != nil {
        return err
    }
    
    // On some platforms, we can signal to clipboard managers not to save
    if preventHistory {
        // Implementation depends on platform
        // X11: Set clipboard targets to exclude history
        // macOS: Use transient pasteboard
        // Windows: Use specific clipboard formats
    }
    
    // Copy to clipboard
    clipboard.Write(clipboard.FmtText, password)
    slog.Info("password copied to clipboard", slog.Duration("timeout", timeout))
    
    // Schedule auto-clear
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    
    go func() {
        <-ctx.Done()
        clipboard.Write(clipboard.FmtText, []byte("")) // Clear
        slog.Info("clipboard cleared")
    }()
    
    return nil
}

// SecureCopy is a convenience wrapper
func SecureCopy(password []byte, opts CopyOptions) error {
    timeout := opts.Timeout
    if timeout == 0 {
        timeout = 30 * time.Second
    }
    return CopySecure(password, timeout, opts.PreventHistory)
}

type CopyOptions struct {
    Timeout        time.Duration
    PreventHistory bool
}
```

---

## Wordlist Integrity Monitoring

```go
package wordlist

import (
    "context"
    "crypto/sha256"
    "crypto/tls"
    "encoding/hex"
    "errors"
    "log/slog"
    "slices"
    "sync"
    "time"
)

// MonitorIntegrity detects if cached wordlists have been tampered with
func (m *Manager) MonitorIntegrity(ctx context.Context, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            for _, source := range m.sources {
                if !source.Enabled {
                    continue
                }
                if !m.isValidCache(m.cachePath(source.ID), source) {
                    slog.Error("wordlist integrity check failed",
                        slog.String("id", source.ID),
                    )
                    // Optionally: delete corrupted cache and re-fetch
                    m.handleCorruptedCache(ctx, source)
                }
            }
        }
    }
}

func (m *Manager) handleCorruptedCache(ctx context.Context, source WordlistSource) {
    slog.Warn("attempting to re-fetch corrupted wordlist", slog.String("id", source.ID))
    cachePath := m.cachePath(source.ID)
    os.Remove(cachePath)
    os.Remove(cachePath + ".meta")
    if err := m.fetchAndCache(ctx, source); err != nil {
        slog.Error("failed to recover corrupted wordlist",
            slog.String("id", source.ID),
            slog.String("error", err.Error()),
        )
    }
}

// Rate limiter to prevent abuse
type rateLimiter struct {
    requests map[string]time.Time
    mu       sync.Mutex
}

func newRateLimiter() *rateLimiter {
    return &rateLimiter{
        requests: make(map[string]time.Time),
    }
}

func (rl *rateLimiter) Allow(sourceID string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    if last, ok := rl.requests[sourceID]; ok {
        if time.Since(last) < 5*time.Second {
            return false // Too soon
        }
    }
    rl.requests[sourceID] = time.Now()
    return true
}

// TLS certificate pinning for critical sources
var trustedCertFingerprints = map[string][]string{
    "www.eff.org": {
        // EFF's cert fingerprint - update with actual value
        "SHA256:...",
    },
    "raw.githubusercontent.com": {
        // GitHub's cert fingerprint - update with actual value
        "SHA256:...",
    },
}

func (m *Manager) verifyTLSConnection(conn *tls.Conn, host string) error {
    if pins, ok := trustedCertFingerprints[host]; ok {
        certs := conn.ConnectionState().PeerCertificates
        if len(certs) == 0 {
            return errors.New("no peer certificates")
        }
        
        hash := sha256.Sum256(certs[0].Raw)
        fingerprint := "SHA256:" + hex.EncodeToString(hash[:])
        
        if !slices.Contains(pins, fingerprint) {
            return fmt.Errorf("certificate pinning failed for %s", host)
        }
    }
    return nil
}

// Cross-verify checksums from multiple sources
type ChecksumSource struct {
    URL    string
    Format string // "sha256sum", "json", etc.
}

var checksumVerification = map[string][]ChecksumSource{
    "eff-large": {
        {URL: "https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt.sha256"},
        // Add archive.org mirror or other trusted sources
    },
    // Add for other critical wordlists
}

func (m *Manager) crossVerifyChecksum(source WordlistSource, data []byte) error {
    hash := sha256.Sum256(data)
    actualHash := hex.EncodeToString(hash[:])
    
    // If we have multiple checksum sources, verify against all
    if sources, ok := checksumVerification[source.ID]; ok {
        verified := false
        for _, csSource := range sources {
            // Fetch checksum from alternative source
            // Compare with actual hash
            // If match found, set verified = true
        }
        if !verified {
            return errors.New("checksum verification failed from all sources")
        }
    }
    
    return nil
}
```

---

## Enhanced Exclusion List with Constant-Time Comparison

```go
// Add to ExclusionList methods:

// ContainsConstantTime provides timing-attack resistant word checking
func (e *ExclusionList) ContainsConstantTime(word string) bool {
    e.mu.RLock()
    defer e.mu.RUnlock()
    
    found := false
    normalized := strings.ToLower(word)
    normalizedBytes := []byte(normalized)
    
    for _, excluded := range e.words {
        if subtle.ConstantTimeCompare(normalizedBytes, []byte(excluded)) == 1 {
            found = true
            // Don't break - continue to maintain constant time
        }
    }
    return found
}
```

---

## Project Structure

```text
glyphic/
├── cmd/glyphic/
│   └── main.go
├── internal/
│   ├── generator/
│   │   ├── generator.go
│   │   ├── entropy.go
│   │   ├── misspell.go
│   │   └── *_test.go
│   ├── wordlist/
│   │   ├── manager.go
│   │   ├── sources.go
│   │   ├── wordlist.go
│   │   ├── exclusions.go
│   │   ├── parser.go
│   │   ├── integrity.go      # Integrity monitoring
│   │   ├── ratelimit.go     # Rate limiting for fetches
│   │   ├── pinning.go       # TLS certificate pinning
│   │   ├── exclusions/
│   │   │   ├── profanity.txt
│   │   │   ├── slurs.txt
│   │   │   ├── sensitive.txt
│   │   │   └── confusing.txt
│   │   └── *_test.go
│   ├── security/
│   │   ├── memory.go        # Memory locking and secure buffers
│   │   ├── random.go        # Crypto random with safety checks
│   │   ├── entropy.go       # System entropy monitoring
│   │   ├── constant.go      # Constant-time operations
│   │   └── *_test.go
│   ├── clipboard/
│   │   ├── clipboard.go     # Secure clipboard operations
│   │   ├── platform_*.go    # Platform-specific implementations
│   │   └── *_test.go
│   ├── config/
│   │   ├── config.go        # Configuration management
│   │   ├── loader.go        # Config file loading
│   │   └── *_test.go
│   ├── tui/
│   │   ├── reveal.go
│   │   ├── styles.go
│   │   └── app.go
│   └── font/
│       ├── parser.go
│       ├── glyphs.go
│       ├── categories.go       # Glyph categorization
│       └── MatrixCodeNfi-YPPj.otf
├── fonts/                       # Font source files
│   ├── MatrixCodeNfi-YPPj.otf
│   └── 19u8fazo22581.webp      # Glyph reference image
├── go.mod
├── go.sum
├── Makefile
├── .goreleaser.yml             # Release automation
└── README.md
```

---

## First Run Experience

```bash
$ glyphic

glyphic v1.0.0 - Quantum-Resistant Password Generator

First run: fetching wordlists from verified sources...

  [✓] EFF Large Wordlist         7,776 words
  [✓] EFF Short Wordlist 1       1,296 words
  [✓] EFF Short Wordlist 2       1,296 words
  [✓] SecureDrop Wordlist        6,800 words
  [✓] BIP-39 English             2,048 words
  [✓] Orchard Street Medium      8,192 words
  [✓] PGP Even Words               256 words
  [✓] PGP Odd Words                256 words

Loaded 8 wordlists (27,920 unique words)
Default exclusions: 847 words

Generating passphrase (6 words from 3+ lists)...

  ꓕᚦ⟐ᛁ◆ᚹ⎊ꓤ-⟒◢ᚱ⍙▼᛭-ꓷ⎋ᚨ◇⟓ᚾ-ᛁ◆ᚹ⎊-ꓤ⟒◢ᚱ-⍙▼᛭ꓷ
  
  [decode animation]
  
  velvet-alpine-crystal-harbor-nebula-theorem

  Entropy: 79.2 bits (quantum-effective: 39.6 bits) — strong
  Sources: eff-large, securedrop, orchard-street-medium
```

---

## Build Configuration

### Makefile

```makefile
.PHONY: all build test clean security-check lint

all: security-check build

# Security checks during build
security-check:
	@echo "Running security checks..."
	@which gosec > /dev/null || go install github.com/securego/gosec/v2/cmd/gosec@latest
	@gosec -quiet ./...
	@which staticcheck > /dev/null || go install honnef.co/go/tools/cmd/staticcheck@latest
	@staticcheck ./...
	@go vet ./...

test:
	go test -v -race -cover ./...

lint:
	golangci-lint run

build: security-check
	go build -ldflags="-s -w" -trimpath -o glyphic ./cmd/glyphic

release:
	goreleaser release --clean

clean:
	rm -f glyphic
```

---

## Configuration File Support

```go
package config

import (
    "os"
    "path/filepath"
    
    "gopkg.in/yaml.v3"
)

type Config struct {
    Generation struct {
        WordCount    int      `yaml:"word_count"`
        Separator    string   `yaml:"separator"`
        MinLists     int      `yaml:"min_lists"`
        Capitalize   string   `yaml:"capitalize"`
        Numbers      int      `yaml:"numbers"`
        Specials     int      `yaml:"specials"`
        SpecialSet   string   `yaml:"special_set"`
    } `yaml:"generation"`
    
    Wordlists struct {
        CacheDir       string   `yaml:"cache_dir"`
        EnabledSources []string `yaml:"enabled_sources"`
        CustomLists    []string `yaml:"custom_lists"`
    } `yaml:"wordlists"`
    
    Exclusions struct {
        UseDefaults bool     `yaml:"use_defaults"`
        CustomFiles []string `yaml:"custom_files"`
        Words       []string `yaml:"words"`
    } `yaml:"exclusions"`
    
    Security struct {
        MinEntropy       int  `yaml:"min_entropy"`
        VerifyChecksums  bool `yaml:"verify_checksums"`
        ClipboardTimeout int  `yaml:"clipboard_timeout"`
    } `yaml:"security"`
    
    UI struct {
        ColorScheme  string `yaml:"color_scheme"`
        RevealSpeed  string `yaml:"reveal_speed"`
        ShowEntropy  bool   `yaml:"show_entropy"`
    } `yaml:"ui"`
}

func Load(path string) (*Config, error) {
    if path == "" {
        home, _ := os.UserHomeDir()
        path = filepath.Join(home, ".config", "glyphic", "config.yaml")
    }
    
    data, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return DefaultConfig(), nil
        }
        return nil, err
    }
    
    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }
    
    return &cfg, nil
}

func (c *Config) Save(path string) error {
    if path == "" {
        home, _ := os.UserHomeDir()
        path = filepath.Join(home, ".config", "glyphic", "config.yaml")
    }
    
    if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
        return err
    }
    
    data, err := yaml.Marshal(c)
    if err != nil {
        return err
    }
    
    return os.WriteFile(path, data, 0600)
}

func DefaultConfig() *Config {
    cfg := &Config{}
    cfg.Generation.WordCount = 6
    cfg.Generation.Separator = "-"
    cfg.Generation.MinLists = 3
    cfg.Generation.Capitalize = "none"
    cfg.Generation.SpecialSet = "!@#$%^&*-_+="
    cfg.Security.ClipboardTimeout = 30
    cfg.UI.ColorScheme = "matrix"
    cfg.UI.RevealSpeed = "normal"
    return cfg
}
```

### Example config.yaml

```yaml
generation:
  word_count: 6
  separator: "-"
  min_lists: 3
  capitalize: "none"
  numbers: 0
  specials: 0
  special_set: "!@#$%^&*-_+="

wordlists:
  cache_dir: ~/.local/share/glyphic/wordlists
  enabled_sources:
    - eff-large
    - eff-short-1
    - bip39-english
  custom_lists:
    - ~/my-wordlist.txt

exclusions:
  use_defaults: true
  custom_files:
    - ~/my-exclusions.txt
  words:
    - badword1
    - badword2

security:
  min_entropy: 60
  verify_checksums: true
  clipboard_timeout: 45

ui:
  color_scheme: matrix
  reveal_speed: normal
  show_entropy: true
```

---

## Security Summary

1. **Verified Sources**: All wordlists fetched via HTTPS with SHA-256 verification
2. **Certificate Pinning**: TLS certificate pinning for critical wordlist sources
3. **Multi-Source Mixing**: Words drawn from 3+ randomly-selected lists per password
4. **Cross-Verification**: Checksums verified from multiple independent sources
5. **Exclusion Lists**: Default exclusions prevent offensive/confusing words
6. **Memory Safety**: All passwords handled as `[]byte`, securely zeroed after use
7. **Memory Locking**: mlock() prevents password memory from swapping to disk
8. **Canary Tokens**: Detect buffer overruns and memory corruption
9. **Crypto Randomness**: Only `crypto/rand`, never `math/rand`
10. **Random Validation**: Startup checks verify PRNG is properly seeded
11. **Entropy Monitoring**: System entropy levels checked (Linux)
12. **Rate Limiting**: Prevents abuse of wordlist download endpoints
13. **Integrity Monitoring**: Periodic verification of cached wordlists
14. **Constant-Time Operations**: Timing-attack resistant exclusion checking
15. **No Persistence**: No logging of passwords, no temp files, no history
16. **Cache Security**: Wordlist cache uses restrictive permissions (0700/0600)
17. **Clipboard Security**: Auto-clear with configurable timeout
18. **Build Security**: Automated security checks (gosec, staticcheck) in build process
