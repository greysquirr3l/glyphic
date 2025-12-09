package font

import (
	"os"
	"strings"
)

// DetectTerminalMode detects the terminal's capability level
// from environment variables (TERM, COLORTERM)
func DetectTerminalMode() TerminalMode {
	term := os.Getenv("TERM")
	colorTerm := os.Getenv("COLORTERM")

	// Check for modern terminals with full Unicode support
	if isModernTerminal(term, colorTerm) {
		return TerminalFull
	}

	// Check for basic terminals with UTF-8 support
	if isBasicTerminal(term) {
		return TerminalBasic
	}

	// Fallback to ASCII-only
	return TerminalDumb
}

// isModernTerminal checks if terminal supports full Unicode
func isModernTerminal(term, colorTerm string) bool {
	// Alacritty, Kitty, iTerm2, modern terminals
	modernTerms := []string{
		"xterm-256color",
		"screen-256color",
		"tmux-256color",
		"alacritty",
		"kitty",
		"iterm",
		"iterm2",
		"vte",
		"gnome",
		"konsole",
	}

	term = strings.ToLower(term)
	for _, modern := range modernTerms {
		if strings.Contains(term, modern) {
			return true
		}
	}

	// Check COLORTERM for truecolor support
	if colorTerm != "" {
		ct := strings.ToLower(colorTerm)
		if strings.Contains(ct, "truecolor") || strings.Contains(ct, "24bit") {
			return true
		}
	}

	return false
}

// isBasicTerminal checks if terminal supports basic UTF-8
func isBasicTerminal(term string) bool {
	// Basic UTF-8 capable terminals
	basicTerms := []string{
		"xterm",
		"vt220",
		"linux",
		"screen",
		"tmux",
		"rxvt",
	}

	term = strings.ToLower(term)
	for _, basic := range basicTerms {
		if strings.Contains(term, basic) {
			return true
		}
	}

	return false
}

// GetCapabilityString returns a human-readable description
func (mode TerminalMode) String() string {
	switch mode {
	case TerminalDumb:
		return "ASCII-only (dumb terminal)"
	case TerminalBasic:
		return "Basic UTF-8 (standard terminal)"
	case TerminalFull:
		return "Full Unicode (modern terminal)"
	default:
		return "Unknown"
	}
}

// SupportsColor checks if terminal likely supports ANSI colors
func SupportsColor() bool {
	term := os.Getenv("TERM")

	// Explicitly no color
	if term == "dumb" || os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Check for color terminal indicators
	if strings.Contains(term, "color") {
		return true
	}

	// Common color-capable terminals
	colorTerms := []string{
		"xterm",
		"screen",
		"tmux",
		"linux",
		"vt100",
		"ansi",
	}

	term = strings.ToLower(term)
	for _, ct := range colorTerms {
		if strings.Contains(term, ct) {
			return true
		}
	}

	return false
}

// SupportsTrueColor checks if terminal supports 24-bit color
func SupportsTrueColor() bool {
	colorTerm := strings.ToLower(os.Getenv("COLORTERM"))
	return strings.Contains(colorTerm, "truecolor") || strings.Contains(colorTerm, "24bit")
}
