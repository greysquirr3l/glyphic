package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/greysquirr3l/glyphic/internal/wordlist"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestGenerator(t *testing.T) *Generator {
	t.Helper()

	// Create temporary wordlists
	tmpDir := t.TempDir()

	wordlist1Path := filepath.Join(tmpDir, "list1.txt")
	wordlist2Path := filepath.Join(tmpDir, "list2.txt")
	wordlist3Path := filepath.Join(tmpDir, "list3.txt")

	words1 := []string{"apple", "banana", "cherry", "date", "elderberry"}
	words2 := []string{"falcon", "goose", "hawk", "ibis", "jay"}
	words3 := []string{"kale", "lettuce", "mustard", "napa", "okra"}

	err := os.WriteFile(wordlist1Path, []byte(strings.Join(words1, "\n")), 0600)
	require.NoError(t, err)
	err = os.WriteFile(wordlist2Path, []byte(strings.Join(words2, "\n")), 0600)
	require.NoError(t, err)
	err = os.WriteFile(wordlist3Path, []byte(strings.Join(words3, "\n")), 0600)
	require.NoError(t, err)

	// Create manager and load wordlists
	cacheDir := filepath.Join(tmpDir, "cache")
	manager, err := wordlist.NewManager(cacheDir)
	require.NoError(t, err)

	err = manager.AddUserWordlist(wordlist1Path, "test1")
	require.NoError(t, err)
	err = manager.AddUserWordlist(wordlist2Path, "test2")
	require.NoError(t, err)
	err = manager.AddUserWordlist(wordlist3Path, "test3")
	require.NoError(t, err)

	// Create exclusion list (no defaults for predictable testing)
	exclusions := wordlist.NewExclusionList(false)

	return New(manager, exclusions)
}

func TestOptionsValidate(t *testing.T) {
	tests := []struct {
		name    string
		opts    Options
		wantErr error
	}{
		{
			name:    "valid default options",
			opts:    DefaultOptions,
			wantErr: nil,
		},
		{
			name: "word count too low",
			opts: Options{
				WordCount:    2,
				MinWordlists: 1,
			},
			wantErr: ErrInvalidWordCount,
		},
		{
			name: "word count too high",
			opts: Options{
				WordCount:    13,
				MinWordlists: 1,
			},
			wantErr: ErrInvalidWordCount,
		},
		{
			name: "number count invalid",
			opts: Options{
				WordCount:    6,
				AddNumbers:   true,
				NumberCount:  5,
				MinWordlists: 1,
			},
			wantErr: ErrInvalidNumberCount,
		},
		{
			name: "special count invalid",
			opts: Options{
				WordCount:    6,
				AddSpecial:   true,
				SpecialCount: 0,
				MinWordlists: 1,
			},
			wantErr: ErrInvalidSpecialCount,
		},
		{
			name: "min wordlists invalid",
			opts: Options{
				WordCount:    6,
				MinWordlists: 0,
			},
			wantErr: ErrInvalidMinWordlists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateBasic(t *testing.T) {
	gen := setupTestGenerator(t)

	opts := Options{
		WordCount:      6,
		Capitalization: CapNone,
		Separator:      SepDash,
		MinWordlists:   3,
	}

	password, err := gen.Generate(opts)
	assert.NoError(t, err)
	assert.NotEmpty(t, password)

	// Should have 6 words separated by dashes
	words := strings.Split(password, "-")
	assert.Len(t, words, 6)

	// Each word should be lowercase (CapNone)
	for _, word := range words {
		assert.Equal(t, strings.ToLower(word), word)
	}
}

func TestGenerateWithNumbers(t *testing.T) {
	gen := setupTestGenerator(t)

	opts := Options{
		WordCount:      4,
		Capitalization: CapNone,
		Separator:      SepDash,
		AddNumbers:     true,
		NumberCount:    3,
		MinWordlists:   2,
	}

	password, err := gen.Generate(opts)
	assert.NoError(t, err)
	assert.NotEmpty(t, password)

	// Should end with 3 digits
	assert.Regexp(t, `\d{3}$`, password)
}

func TestGenerateWithSpecialChars(t *testing.T) {
	gen := setupTestGenerator(t)

	opts := Options{
		WordCount:      4,
		Capitalization: CapNone,
		Separator:      SepDash,
		AddSpecial:     true,
		SpecialCount:   2,
		MinWordlists:   2,
	}

	password, err := gen.Generate(opts)
	assert.NoError(t, err)
	assert.NotEmpty(t, password)

	// Should end with 2 special characters
	// Count special chars at the end
	specialCount := 0
	for i := len(password) - 1; i >= 0; i-- {
		r := rune(password[i])
		isSpecial := false
		for _, sc := range SpecialChars {
			if r == sc {
				isSpecial = true
				break
			}
		}
		if !isSpecial {
			break
		}
		specialCount++
	}
	assert.Equal(t, 2, specialCount)
}

func TestCapitalizationModes(t *testing.T) {
	gen := setupTestGenerator(t)

	tests := []struct {
		name  string
		mode  CapitalizationMode
		check func(t *testing.T, password string)
	}{
		{
			name: "CapNone",
			mode: CapNone,
			check: func(t *testing.T, password string) {
				words := strings.Split(password, "-")
				for _, word := range words {
					assert.Equal(t, strings.ToLower(word), word)
				}
			},
		},
		{
			name: "CapFirst",
			mode: CapFirst,
			check: func(t *testing.T, password string) {
				words := strings.Split(password, "-")
				for _, word := range words {
					if len(word) > 0 {
						assert.True(t, word[0] >= 'A' && word[0] <= 'Z', "first letter should be uppercase")
					}
				}
			},
		},
		{
			name: "CapAll",
			mode: CapAll,
			check: func(t *testing.T, password string) {
				words := strings.Split(password, "-")
				for _, word := range words {
					assert.Equal(t, strings.ToUpper(word), word)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{
				WordCount:      4,
				Capitalization: tt.mode,
				Separator:      SepDash,
				MinWordlists:   2,
			}

			password, err := gen.Generate(opts)
			assert.NoError(t, err)
			assert.NotEmpty(t, password)

			tt.check(t, password)
		})
	}
}

func TestSeparatorModes(t *testing.T) {
	gen := setupTestGenerator(t)

	tests := []struct {
		name      string
		separator SeparatorMode
		customSep string
		expected  string
	}{
		{
			name:      "SepNone",
			separator: SepNone,
			expected:  "",
		},
		{
			name:      "SepSpace",
			separator: SepSpace,
			expected:  " ",
		},
		{
			name:      "SepDash",
			separator: SepDash,
			expected:  "-",
		},
		{
			name:      "SepUnderscore",
			separator: SepUnderscore,
			expected:  "_",
		},
		{
			name:      "SepCustom",
			separator: SepCustom,
			customSep: "|",
			expected:  "|",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{
				WordCount:      3,
				Capitalization: CapNone,
				Separator:      tt.separator,
				CustomSep:      tt.customSep,
				MinWordlists:   2,
			}

			password, err := gen.Generate(opts)
			assert.NoError(t, err)
			assert.NotEmpty(t, password)

			if tt.expected != "" {
				assert.Contains(t, password, tt.expected)
			}
		})
	}
}

func TestGenerateMultiple(t *testing.T) {
	gen := setupTestGenerator(t)

	opts := Options{
		WordCount:      4,
		Capitalization: CapFirst,
		Separator:      SepDash,
		MinWordlists:   2,
	}

	count := 10
	passwords, err := gen.GenerateMultiple(count, opts)
	assert.NoError(t, err)
	assert.Len(t, passwords, count)

	// Each password should be unique (very high probability)
	seen := make(map[string]bool)
	for _, pwd := range passwords {
		assert.False(t, seen[pwd], "duplicate password generated")
		seen[pwd] = true
	}
}

func TestExclusionFiltering(t *testing.T) {
	gen := setupTestGenerator(t)

	// Add some words to exclusion list
	gen.exclusions.Add("apple", "banana", "falcon")

	opts := Options{
		WordCount:      4,
		Capitalization: CapNone,
		Separator:      SepDash,
		MinWordlists:   3,
	}

	// Generate multiple passwords and verify excluded words don't appear
	for range 20 {
		password, err := gen.Generate(opts)
		assert.NoError(t, err)

		assert.NotContains(t, password, "apple")
		assert.NotContains(t, password, "banana")
		assert.NotContains(t, password, "falcon")
	}
}

func TestEstimateEntropy(t *testing.T) {
	gen := setupTestGenerator(t)

	tests := []struct {
		name        string
		opts        Options
		expectRange [2]float64 // min, max expected entropy
	}{
		{
			name: "basic 6 words",
			opts: Options{
				WordCount:    6,
				MinWordlists: 3,
			},
			expectRange: [2]float64{10.0, 50.0}, // rough estimate
		},
		{
			name: "with numbers",
			opts: Options{
				WordCount:    6,
				AddNumbers:   true,
				NumberCount:  2,
				MinWordlists: 3,
			},
			expectRange: [2]float64{15.0, 60.0}, // should be higher
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entropy, err := gen.EstimateEntropy(tt.opts)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, entropy, tt.expectRange[0])
			assert.LessOrEqual(t, entropy, tt.expectRange[1])
		})
	}
}

func TestApplyCapitalization(t *testing.T) {
	words := []string{"hello", "world", "test"}

	tests := []struct {
		name  string
		mode  CapitalizationMode
		check func(t *testing.T, result []string)
	}{
		{
			name: "CapNone",
			mode: CapNone,
			check: func(t *testing.T, result []string) {
				for _, word := range result {
					assert.Equal(t, strings.ToLower(word), word)
				}
			},
		},
		{
			name: "CapFirst",
			mode: CapFirst,
			check: func(t *testing.T, result []string) {
				assert.Equal(t, "Hello", result[0])
				assert.Equal(t, "World", result[1])
				assert.Equal(t, "Test", result[2])
			},
		},
		{
			name: "CapAll",
			mode: CapAll,
			check: func(t *testing.T, result []string) {
				assert.Equal(t, "HELLO", result[0])
				assert.Equal(t, "WORLD", result[1])
				assert.Equal(t, "TEST", result[2])
			},
		},
		{
			name: "CapAlternating",
			mode: CapAlternating,
			check: func(t *testing.T, result []string) {
				assert.Equal(t, "HELLO", result[0])
				assert.Equal(t, "world", result[1])
				assert.Equal(t, "TEST", result[2])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := applyCapitalization(words, tt.mode)
			assert.NoError(t, err)
			tt.check(t, result)
		})
	}
}

func TestGenerateRandomNumbers(t *testing.T) {
	for count := 1; count <= 4; count++ {
		t.Run(fmt.Sprintf("count=%d", count), func(t *testing.T) {
			numbers, err := generateRandomNumbers(count)
			assert.NoError(t, err)
			assert.Len(t, numbers, count)

			// Should only contain digits
			for _, r := range numbers {
				assert.True(t, r >= '0' && r <= '9')
			}
		})
	}
}

func TestGenerateRandomSpecialChars(t *testing.T) {
	for count := 1; count <= 4; count++ {
		t.Run(fmt.Sprintf("count=%d", count), func(t *testing.T) {
			specials, err := generateRandomSpecialChars(count)
			assert.NoError(t, err)
			assert.Len(t, specials, count)

			// Each character should be in SpecialChars
			for _, r := range specials {
				found := false
				for _, sc := range SpecialChars {
					if r == sc {
						found = true
						break
					}
				}
				assert.True(t, found, "character %c not in SpecialChars", r)
			}
		})
	}
}
