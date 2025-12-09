package generator

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/greysquirr3l/glyphic/internal/security"
	"github.com/greysquirr3l/glyphic/internal/wordlist"
)

// CapitalizationMode defines how words should be capitalized
type CapitalizationMode int

const (
	CapNone        CapitalizationMode = iota // no capitalization
	CapFirst                                 // capitalize first letter of each word
	CapRandom                                // capitalize random letters
	CapAll                                   // all uppercase
	CapAlternating                           // alternate between upper and lower case words
)

// SeparatorMode defines how words should be separated
type SeparatorMode int

const (
	SepNone       SeparatorMode = iota // no separator
	SepSpace                           // space between words
	SepDash                            // dash between words
	SepUnderscore                      // underscore between words
	SepCustom                          // custom separator
)

// Options configures password generation
type Options struct {
	WordCount      int                // Number of words (3-12)
	Capitalization CapitalizationMode // How to capitalize
	AddNumbers     bool               // Add random numbers
	NumberCount    int                // How many numbers to add (1-4)
	AddSpecial     bool               // Add special characters
	SpecialCount   int                // How many special chars to add (1-4)
	Separator      SeparatorMode      // Word separator style
	CustomSep      string             // Custom separator if SepCustom
	MinWordlists   int                // Minimum different wordlists to use (default 3)
}

// DefaultOptions provides secure default settings
var DefaultOptions = Options{
	WordCount:      6,
	Capitalization: CapFirst,
	AddNumbers:     false,
	NumberCount:    2,
	AddSpecial:     false,
	SpecialCount:   1,
	Separator:      SepDash,
	MinWordlists:   3,
}

// SpecialChars is the set of allowed special characters
var SpecialChars = []rune{'!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '-', '_', '=', '+', '[', ']', '{', '}', '|', ';', ':', ',', '.', '?'}

var (
	ErrInvalidWordCount    = errors.New("word count must be between 3 and 12")
	ErrInvalidNumberCount  = errors.New("number count must be between 1 and 4")
	ErrInvalidSpecialCount = errors.New("special character count must be between 1 and 4")
	ErrInvalidMinWordlists = errors.New("minimum wordlists must be at least 1")
	ErrNotEnoughWordlists  = errors.New("not enough different wordlists available")
	ErrNoWordsAvailable    = errors.New("no words available after exclusion filtering")
)

// Generator generates passwords using wordlists
type Generator struct {
	manager    *wordlist.Manager
	exclusions *wordlist.ExclusionList
}

// New creates a new password generator
func New(manager *wordlist.Manager, exclusions *wordlist.ExclusionList) *Generator {
	return &Generator{
		manager:    manager,
		exclusions: exclusions,
	}
}

// Validate checks if options are valid
func (o *Options) Validate() error {
	if o.WordCount < 3 || o.WordCount > 12 {
		return ErrInvalidWordCount
	}
	if o.AddNumbers && (o.NumberCount < 1 || o.NumberCount > 4) {
		return ErrInvalidNumberCount
	}
	if o.AddSpecial && (o.SpecialCount < 1 || o.SpecialCount > 4) {
		return ErrInvalidSpecialCount
	}
	if o.MinWordlists < 1 {
		return ErrInvalidMinWordlists
	}
	return nil
}

// Generate creates a single password with the given options
func (g *Generator) Generate(opts Options) (string, error) {
	if err := opts.Validate(); err != nil {
		return "", fmt.Errorf("invalid options: %w", err)
	}

	// Select random wordlists (at least MinWordlists different ones)
	numLists := opts.MinWordlists
	if opts.WordCount < numLists {
		numLists = opts.WordCount
	}

	lists, err := g.manager.SelectRandomLists(numLists)
	if err != nil {
		return "", fmt.Errorf("failed to select wordlists: %w", err)
	}

	if len(lists) < numLists {
		return "", fmt.Errorf("%w: need %d, got %d", ErrNotEnoughWordlists, numLists, len(lists))
	}

	// Select words from different wordlists
	words := make([]string, opts.WordCount)
	for i := range opts.WordCount {
		// Round-robin through wordlists
		listIdx := i % len(lists)
		list := lists[listIdx]

		// Filter words by exclusion list
		availableWords := g.exclusions.Filter(list.Words)
		if len(availableWords) == 0 {
			return "", fmt.Errorf("%w: wordlist %s has no available words after filtering", ErrNoWordsAvailable, list.Source.ID)
		}

		// Select random word
		wordIdx, err := security.SecureRandomIndex(len(availableWords))
		if err != nil {
			return "", fmt.Errorf("failed to select random word: %w", err)
		}

		words[i] = availableWords[wordIdx]
	}

	// Apply capitalization
	words, err = applyCapitalization(words, opts.Capitalization)
	if err != nil {
		return "", fmt.Errorf("failed to apply capitalization: %w", err)
	}

	// Join words with separator
	separator := getSeparator(opts.Separator, opts.CustomSep)
	password := strings.Join(words, separator)

	// Add numbers if requested
	if opts.AddNumbers {
		numbers, err := generateRandomNumbers(opts.NumberCount)
		if err != nil {
			return "", fmt.Errorf("failed to generate numbers: %w", err)
		}
		password += numbers
	}

	// Add special characters if requested
	if opts.AddSpecial {
		specials, err := generateRandomSpecialChars(opts.SpecialCount)
		if err != nil {
			return "", fmt.Errorf("failed to generate special characters: %w", err)
		}
		password += specials
	}

	return password, nil
}

// GenerateMultiple creates multiple passwords
func (g *Generator) GenerateMultiple(count int, opts Options) ([]string, error) {
	if count < 1 {
		return nil, errors.New("count must be at least 1")
	}

	passwords := make([]string, count)
	for i := range count {
		pwd, err := g.Generate(opts)
		if err != nil {
			return nil, fmt.Errorf("failed to generate password %d: %w", i+1, err)
		}
		passwords[i] = pwd
	}

	return passwords, nil
}

// applyCapitalization applies the specified capitalization mode to words
func applyCapitalization(words []string, mode CapitalizationMode) ([]string, error) {
	result := make([]string, len(words))

	for i, word := range words {
		switch mode {
		case CapNone:
			result[i] = strings.ToLower(word)

		case CapFirst:
			if len(word) == 0 {
				result[i] = word
			} else {
				result[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
			}

		case CapRandom:
			runes := []rune(strings.ToLower(word))
			for j := range runes {
				// 50% chance to capitalize each character
				shouldCap, err := security.SecureRandomIndex(2)
				if err != nil {
					return nil, err
				}
				if shouldCap == 1 {
					runes[j] = []rune(strings.ToUpper(string(runes[j])))[0]
				}
			}
			result[i] = string(runes)

		case CapAll:
			result[i] = strings.ToUpper(word)

		case CapAlternating:
			if i%2 == 0 {
				result[i] = strings.ToUpper(word)
			} else {
				result[i] = strings.ToLower(word)
			}

		default:
			result[i] = word
		}
	}

	return result, nil
}

// getSeparator returns the appropriate separator string
func getSeparator(mode SeparatorMode, custom string) string {
	switch mode {
	case SepNone:
		return ""
	case SepSpace:
		return " "
	case SepDash:
		return "-"
	case SepUnderscore:
		return "_"
	case SepCustom:
		return custom
	default:
		return ""
	}
}

// generateRandomNumbers generates N random digits
func generateRandomNumbers(count int) (string, error) {
	var sb strings.Builder
	for range count {
		digit, err := security.SecureRandomIndex(10)
		if err != nil {
			return "", err
		}
		sb.WriteString(strconv.Itoa(digit))
	}
	return sb.String(), nil
}

// generateRandomSpecialChars generates N random special characters
func generateRandomSpecialChars(count int) (string, error) {
	var sb strings.Builder
	for range count {
		idx, err := security.SecureRandomIndex(len(SpecialChars))
		if err != nil {
			return "", err
		}
		sb.WriteRune(SpecialChars[idx])
	}
	return sb.String(), nil
}

// EstimateEntropy calculates the approximate entropy bits for given options
func (g *Generator) EstimateEntropy(opts Options) (float64, error) {
	if err := opts.Validate(); err != nil {
		return 0, err
	}

	// Get average wordlist size
	lists, err := g.manager.SelectRandomLists(opts.MinWordlists)
	if err != nil {
		return 0, err
	}

	var totalWords int
	for _, list := range lists {
		availableWords := g.exclusions.Filter(list.Words)
		totalWords += len(availableWords)
	}

	if totalWords == 0 {
		return 0, ErrNoWordsAvailable
	}

	avgWordlistSize := float64(totalWords) / float64(len(lists))

	// Calculate entropy: log2(possibilities)
	// For words: word_count * log2(avg_wordlist_size)
	entropy := float64(opts.WordCount) * math.Log2(avgWordlistSize)

	// Add entropy for numbers
	if opts.AddNumbers {
		entropy += float64(opts.NumberCount) * math.Log2(10)
	}

	// Add entropy for special characters
	if opts.AddSpecial {
		entropy += float64(opts.SpecialCount) * math.Log2(float64(len(SpecialChars)))
	}

	return entropy, nil
}
