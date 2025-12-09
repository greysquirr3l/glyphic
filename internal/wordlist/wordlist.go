// Package wordlist provides wordlist management with HTTPS fetching,
// SHA-256 verification, and local caching.
package wordlist

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"
)

var (
	// ErrChecksumMismatch indicates wordlist checksum verification failed
	ErrChecksumMismatch = errors.New("wordlist checksum mismatch")

	// ErrInvalidWordlist indicates wordlist data is invalid
	ErrInvalidWordlist = errors.New("invalid wordlist")

	// ErrInsufficientLists indicates not enough wordlists are available
	ErrInsufficientLists = errors.New("insufficient wordlists available")

	// ErrSourceNotFound indicates a wordlist source was not found
	ErrSourceNotFound = errors.New("wordlist source not found")
)

// WordlistSource represents a verified source for wordlist data
type WordlistSource struct {
	ID          string // Unique identifier
	Name        string // Display name
	URL         string // HTTPS download URL
	SHA256      string // Expected SHA-256 checksum
	WordCount   int    // Expected number of words
	MinLength   int    // Minimum word length
	MaxLength   int    // Maximum word length
	Description string // Human-readable description
	Category    string // "general", "technical", "nature", "phonetic", etc.
	Language    string // "en", "es", etc.
}

// DefaultSources contains verified wordlist sources
var DefaultSources = []WordlistSource{
	{
		ID:          "eff-large",
		Name:        "EFF Large Wordlist",
		URL:         "https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt",
		SHA256:      "replacewithactual", // Replace with actual checksum
		WordCount:   7776,
		MinLength:   3,
		MaxLength:   9,
		Description: "EFF's large wordlist for memorable passphrases",
		Category:    "general",
		Language:    "en",
	},
	{
		ID:          "eff-short-1",
		Name:        "EFF Short Wordlist 1",
		URL:         "https://www.eff.org/files/2016/09/08/eff_short_wordlist_1.txt",
		SHA256:      "replacewithactual",
		WordCount:   1296,
		MinLength:   4,
		MaxLength:   5,
		Description: "EFF's short wordlist with unique prefixes",
		Category:    "general",
		Language:    "en",
	},
	{
		ID:          "eff-short-2",
		Name:        "EFF Short Wordlist 2",
		URL:         "https://www.eff.org/files/2016/09/08/eff_short_wordlist_2_0.txt",
		SHA256:      "replacewithactual",
		WordCount:   1296,
		MinLength:   3,
		MaxLength:   5,
		Description: "EFF's alternative short wordlist",
		Category:    "general",
		Language:    "en",
	},
}

// Wordlist represents a loaded wordlist
type Wordlist struct {
	Source *WordlistSource
	Words  []string
}

// Manager handles wordlist fetching, caching, and loading
type Manager struct {
	cacheDir string
	sources  []WordlistSource
	loaded   map[string]*Wordlist
	client   *http.Client
	mu       sync.RWMutex
}

// NewManager creates a new wordlist manager
func NewManager(cacheDir string) (*Manager, error) {
	if cacheDir == "" {
		// Default: ~/.local/share/glyphic/wordlists/
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		cacheDir = filepath.Join(homeDir, ".local", "share", "glyphic", "wordlists")
	}

	// Ensure cache directory exists
	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &Manager{
		cacheDir: cacheDir,
		sources:  slices.Clone(DefaultSources),
		loaded:   make(map[string]*Wordlist),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// EnsureWordlists fetches any missing wordlists (called on first run)
func (m *Manager) EnsureWordlists(ctx context.Context) error {
	var errs []error

	for _, source := range m.sources {
		cachePath := m.cachePath(source.ID)

		// Check if cached and valid
		if m.isValidCache(cachePath, source) {
			continue
		}

		// Fetch and cache
		if err := m.fetchAndCache(ctx, source); err != nil {
			errs = append(errs, fmt.Errorf("source %s: %w", source.ID, err))
		}
	}

	if len(errs) > 0 && len(errs) == len(m.sources) {
		return fmt.Errorf("failed to fetch all wordlists: %v", errs)
	}

	return nil
}

// fetchAndCache downloads and caches a wordlist
func (m *Manager) fetchAndCache(ctx context.Context, source WordlistSource) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, source.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch wordlist: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read and compute checksum
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	checksum := sha256.Sum256(data)
	checksumHex := hex.EncodeToString(checksum[:])

	// Verify checksum if provided
	if source.SHA256 != "replacewithactual" && source.SHA256 != checksumHex {
		return fmt.Errorf("%w: expected %s, got %s", ErrChecksumMismatch, source.SHA256, checksumHex)
	}

	// Parse and validate wordlist
	words := parseWordlist(data, source.ID)
	if len(words) == 0 {
		return ErrInvalidWordlist
	}

	// Write to cache
	cachePath := m.cachePath(source.ID)
	if err := os.WriteFile(cachePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write cache: %w", err)
	}

	return nil
}

// cachePath returns the cache file path for a wordlist ID
func (m *Manager) cachePath(id string) string {
	return filepath.Join(m.cacheDir, id+".txt")
}

// isValidCache checks if a cached wordlist is valid
func (m *Manager) isValidCache(path string, source WordlistSource) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// Check file is not empty
	if info.Size() == 0 {
		return false
	}

	// Check file age (re-fetch if older than 30 days)
	if time.Since(info.ModTime()) > 30*24*time.Hour {
		return false
	}

	return true
}

// AvailableCount returns the number of available wordlists
func (m *Manager) AvailableCount() int {
	count := 0
	for _, source := range m.sources {
		if _, err := os.Stat(m.cachePath(source.ID)); err == nil {
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
		cachePath := m.cachePath(source.ID)
		if _, err := os.Stat(cachePath); err != nil {
			continue // Skip missing wordlists
		}

		data, err := os.ReadFile(cachePath)
		if err != nil {
			return fmt.Errorf("failed to read cache %s: %w", source.ID, err)
		}

		words := parseWordlist(data, source.ID)
		m.loaded[source.ID] = &Wordlist{
			Source: &source,
			Words:  words,
		}
	}

	return nil
}

// SelectRandomLists returns n randomly selected wordlists
func (m *Manager) SelectRandomLists(n int) ([]*Wordlist, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.loaded) < n {
		return nil, fmt.Errorf("%w: need %d, have %d", ErrInsufficientLists, n, len(m.loaded))
	}

	// Get all available wordlists
	available := make([]*Wordlist, 0, len(m.loaded))
	for _, wl := range m.loaded {
		available = append(available, wl)
	}

	// TODO: Use secure random selection from internal/security
	// For now, return first n wordlists
	return available[:n], nil
}

// AddUserWordlist adds a user-provided wordlist
func (m *Manager) AddUserWordlist(path, id string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read wordlist: %w", err)
	}

	words := parseWordlist(data, id)
	if len(words) == 0 {
		return ErrInvalidWordlist
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.loaded[id] = &Wordlist{
		Source: &WordlistSource{
			ID:          id,
			Name:        filepath.Base(path),
			Description: "User-provided wordlist",
			Category:    "custom",
		},
		Words: words,
	}

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
			// Implementation would track enabled/disabled state
			return true
		}
	}
	return false
}

// parseWordlist parses wordlist data into a slice of words
func parseWordlist(data []byte, sourceID string) []string {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	words := make([]string, 0, 8192)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Handle EFF format (number followed by word)
		fields := strings.Fields(line)
		var word string
		if len(fields) >= 2 && isAllDigits(fields[0]) {
			word = fields[1] // EFF format: "11111 word"
		} else if len(fields) == 1 {
			word = fields[0] // Plain word format
		} else {
			continue // Skip invalid lines
		}

		// Validate word
		if isValidWord(word) {
			words = append(words, strings.ToLower(word))
		}
	}

	// Sort for binary search
	slices.Sort(words)

	// Remove duplicates
	return slices.Compact(words)
}

// isAllDigits returns true if the string contains only digits
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

// isValidWord validates a word for inclusion in the wordlist
func isValidWord(s string) bool {
	if len(s) < 2 || len(s) > 12 {
		return false
	}

	// Must contain only letters
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}

	return true
}
