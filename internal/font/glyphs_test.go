package font

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatrixGlyphsNotEmpty(t *testing.T) {
	assert.NotEmpty(t, MatrixGlyphs.ASCII, "ASCII glyphs should not be empty")
	assert.NotEmpty(t, MatrixGlyphs.Basic, "Basic glyphs should not be empty")
	assert.NotEmpty(t, MatrixGlyphs.Extended, "Extended glyphs should not be empty")
}

func TestMatrixGlyphsValid(t *testing.T) {
	tests := []struct {
		name   string
		glyphs []rune
	}{
		{"ASCII", MatrixGlyphs.ASCII},
		{"Basic", MatrixGlyphs.Basic},
		{"Extended", MatrixGlyphs.Extended},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, r := range tt.glyphs {
				// Each glyph should be valid UTF-8
				assert.NotEqual(t, utf8.RuneError, r, "glyph should be valid UTF-8")

				// Should be printable
				assert.True(t, IsSafe(r), "glyph %U should be safe", r)
			}
		})
	}
}

func TestASCIIGlyphsInRange(t *testing.T) {
	// ASCII glyphs should be in printable range (33-126)
	for _, r := range MatrixGlyphs.ASCII {
		assert.GreaterOrEqual(t, int(r), 33, "ASCII glyph below printable range")
		assert.LessOrEqual(t, int(r), 126, "ASCII glyph above printable range")
	}
}

func TestGetGlyphSet(t *testing.T) {
	tests := []struct {
		name           string
		mode           TerminalMode
		expectASCII    bool
		expectBasic    bool
		expectExtended bool
	}{
		{
			name:           "TerminalDumb",
			mode:           TerminalDumb,
			expectASCII:    true,
			expectBasic:    false,
			expectExtended: false,
		},
		{
			name:           "TerminalBasic",
			mode:           TerminalBasic,
			expectASCII:    true,
			expectBasic:    true,
			expectExtended: false,
		},
		{
			name:           "TerminalFull",
			mode:           TerminalFull,
			expectASCII:    true,
			expectBasic:    true,
			expectExtended: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := GetGlyphSet(tt.mode)
			require.NotNil(t, gs)
			assert.Equal(t, tt.mode, gs.Mode)
			assert.NotEmpty(t, gs.Glyphs)

			// Check glyph count expectations
			if tt.expectASCII && !tt.expectBasic {
				assert.Equal(t, len(MatrixGlyphs.ASCII), gs.Count())
			}
			if tt.expectBasic && !tt.expectExtended {
				assert.Equal(t, len(MatrixGlyphs.ASCII)+len(MatrixGlyphs.Basic), gs.Count())
			}
			if tt.expectExtended {
				expected := len(MatrixGlyphs.ASCII) + len(MatrixGlyphs.Basic) + len(MatrixGlyphs.Extended)
				assert.Equal(t, expected, gs.Count())
			}
		})
	}
}

func TestGlyphSetSelectRandomGlyph(t *testing.T) {
	gs := GetGlyphSet(TerminalFull)
	require.NotNil(t, gs)

	// Generate multiple random glyphs
	seen := make(map[rune]bool)
	for range 100 {
		glyph, err := gs.SelectRandomGlyph()
		assert.NoError(t, err)
		assert.NotEqual(t, rune(0), glyph)

		// Should be in glyph set
		assert.True(t, gs.IsValidGlyph(glyph))

		seen[glyph] = true
	}

	// Should have some variety (not all the same)
	assert.Greater(t, len(seen), 10, "should generate varied glyphs")
}

func TestGlyphSetCount(t *testing.T) {
	tests := []struct {
		name string
		mode TerminalMode
	}{
		{"Dumb", TerminalDumb},
		{"Basic", TerminalBasic},
		{"Full", TerminalFull},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := GetGlyphSet(tt.mode)
			count := gs.Count()
			assert.Greater(t, count, 0)
			assert.Equal(t, len(gs.Glyphs), count)
		})
	}
}

func TestGlyphSetIsValidGlyph(t *testing.T) {
	gs := GetGlyphSet(TerminalDumb) // ASCII only

	// Test valid glyphs
	for _, glyph := range gs.Glyphs {
		assert.True(t, gs.IsValidGlyph(glyph))
	}

	// Test invalid glyphs (not in set)
	invalidGlyphs := []rune{'a', 'A', '0', '9', 'α', 'Σ'}
	for _, glyph := range invalidGlyphs {
		assert.False(t, gs.IsValidGlyph(glyph), "glyph %c should not be valid", glyph)
	}
}

func TestIsSafe(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		{"ASCII letter", 'A', true},
		{"ASCII digit", '5', true},
		{"ASCII symbol", '!', true},
		{"Greek letter", 'α', true},
		{"Cyrillic letter", 'А', true},
		{"Control character", '\x00', false},
		{"Newline", '\n', false},
		{"Tab", '\t', false},
		{"RuneError", utf8.RuneError, false},
		{"Valid emoji", '😀', true},
		{"Zero-width joiner", '\u200D', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSafe(tt.r)
			assert.Equal(t, tt.want, got, "IsSafe(%U)", tt.r)
		})
	}
}

func TestValidateGlyphSet(t *testing.T) {
	tests := []struct {
		name    string
		glyphs  []rune
		wantErr bool
	}{
		{
			name:    "valid glyphs",
			glyphs:  []rune{'A', 'B', 'C', '!', '@', '#'},
			wantErr: false,
		},
		{
			name:    "empty set",
			glyphs:  []rune{},
			wantErr: true,
		},
		{
			name:    "with control characters",
			glyphs:  []rune{'A', '\n', 'B'},
			wantErr: true,
		},
		{
			name:    "with rune error",
			glyphs:  []rune{'A', utf8.RuneError, 'B'},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGlyphSet(tt.glyphs)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMatrixGlyphSetsPassValidation(t *testing.T) {
	// All predefined glyph sets should pass validation
	tests := []struct {
		name   string
		glyphs []rune
	}{
		{"ASCII", MatrixGlyphs.ASCII},
		{"Basic", MatrixGlyphs.Basic},
		{"Extended", MatrixGlyphs.Extended},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGlyphSet(tt.glyphs)
			assert.NoError(t, err, "predefined glyph set %s should be valid", tt.name)
		})
	}
}

func TestGlyphSetNoDuplicates(t *testing.T) {
	tests := []struct {
		name   string
		glyphs []rune
	}{
		{"ASCII", MatrixGlyphs.ASCII},
		{"Basic", MatrixGlyphs.Basic},
		{"Extended", MatrixGlyphs.Extended},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seen := make(map[rune]bool)
			for _, r := range tt.glyphs {
				assert.False(t, seen[r], "duplicate glyph %U in %s", r, tt.name)
				seen[r] = true
			}
		})
	}
}

func BenchmarkSelectRandomGlyph(b *testing.B) {
	gs := GetGlyphSet(TerminalFull)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = gs.SelectRandomGlyph()
	}
}

func BenchmarkIsValidGlyph(b *testing.B) {
	gs := GetGlyphSet(TerminalFull)
	testGlyph := gs.Glyphs[0]
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = gs.IsValidGlyph(testGlyph)
	}
}

func BenchmarkIsSafe(b *testing.B) {
	testRune := 'A'
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = IsSafe(testRune)
	}
}
