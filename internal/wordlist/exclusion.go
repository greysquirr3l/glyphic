// Package wordlist - exclusion list functionality
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

// ExclusionList manages words to exclude from password generation
type ExclusionList struct {
	words []string // Sorted for binary search
	mu    sync.RWMutex
}

// NewExclusionList creates a new exclusion list
func NewExclusionList(useDefaults bool) *ExclusionList {
	el := &ExclusionList{}
	if useDefaults {
		el.words = loadDefaultExclusions()
	}
	return el
}

// loadDefaultExclusions loads embedded default exclusion lists
func loadDefaultExclusions() []string {
	var words []string

	// Load all embedded exclusion files
	exclusionFiles := []string{
		"exclusions/profanity.txt",
		"exclusions/slurs.txt",
		"exclusions/sensitive.txt",
		"exclusions/confusing.txt",
	}

	for _, filename := range exclusionFiles {
		data, err := embeddedExclusions.ReadFile(filename)
		if err != nil {
			continue // Skip missing files
		}

		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" && !strings.HasPrefix(line, "#") {
				words = append(words, strings.ToLower(line))
			}
		}
	}

	// Sort and remove duplicates
	slices.Sort(words)
	return slices.Compact(words)
}

// LoadFile loads exclusions from a file
func (e *ExclusionList) LoadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	e.mu.Lock()
	defer e.mu.Unlock()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			e.words = append(e.words, strings.ToLower(line))
		}
	}

	// Re-sort and compact after adding new words
	slices.Sort(e.words)
	e.words = slices.Compact(e.words)

	return scanner.Err()
}

// Add adds words to the exclusion list
func (e *ExclusionList) Add(words ...string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, word := range words {
		e.words = append(e.words, strings.ToLower(word))
	}

	slices.Sort(e.words)
	e.words = slices.Compact(e.words)
}

// Contains checks if a word is in the exclusion list
func (e *ExclusionList) Contains(word string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	_, found := slices.BinarySearch(e.words, strings.ToLower(word))
	return found
}

// Filter removes excluded words from a slice
func (e *ExclusionList) Filter(words []string) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return slices.DeleteFunc(words, func(word string) bool {
		_, found := slices.BinarySearch(e.words, strings.ToLower(word))
		return found
	})
}

// Count returns the number of excluded words
func (e *ExclusionList) Count() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.words)
}

// Disable clears the exclusion list
func (e *ExclusionList) Disable() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.words = nil
}
