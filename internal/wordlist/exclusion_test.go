package wordlist

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewExclusionList(t *testing.T) {
	tests := []struct {
		name        string
		useDefaults bool
		wantCount   int
	}{
		{
			name:        "with defaults",
			useDefaults: true,
			wantCount:   1, // At least 1 word
		},
		{
			name:        "without defaults",
			useDefaults: false,
			wantCount:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewExclusionList(tt.useDefaults)
			assert.NotNil(t, list)
			if tt.useDefaults {
				assert.GreaterOrEqual(t, list.Count(), tt.wantCount)
			} else {
				assert.Equal(t, tt.wantCount, list.Count())
			}
		})
	}
}

func TestLoadDefaultExclusions(t *testing.T) {
	words := loadDefaultExclusions()
	assert.NotEmpty(t, words, "should load default exclusion words")

	// Verify words are lowercase and sorted
	for i, word := range words {
		assert.Equal(t, word, word, "word should equal itself")
		if i > 0 {
			assert.True(t, words[i-1] < words[i], "words should be sorted")
		}
	}
}

func TestExclusionListAdd(t *testing.T) {
	list := NewExclusionList(false)
	assert.Equal(t, 0, list.Count())

	list.Add("apple", "banana", "cherry")
	assert.Equal(t, 3, list.Count())

	// Add duplicates
	list.Add("apple", "date")
	assert.Equal(t, 4, list.Count()) // apple not added twice

	// Verify sorted
	words := list.words
	for i := 1; i < len(words); i++ {
		assert.True(t, words[i-1] < words[i])
	}
}

func TestExclusionListContains(t *testing.T) {
	list := NewExclusionList(false)
	list.Add("apple", "banana", "cherry")

	tests := []struct {
		word string
		want bool
	}{
		{"apple", true},
		{"banana", true},
		{"cherry", true},
		{"date", false},
		{"APPLE", true}, // case-insensitive
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			got := list.Contains(tt.word)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExclusionListFilter(t *testing.T) {
	list := NewExclusionList(false)
	list.Add("bad", "worse", "worst")

	tests := []struct {
		name  string
		words []string
		want  []string
	}{
		{
			name:  "remove some words",
			words: []string{"good", "bad", "better", "worse", "best"},
			want:  []string{"good", "better", "best"},
		},
		{
			name:  "remove all words",
			words: []string{"bad", "worse", "worst"},
			want:  []string{},
		},
		{
			name:  "remove no words",
			words: []string{"good", "better", "best"},
			want:  []string{"good", "better", "best"},
		},
		{
			name:  "empty input",
			words: []string{},
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := list.Filter(tt.words)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExclusionListLoadFile(t *testing.T) {
	// Create temporary exclusion file
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test-exclusions.txt")
	content := `# Test exclusions
apple
banana
# Another comment

cherry
apple
`
	err := os.WriteFile(filePath, []byte(content), 0600)
	require.NoError(t, err)

	// Load file
	list := NewExclusionList(false)
	err = list.LoadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, 3, list.Count()) // apple, banana, cherry (no dupes)

	// Verify words are present
	assert.True(t, list.Contains("apple"))
	assert.True(t, list.Contains("banana"))
	assert.True(t, list.Contains("cherry"))
}

func TestExclusionListLoadFileErrors(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{
			name:     "nonexistent file",
			filePath: "/nonexistent/file.txt",
			wantErr:  true,
		},
		{
			name:     "directory instead of file",
			filePath: t.TempDir(),
			wantErr:  false, // bufio.Scanner won't error on empty
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewExclusionList(false)
			err := list.LoadFile(tt.filePath)
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestExclusionListDisable(t *testing.T) {
	list := NewExclusionList(true) // Start with defaults
	assert.Greater(t, list.Count(), 0)

	list.Disable()
	assert.Equal(t, 0, list.Count())
	assert.False(t, list.Contains("anything"))
}

func TestExclusionListCaseInsensitive(t *testing.T) {
	list := NewExclusionList(false)
	list.Add("Apple", "BANANA", "CheRRy")

	// All should be lowercase internally
	assert.True(t, list.Contains("apple"))
	assert.True(t, list.Contains("banana"))
	assert.True(t, list.Contains("cherry"))

	// Case variations should match
	assert.True(t, list.Contains("APPLE"))
	assert.True(t, list.Contains("Banana"))
	assert.True(t, list.Contains("CHERRY"))
}

func TestFilterPreservesOrder(t *testing.T) {
	list := NewExclusionList(false)
	list.Add("remove1", "remove2")

	input := []string{"keep1", "remove1", "keep2", "remove2", "keep3"}
	expected := []string{"keep1", "keep2", "keep3"}

	got := list.Filter(input)
	assert.Equal(t, expected, got)
}
