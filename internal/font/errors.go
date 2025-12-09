package font

import (
	"errors"
	"fmt"
)

var (
	// ErrNoGlyphsAvailable indicates the glyph set is empty
	ErrNoGlyphsAvailable = errors.New("no glyphs available")

	// ErrInvalidGlyphIndex indicates an out-of-bounds glyph index
	ErrInvalidGlyphIndex = errors.New("invalid glyph index")
)

// UnsafeGlyphError indicates glyphs that are unsafe for terminal display
type UnsafeGlyphError struct {
	Glyphs []rune
}

func (e *UnsafeGlyphError) Error() string {
	return fmt.Sprintf("unsafe glyphs detected: %v (count: %d)", e.Glyphs, len(e.Glyphs))
}
