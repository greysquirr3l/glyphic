package font

import (
	"unicode"
	"unicode/utf8"

	"github.com/greysquirr3l/glyphic/internal/security"
)

// TerminalMode represents terminal capability levels
type TerminalMode int

const (
	TerminalDumb  TerminalMode = iota // ASCII fallback (33-126)
	TerminalBasic                     // Basic UTF-8 (Latin + basic symbols)
	TerminalFull                      // Full Unicode (xterm-256color, kitty, alacritty)
)

// MatrixGlyphs contains curated glyphs from Matrix Code NFI font
// Organized by terminal capability level for graceful degradation
var MatrixGlyphs = struct {
	ASCII    []rune // Safe ASCII glyphs (33-126) for dumb terminals
	Basic    []rune // Basic UTF-8 glyphs for standard terminals
	Extended []rune // Extended Unicode glyphs for modern terminals
}{
	// ASCII: Printable ASCII excluding letters/numbers
	// Safe for vt100, dumb terminals
	ASCII: []rune{
		'!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
		':', ';', '<', '=', '>', '?', '@', '[', '\\', ']', '^', '_', '`', '{', '|', '}', '~',
	},

	// Basic: Latin-1 Supplement + common symbols
	// Safe for vt220, linux console
	Basic: []rune{
		// Latin-1 Supplement (U+00A0 - U+00FF)
		'¡', '¢', '£', '¤', '¥', '¦', '§', '¨', '©', 'ª', '«', '¬', '®', '¯',
		'°', '±', '²', '³', '´', 'µ', '¶', '·', '¸', '¹', 'º', '»', '¼', '½', '¾', '¿',
		'×', '÷',
		// Box Drawing (U+2500 - U+257F) - subset
		'─', '│', '┌', '┐', '└', '┘', '├', '┤', '┬', '┴', '┼',
		'═', '║', '╔', '╗', '╚', '╝', '╠', '╣', '╦', '╩', '╬',
		// Block Elements (U+2580 - U+259F)
		'▀', '▄', '█', '▌', '▐', '░', '▒', '▓',
	},

	// Extended: Full Unicode glyphs for modern terminals
	// Requires xterm-256color, kitty, alacritty, etc.
	Extended: []rune{
		// Greek and Coptic (U+0370 - U+03FF)
		'Α', 'Β', 'Γ', 'Δ', 'Ε', 'Ζ', 'Η', 'Θ', 'Ι', 'Κ', 'Λ', 'Μ',
		'Ν', 'Ξ', 'Ο', 'Π', 'Ρ', 'Σ', 'Τ', 'Υ', 'Φ', 'Χ', 'Ψ', 'Ω',
		'α', 'β', 'γ', 'δ', 'ε', 'ζ', 'η', 'θ', 'ι', 'κ', 'λ', 'μ',
		'ν', 'ξ', 'ο', 'π', 'ρ', 'ς', 'σ', 'τ', 'υ', 'φ', 'χ', 'ψ', 'ω',

		// Cyrillic (U+0400 - U+04FF) - subset
		'Ѐ', 'Ё', 'Ђ', 'Ѓ', 'Є', 'Ѕ', 'І', 'Ї', 'Ј', 'Љ', 'Њ', 'Ћ',
		'А', 'Б', 'В', 'Г', 'Д', 'Е', 'Ж', 'З', 'И', 'Й', 'К', 'Л',
		'М', 'Н', 'О', 'П', 'Р', 'С', 'Т', 'У', 'Ф', 'Х', 'Ц', 'Ч',
		'а', 'б', 'в', 'г', 'д', 'е', 'ж', 'з', 'и', 'й', 'к', 'л',
		'м', 'н', 'о', 'п', 'р', 'с', 'т', 'у', 'ф', 'х', 'ц', 'ч',

		// Hebrew (U+0590 - U+05FF) - subset
		'א', 'ב', 'ג', 'ד', 'ה', 'ו', 'ז', 'ח', 'ט', 'י', 'כ', 'ל',
		'מ', 'נ', 'ס', 'ע', 'פ', 'צ', 'ק', 'ר', 'ש', 'ת',

		// Arabic (U+0600 - U+06FF) - subset
		'ء', 'آ', 'أ', 'ؤ', 'إ', 'ئ', 'ا', 'ب', 'ة', 'ت', 'ث', 'ج',
		'ح', 'خ', 'د', 'ذ', 'ر', 'ز', 'س', 'ش', 'ص', 'ض', 'ط', 'ظ',

		// CJK Symbols (U+3000 - U+303F)
		'〒', '〓', '〔', '〕', '〖', '〗', '〘', '〙', '〚', '〛', '〜', '〝', '〞', '〟',

		// Hiragana (U+3040 - U+309F) - subset
		'あ', 'い', 'う', 'え', 'お', 'か', 'き', 'く', 'け', 'こ',
		'さ', 'し', 'す', 'せ', 'そ', 'た', 'ち', 'つ', 'て', 'と',

		// Katakana (U+30A0 - U+30FF) - subset
		'ア', 'イ', 'ウ', 'エ', 'オ', 'カ', 'キ', 'ク', 'ケ', 'コ',
		'サ', 'シ', 'ス', 'セ', 'ソ', 'タ', 'チ', 'ツ', 'テ', 'ト',

		// Mathematical Operators (U+2200 - U+22FF)
		'∀', '∂', '∃', '∅', '∇', '∈', '∉', '∋', '∏', '∑', '−', '∓', '∔', '∕',
		'∗', '∘', '√', '∝', '∞', '∟', '∠', '∡', '∢', '∣', '∤', '∥', '∦', '∧',
		'∨', '∩', '∪', '∫', '∬', '∭', '∮', '∯', '∰', '∱', '∲', '∳', '∴', '∵',

		// Geometric Shapes (U+25A0 - U+25FF)
		'■', '□', '▢', '▣', '▤', '▥', '▦', '▧', '▨', '▩', '▪', '▫', '▬', '▭',
		'▮', '▯', '▰', '▱', '▲', '△', '▴', '▵', '▶', '▷', '▸', '▹', '►', '▻',
		'▼', '▽', '▾', '▿', '◀', '◁', '◂', '◃', '◄', '◅', '◆', '◇', '◈', '◉',
		'◊', '○', '◌', '◍', '◎', '●', '◐', '◑', '◒', '◓', '◔', '◕', '◖', '◗',

		// Miscellaneous Symbols (U+2600 - U+26FF) - subset
		'☀', '☁', '☂', '☃', '☄', '★', '☆', '☇', '☈', '☉', '☊', '☋', '☌', '☍',
		'☎', '☏', '☐', '☑', '☒', '☓', '☔', '☕', '☖', '☗', '☘', '☙', '☚', '☛',

		// Dingbats (U+2700 - U+27BF) - subset
		'✁', '✂', '✃', '✄', '✅', '✆', '✇', '✈', '✉', '✊', '✋', '✌', '✍', '✎',
		'✏', '✐', '✑', '✒', '✓', '✔', '✕', '✖', '✗', '✘', '✙', '✚', '✛', '✜',

		// Braille Patterns (U+2800 - U+28FF) - subset
		'⠀', '⠁', '⠂', '⠃', '⠄', '⠅', '⠆', '⠇', '⠈', '⠉', '⠊', '⠋', '⠌', '⠍',
		'⠎', '⠏', '⠐', '⠑', '⠒', '⠓', '⠔', '⠕', '⠖', '⠗', '⠘', '⠙', '⠚', '⠛',
	},
}

// GlyphSet combines glyphs for a specific terminal mode
type GlyphSet struct {
	Mode   TerminalMode
	Glyphs []rune
}

// GetGlyphSet returns the appropriate glyph set for the terminal mode
func GetGlyphSet(mode TerminalMode) *GlyphSet {
	var glyphs []rune

	switch mode {
	case TerminalDumb:
		glyphs = MatrixGlyphs.ASCII
	case TerminalBasic:
		glyphs = append(append([]rune{}, MatrixGlyphs.ASCII...), MatrixGlyphs.Basic...)
	case TerminalFull:
		glyphs = append(append(append([]rune{}, MatrixGlyphs.ASCII...), MatrixGlyphs.Basic...), MatrixGlyphs.Extended...)
	default:
		glyphs = MatrixGlyphs.ASCII
	}

	return &GlyphSet{
		Mode:   mode,
		Glyphs: glyphs,
	}
}

// SelectRandomGlyph returns a cryptographically random glyph
func (g *GlyphSet) SelectRandomGlyph() (rune, error) {
	if len(g.Glyphs) == 0 {
		return 0, ErrNoGlyphsAvailable
	}

	idx, err := security.SecureRandomIndex(len(g.Glyphs))
	if err != nil {
		return 0, err
	}

	return g.Glyphs[idx], nil
}

// Count returns the number of glyphs available
func (g *GlyphSet) Count() int {
	return len(g.Glyphs)
}

// IsValidGlyph checks if a rune is in the glyph set
func (g *GlyphSet) IsValidGlyph(r rune) bool {
	for _, glyph := range g.Glyphs {
		if glyph == r {
			return true
		}
	}
	return false
}

// IsSafe checks if a rune is safe to display in the terminal
// Safe means: printable, not a control character, valid UTF-8
func IsSafe(r rune) bool {
	// Must be valid UTF-8
	if r == utf8.RuneError {
		return false
	}

	// Must be printable (not control character)
	if !unicode.IsPrint(r) {
		return false
	}

	// Must not be a combining character (can cause rendering issues)
	if unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Mc, r) || unicode.Is(unicode.Me, r) {
		return false
	}

	// Must not be a zero-width character
	if unicode.Is(unicode.Cf, r) {
		return false
	}

	return true
}

// ValidateGlyphSet ensures all glyphs in the set are safe
func ValidateGlyphSet(glyphs []rune) error {
	if len(glyphs) == 0 {
		return ErrNoGlyphsAvailable
	}

	unsafe := make([]rune, 0)
	for _, r := range glyphs {
		if !IsSafe(r) {
			unsafe = append(unsafe, r)
		}
	}

	if len(unsafe) > 0 {
		return &UnsafeGlyphError{Glyphs: unsafe}
	}

	return nil
}
