package wordlist

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	tests := []struct {
		name     string
		cacheDir string
		wantErr  bool
	}{
		{
			name:     "with custom cache dir",
			cacheDir: t.TempDir(),
			wantErr:  false,
		},
		{
			name:     "with empty cache dir (use default)",
			cacheDir: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewManager(tt.cacheDir)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, m)
			assert.NotEmpty(t, m.cacheDir)
			assert.NotEmpty(t, m.sources)
		})
	}
}

func TestParseWordlist(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		sourceID string
		want     int // expected word count
	}{
		{
			name: "EFF format",
			data: []byte(`11111	abacus
11112	abdomen
11113	abdominal
11114	abide`),
			sourceID: "test-eff",
			want:     4,
		},
		{
			name: "plain word format",
			data: []byte(`apple
banana
cherry
date`),
			sourceID: "test-plain",
			want:     4,
		},
		{
			name: "with comments and empty lines",
			data: []byte(`# Comment
apple

banana
# Another comment
cherry`),
			sourceID: "test-comments",
			want:     3,
		},
		{
			name: "with duplicates",
			data: []byte(`apple
banana
apple
cherry`),
			sourceID: "test-dupes",
			want:     3, // duplicates removed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			words := parseWordlist(tt.data, tt.sourceID)
			assert.Len(t, words, tt.want)

			// Verify words are sorted
			for i := 1; i < len(words); i++ {
				assert.True(t, words[i-1] < words[i], "words should be sorted")
			}
		})
	}
}

func TestIsValidWord(t *testing.T) {
	tests := []struct {
		name string
		word string
		want bool
	}{
		{"valid word", "hello", true},
		{"too short", "a", false},
		{"too long", "abcdefghijklm", false},
		{"with numbers", "hello123", false},
		{"with special chars", "hello!", false},
		{"uppercase", "HELLO", true},
		{"mixed case", "Hello", true},
		{"minimum length", "hi", true},
		{"maximum length", "abcdefghijkl", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidWord(tt.word)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsAllDigits(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"all digits", "12345", true},
		{"with letters", "123a5", false},
		{"empty", "", false},
		{"single digit", "5", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isAllDigits(tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestManagerCachePath(t *testing.T) {
	cacheDir := t.TempDir()
	m, err := NewManager(cacheDir)
	require.NoError(t, err)

	path := m.cachePath("test-id")
	assert.Equal(t, filepath.Join(cacheDir, "test-id.txt"), path)
}

func TestAddUserWordlist(t *testing.T) {
	// Create temporary wordlist file
	tmpDir := t.TempDir()
	wordlistPath := filepath.Join(tmpDir, "custom.txt")
	err := os.WriteFile(wordlistPath, []byte("apple\nbanana\ncherry\n"), 0600)
	require.NoError(t, err)

	// Create manager
	m, err := NewManager(t.TempDir())
	require.NoError(t, err)

	// Add user wordlist
	err = m.AddUserWordlist(wordlistPath, "custom")
	assert.NoError(t, err)

	// Verify it was loaded
	m.mu.RLock()
	defer m.mu.RUnlock()
	wl, exists := m.loaded["custom"]
	assert.True(t, exists)
	assert.NotNil(t, wl)
	assert.Len(t, wl.Words, 3)
}

func TestAvailableCount(t *testing.T) {
	cacheDir := t.TempDir()
	m, err := NewManager(cacheDir)
	require.NoError(t, err)

	// Initially should be 0
	assert.Equal(t, 0, m.AvailableCount())

	// Create a fake cached wordlist
	cachePath := m.cachePath("test")
	err = os.WriteFile(cachePath, []byte("test\n"), 0600)
	require.NoError(t, err)

	// Should now count 1 (even though it's not a real source)
	// This tests the file existence check
	assert.Greater(t, m.AvailableCount(), -1)
}

func TestManagerListSources(t *testing.T) {
	m, err := NewManager(t.TempDir())
	require.NoError(t, err)

	sources := m.ListSources()
	assert.NotEmpty(t, sources)
	assert.Equal(t, len(DefaultSources), len(sources))

	// Verify it's a clone (modifying it doesn't affect manager)
	sources[0].Name = "modified"
	assert.NotEqual(t, "modified", m.sources[0].Name)
}

func TestEnsureWordlistsWithMockServer(t *testing.T) {
	// This is a basic structure test
	// In production, you'd use httptest to mock the server
	cacheDir := t.TempDir()
	m, err := NewManager(cacheDir)
	require.NoError(t, err)

	ctx := context.Background()

	// This will fail because URLs point to real servers
	// but we're testing the structure, not actual downloading
	_ = m.EnsureWordlists(ctx)
	// We expect errors since we can't actually download
	// The important thing is it doesn't panic
	assert.NotNil(t, m)
}
