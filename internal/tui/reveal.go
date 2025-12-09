package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/greysquirr3l/glyphic/internal/font"
	"github.com/greysquirr3l/glyphic/internal/security"
)

// Speed defines animation speed presets
type Speed int

const (
	SpeedSlow   Speed = 150 // milliseconds per frame
	SpeedNormal Speed = 75
	SpeedFast   Speed = 30
)

// RevealOptions configures the reveal animation
type RevealOptions struct {
	Scheme       ColorScheme       // Color scheme to use
	Speed        Speed             // Animation speed
	TerminalMode font.TerminalMode // Terminal capability level
	ShowEntropy  bool              // Show entropy calculation
}

// DefaultRevealOptions provides sensible defaults
var DefaultRevealOptions = RevealOptions{
	Scheme:       MatrixScheme,
	Speed:        SpeedNormal,
	TerminalMode: font.DetectTerminalMode(),
	ShowEntropy:  false,
}

// CharState represents the reveal state of a character
type CharState struct {
	Target   rune // The final character to reveal
	Current  rune // The current displayed character
	Revealed bool // Whether fully revealed
}

// RevealModel is the Bubble Tea model for the reveal animation
type RevealModel struct {
	password    string
	chars       []CharState
	glyphSet    *font.GlyphSet
	opts        RevealOptions
	currentStep int
	totalSteps  int
	done        bool
	width       int
	height      int
}

// tickMsg is sent on each animation frame
type tickMsg time.Time

// NewRevealModel creates a new reveal animation model
func NewRevealModel(password string, opts RevealOptions) RevealModel {
	// Calculate total steps: ~10 frames per character
	totalSteps := len(password) * 10

	// Initialize character states
	chars := make([]CharState, len(password))
	for i, r := range password {
		chars[i] = CharState{
			Target:   r,
			Current:  r, // Start with target, will scramble
			Revealed: false,
		}
	}

	// Get appropriate glyph set for terminal
	glyphSet := font.GetGlyphSet(opts.TerminalMode)

	return RevealModel{
		password:    password,
		chars:       chars,
		glyphSet:    glyphSet,
		opts:        opts,
		currentStep: 0,
		totalSteps:  totalSteps,
		done:        false,
	}
}

// Init initializes the model
func (m RevealModel) Init() tea.Cmd {
	return tick(m.opts.Speed)
}

// Update handles messages
func (m RevealModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.done = true
			return m, tea.Quit
		case "enter", " ":
			// Skip to end
			for i := range m.chars {
				m.chars[i].Revealed = true
				m.chars[i].Current = m.chars[i].Target
			}
			m.done = true
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		if m.done {
			return m, nil
		}

		m.currentStep++

		// Update character states based on progress
		m.updateReveal()

		// Check if done
		if m.currentStep >= m.totalSteps {
			// Ensure all characters are fully revealed
			for i := range m.chars {
				m.chars[i].Revealed = true
				m.chars[i].Current = m.chars[i].Target
			}
			m.done = true
			// Don't auto-quit, let the password stay visible
			return m, nil
		}

		return m, tick(m.opts.Speed)
	}

	return m, nil
}

// updateReveal updates the reveal state for each character
func (m *RevealModel) updateReveal() {
	progress := float64(m.currentStep) / float64(m.totalSteps)

	for i := range m.chars {
		if m.chars[i].Revealed {
			continue
		}

		// Calculate this character's reveal threshold
		// Characters reveal in sequence with overlap
		charProgress := float64(i) / float64(len(m.chars))
		revealStart := charProgress * 0.5   // Start early
		revealEnd := charProgress*0.5 + 0.6 // Overlap

		if progress >= revealEnd {
			// Fully revealed
			m.chars[i].Revealed = true
			m.chars[i].Current = m.chars[i].Target
		} else if progress >= revealStart {
			// Scrambling/revealing phase
			// Randomly scramble with decreasing probability
			localProgress := (progress - revealStart) / (revealEnd - revealStart)

			if localProgress > 0.7 {
				// High chance of revealing
				shouldReveal, _ := security.SecureRandomIndex(10)
				if shouldReveal > 2 { // 70% chance
					m.chars[i].Current = m.chars[i].Target
				} else {
					glyph, err := m.glyphSet.SelectRandomGlyph()
					if err == nil {
						m.chars[i].Current = glyph
					}
				}
			} else {
				// Still scrambling
				glyph, err := m.glyphSet.SelectRandomGlyph()
				if err == nil {
					m.chars[i].Current = glyph
				}
			}
		} else {
			// Not started yet, show random glyph
			glyph, err := m.glyphSet.SelectRandomGlyph()
			if err == nil {
				m.chars[i].Current = glyph
			}
		}
	}
}

// View renders the current state
func (m RevealModel) View() string {
	if m.width == 0 {
		m.width = 80 // Default width
	}

	var b strings.Builder

	// Build the displayed password
	for i, char := range m.chars {
		var style lipgloss.Style

		if char.Revealed {
			// Revealed - use revealed color
			style = lipgloss.NewStyle().Foreground(m.opts.Scheme.Revealed)
		} else {
			// Scrambled or revealing
			progress := float64(m.currentStep) / float64(m.totalSteps)
			charProgress := float64(i) / float64(len(m.chars))

			if progress > charProgress {
				// Mid-reveal - use revealing color
				style = lipgloss.NewStyle().Foreground(m.opts.Scheme.Revealing)
			} else {
				// Scrambled - use scrambled color
				style = lipgloss.NewStyle().Foreground(m.opts.Scheme.Scrambled)
			}
		}

		b.WriteString(style.Render(string(char.Current)))
	}

	password := b.String()

	// Center the password
	padding := (m.width - len(m.password)) / 2
	if padding < 0 {
		padding = 0
	}

	result := strings.Repeat(" ", padding) + password

	// Add entropy info if requested
	if m.opts.ShowEntropy && m.done {
		result += "\n\n"
		entropyInfo := fmt.Sprintf("Length: %d characters", len(m.password))
		result += lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Render(entropyInfo)
	}

	// Add help text
	result += "\n\n"
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Italic(true)
	if !m.done {
		result += helpStyle.Render("Press SPACE to skip • ESC to quit")
	} else {
		result += helpStyle.Render("Press q or ESC to exit")
	}

	// Center vertically
	if m.height > 0 {
		verticalPadding := (m.height - 3) / 2
		if verticalPadding > 0 {
			result = strings.Repeat("\n", verticalPadding) + result
		}
	}

	return result
}

// tick creates a tick command with the specified delay
func tick(speed Speed) tea.Cmd {
	return tea.Tick(time.Duration(speed)*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Reveal runs the reveal animation and returns when complete
func Reveal(password string, opts RevealOptions) error {
	m := NewRevealModel(password, opts)
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run reveal animation: %w", err)
	}

	return nil
}
