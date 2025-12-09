package font

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectTerminalMode(t *testing.T) {
	tests := []struct {
		name      string
		term      string
		colorTerm string
		want      TerminalMode
	}{
		{
			name: "xterm-256color",
			term: "xterm-256color",
			want: TerminalFull,
		},
		{
			name: "alacritty",
			term: "alacritty",
			want: TerminalFull,
		},
		{
			name: "kitty",
			term: "kitty",
			want: TerminalFull,
		},
		{
			name:      "truecolor support",
			term:      "xterm",
			colorTerm: "truecolor",
			want:      TerminalFull,
		},
		{
			name: "basic xterm",
			term: "xterm",
			want: TerminalBasic,
		},
		{
			name: "linux console",
			term: "linux",
			want: TerminalBasic,
		},
		{
			name: "dumb terminal",
			term: "dumb",
			want: TerminalDumb,
		},
		{
			name: "empty terminal",
			term: "",
			want: TerminalDumb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env
			origTerm := os.Getenv("TERM")
			origColorTerm := os.Getenv("COLORTERM")
			defer func() {
				_ = os.Setenv("TERM", origTerm)
				_ = os.Setenv("COLORTERM", origColorTerm)
			}()

			// Set test env
			_ = os.Setenv("TERM", tt.term)
			_ = os.Setenv("COLORTERM", tt.colorTerm)

			got := DetectTerminalMode()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsModernTerminal(t *testing.T) {
	tests := []struct {
		name      string
		term      string
		colorTerm string
		want      bool
	}{
		{"xterm-256color", "xterm-256color", "", true},
		{"screen-256color", "screen-256color", "", true},
		{"alacritty", "alacritty", "", true},
		{"kitty", "kitty", "", true},
		{"iterm2", "iterm2", "", true},
		{"truecolor", "xterm", "truecolor", true},
		{"24bit", "xterm", "24bit", true},
		{"basic xterm", "xterm", "", false},
		{"vt100", "vt100", "", false},
		{"dumb", "dumb", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isModernTerminal(tt.term, tt.colorTerm)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsBasicTerminal(t *testing.T) {
	tests := []struct {
		name string
		term string
		want bool
	}{
		{"xterm", "xterm", true},
		{"vt220", "vt220", true},
		{"linux", "linux", true},
		{"screen", "screen", true},
		{"tmux", "tmux", true},
		{"rxvt", "rxvt", true},
		{"dumb", "dumb", false},
		{"unknown", "unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isBasicTerminal(tt.term)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTerminalModeString(t *testing.T) {
	tests := []struct {
		mode TerminalMode
		want string
	}{
		{TerminalDumb, "ASCII-only (dumb terminal)"},
		{TerminalBasic, "Basic UTF-8 (standard terminal)"},
		{TerminalFull, "Full Unicode (modern terminal)"},
		{TerminalMode(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.mode.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSupportsColor(t *testing.T) {
	tests := []struct {
		name    string
		term    string
		noColor string
		want    bool
	}{
		{"xterm-256color", "xterm-256color", "", true},
		{"xterm", "xterm", "", true},
		{"linux", "linux", "", true},
		{"screen", "screen", "", true},
		{"dumb", "dumb", "", false},
		{"with NO_COLOR", "xterm-256color", "1", false},
		{"empty", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env
			origTerm := os.Getenv("TERM")
			origNoColor := os.Getenv("NO_COLOR")
			defer func() {
				_ = os.Setenv("TERM", origTerm)
				_ = os.Setenv("NO_COLOR", origNoColor)
			}()

			// Set test env
			_ = os.Setenv("TERM", tt.term)
			if tt.noColor != "" {
				_ = os.Setenv("NO_COLOR", tt.noColor)
			} else {
				_ = os.Unsetenv("NO_COLOR")
			}

			got := SupportsColor()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSupportsTrueColor(t *testing.T) {
	tests := []struct {
		name      string
		colorTerm string
		want      bool
	}{
		{"truecolor", "truecolor", true},
		{"24bit", "24bit", true},
		{"TrueColor", "TrueColor", true},
		{"24BIT", "24BIT", true},
		{"empty", "", false},
		{"other", "other", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env
			origColorTerm := os.Getenv("COLORTERM")
			defer func() { _ = os.Setenv("COLORTERM", origColorTerm) }()

			// Set test env
			_ = os.Setenv("COLORTERM", tt.colorTerm)

			got := SupportsTrueColor()
			assert.Equal(t, tt.want, got)
		})
	}
}

func BenchmarkDetectTerminalMode(b *testing.B) {
	_ = os.Setenv("TERM", "xterm-256color")
	_ = os.Setenv("COLORTERM", "truecolor")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = DetectTerminalMode()
	}
}
