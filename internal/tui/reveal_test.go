package tui

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/greysquirr3l/glyphic/internal/font"
	"github.com/stretchr/testify/assert"
)

func TestNewRevealModel(t *testing.T) {
	tests := []struct {
		name     string
		password string
		opts     RevealOptions
	}{
		{
			name:     "default options",
			password: "test-password",
			opts: RevealOptions{
				Scheme:       MatrixScheme,
				Speed:        SpeedNormal,
				TerminalMode: font.TerminalFull,
			},
		},
		{
			name:     "cyber scheme fast speed",
			password: "another-test",
			opts: RevealOptions{
				Scheme:       CyberScheme,
				Speed:        SpeedFast,
				TerminalMode: font.TerminalFull,
			},
		},
		{
			name:     "mono scheme slow speed",
			password: "slow-reveal",
			opts: RevealOptions{
				Scheme:       MonoScheme,
				Speed:        SpeedSlow,
				TerminalMode: font.TerminalBasic,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewRevealModel(tt.password, tt.opts)
			assert.NotNil(t, model)
			assert.Equal(t, tt.password, model.password)
			assert.Equal(t, tt.opts.Scheme.Name, model.opts.Scheme.Name)
			assert.Equal(t, len(tt.password), len(model.chars))
			assert.Equal(t, 0, model.currentStep)
			assert.False(t, model.done)
			assert.NotNil(t, model.glyphSet)

			// Check all characters start unrevealed
			for i, char := range model.chars {
				assert.False(t, char.Revealed, "char %d should start unrevealed", i)
				assert.Equal(t, rune(tt.password[i]), char.Target)
			}
		})
	}
}

func TestRevealModelInit(t *testing.T) {
	model := NewRevealModel("test", RevealOptions{
		Scheme:       MatrixScheme,
		Speed:        SpeedNormal,
		TerminalMode: font.TerminalFull,
	})

	cmd := model.Init()
	assert.NotNil(t, cmd)
}

func TestRevealModelUpdate(t *testing.T) {
	t.Run("tick message progresses animation", func(t *testing.T) {
		model := NewRevealModel("test", RevealOptions{
			Scheme:       MatrixScheme,
			Speed:        SpeedNormal,
			TerminalMode: font.TerminalFull,
		})

		initialStep := model.currentStep

		// Tick should advance step
		newModel, cmd := model.Update(tickMsg{})
		assert.NotNil(t, cmd)
		m := newModel.(RevealModel)
		assert.Greater(t, m.currentStep, initialStep)
	})

	t.Run("completes after totalSteps", func(t *testing.T) {
		model := NewRevealModel("ab", RevealOptions{
			Scheme:       MatrixScheme,
			Speed:        SpeedNormal,
			TerminalMode: font.TerminalFull,
		})

		// Fast-forward to completion
		model.currentStep = model.totalSteps

		newModel, cmd := model.Update(tickMsg{})
		m := newModel.(RevealModel)
		assert.True(t, m.done)
		// Should not auto-quit, password stays visible
		assert.Nil(t, cmd)
	})

	t.Run("quit key exits", func(t *testing.T) {
		model := NewRevealModel("test", RevealOptions{
			Scheme:       MatrixScheme,
			Speed:        SpeedNormal,
			TerminalMode: font.TerminalFull,
		})

		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
		newModel, cmd := model.Update(msg)
		m := newModel.(RevealModel)
		assert.True(t, m.done)
		assert.NotNil(t, cmd)
	})

	t.Run("ctrl-c exits", func(t *testing.T) {
		model := NewRevealModel("test", RevealOptions{
			Scheme:       MatrixScheme,
			Speed:        SpeedNormal,
			TerminalMode: font.TerminalFull,
		})

		msg := tea.KeyMsg{Type: tea.KeyCtrlC}
		newModel, cmd := model.Update(msg)
		m := newModel.(RevealModel)
		assert.True(t, m.done)
		assert.NotNil(t, cmd)
	})

	t.Run("enter completes reveal immediately", func(t *testing.T) {
		model := NewRevealModel("test", RevealOptions{
			Scheme:       MatrixScheme,
			Speed:        SpeedNormal,
			TerminalMode: font.TerminalFull,
		})

		msg := tea.KeyMsg{Type: tea.KeyEnter}
		newModel, _ := model.Update(msg)
		m := newModel.(RevealModel)

		// All characters should be revealed
		assert.True(t, m.done)
		for i, char := range m.chars {
			assert.True(t, char.Revealed, "char %d should be revealed", i)
			assert.Equal(t, char.Target, char.Current, "char %d should show target", i)
		}
	})

	t.Run("space completes reveal immediately", func(t *testing.T) {
		model := NewRevealModel("test", RevealOptions{
			Scheme:       MatrixScheme,
			Speed:        SpeedNormal,
			TerminalMode: font.TerminalFull,
		})

		msg := tea.KeyMsg{Type: tea.KeySpace}
		newModel, _ := model.Update(msg)
		m := newModel.(RevealModel)

		assert.True(t, m.done)
	})

	t.Run("window size updates dimensions", func(t *testing.T) {
		model := NewRevealModel("test", RevealOptions{
			Scheme:       MatrixScheme,
			Speed:        SpeedNormal,
			TerminalMode: font.TerminalFull,
		})

		msg := tea.WindowSizeMsg{Width: 120, Height: 40}
		newModel, _ := model.Update(msg)
		m := newModel.(RevealModel)

		assert.Equal(t, 120, m.width)
		assert.Equal(t, 40, m.height)
	})
}

func TestRevealModelView(t *testing.T) {
	t.Run("renders view", func(t *testing.T) {
		model := NewRevealModel("test", RevealOptions{
			Scheme:       MatrixScheme,
			Speed:        SpeedNormal,
			TerminalMode: font.TerminalFull,
		})

		view := model.View()
		assert.NotEmpty(t, view)
		// Should contain help text
		assert.Contains(t, view, "skip")
	})

	t.Run("fully revealed shows password", func(t *testing.T) {
		model := NewRevealModel("test", RevealOptions{
			Scheme:       MatrixScheme,
			Speed:        SpeedNormal,
			TerminalMode: font.TerminalFull,
		})

		// Reveal all characters
		for i := range model.chars {
			model.chars[i].Revealed = true
			model.chars[i].Current = model.chars[i].Target
		}
		model.done = true

		view := model.View()
		assert.NotEmpty(t, view)
		// Should contain the password characters
		assert.Contains(t, view, "t")
		assert.Contains(t, view, "e")
		assert.Contains(t, view, "s")
	})

	t.Run("all color schemes render", func(t *testing.T) {
		schemes := []ColorScheme{
			MatrixScheme,
			CyberScheme,
			FireScheme,
			VaporScheme,
			MonoScheme,
			NordScheme,
			GruvboxScheme,
		}

		for _, scheme := range schemes {
			t.Run(scheme.Name, func(t *testing.T) {
				model := NewRevealModel("test", RevealOptions{
					Scheme:       scheme,
					Speed:        SpeedNormal,
					TerminalMode: font.TerminalFull,
				})

				view := model.View()
				assert.NotEmpty(t, view)
			})
		}
	})

	t.Run("shows entropy when done and requested", func(t *testing.T) {
		model := NewRevealModel("test", RevealOptions{
			Scheme:       MatrixScheme,
			Speed:        SpeedNormal,
			TerminalMode: font.TerminalFull,
			ShowEntropy:  true,
		})

		model.done = true

		view := model.View()
		assert.Contains(t, view, "Length")
	})
}

func TestSpeedConstants(t *testing.T) {
	tests := []struct {
		name     string
		speed    Speed
		expected time.Duration
	}{
		{"slow", SpeedSlow, 150 * time.Millisecond},
		{"normal", SpeedNormal, 75 * time.Millisecond},
		{"fast", SpeedFast, 30 * time.Millisecond},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, time.Duration(tt.speed)*time.Millisecond)
		})
	}
}

func TestRevealAnimation(t *testing.T) {
	t.Run("characters update during animation", func(t *testing.T) {
		password := "abc"
		model := NewRevealModel(password, RevealOptions{
			Scheme:       MatrixScheme,
			Speed:        SpeedNormal,
			TerminalMode: font.TerminalFull,
		})

		// Simulate several animation steps
		for range 50 {
			newModel, _ := model.Update(tickMsg{})
			model = newModel.(RevealModel)

			if model.done {
				break
			}
		}

		// At least some progress should have been made
		assert.Greater(t, model.currentStep, 0)
	})

	t.Run("handles empty password", func(t *testing.T) {
		model := NewRevealModel("", RevealOptions{
			Scheme:       MatrixScheme,
			Speed:        SpeedNormal,
			TerminalMode: font.TerminalFull,
		})

		view := model.View()
		assert.NotEmpty(t, view) // Should still render help text
	})

	t.Run("handles long password", func(t *testing.T) {
		longPassword := strings.Repeat("word-", 20) // 100 characters
		model := NewRevealModel(longPassword, RevealOptions{
			Scheme:       MatrixScheme,
			Speed:        SpeedFast,
			TerminalMode: font.TerminalFull,
		})

		// Should initialize without panic
		assert.NotNil(t, model)
		assert.Equal(t, len(longPassword), len(model.chars))

		// Test a few ticks
		for range 10 {
			newModel, _ := model.Update(tickMsg{})
			model = newModel.(RevealModel)
		}

		assert.Greater(t, model.currentStep, 0)
	})

	t.Run("handles special characters", func(t *testing.T) {
		password := "test!@#$%^&*()_+-=[]{}|"
		model := NewRevealModel(password, RevealOptions{
			Scheme:       MatrixScheme,
			Speed:        SpeedNormal,
			TerminalMode: font.TerminalFull,
		})

		assert.NotNil(t, model)
		assert.Equal(t, len(password), len(model.chars))

		// Should render without panic
		view := model.View()
		assert.NotEmpty(t, view)
	})
}

func TestDefaultRevealOptions(t *testing.T) {
	opts := DefaultRevealOptions

	assert.Equal(t, MatrixScheme.Name, opts.Scheme.Name)
	assert.Equal(t, SpeedNormal, opts.Speed)
	assert.False(t, opts.ShowEntropy)
	// TerminalMode is detected, so it could be any valid value
	assert.GreaterOrEqual(t, int(opts.TerminalMode), int(font.TerminalDumb))
	assert.LessOrEqual(t, int(opts.TerminalMode), int(font.TerminalFull))
}

func TestTickCommand(t *testing.T) {
	cmd := tick(SpeedFast)
	assert.NotNil(t, cmd)
}

// Benchmarks
func BenchmarkRevealModelUpdate(b *testing.B) {
	model := NewRevealModel("test-password-123", RevealOptions{
		Scheme:       MatrixScheme,
		Speed:        SpeedNormal,
		TerminalMode: font.TerminalFull,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = model.Update(tickMsg{})
	}
}

func BenchmarkRevealModelView(b *testing.B) {
	model := NewRevealModel("test-password-123", RevealOptions{
		Scheme:       MatrixScheme,
		Speed:        SpeedNormal,
		TerminalMode: font.TerminalFull,
	})

	model.width = 80
	model.height = 24

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = model.View()
	}
}

func BenchmarkUpdateReveal(b *testing.B) {
	model := NewRevealModel("test-password-with-many-characters", RevealOptions{
		Scheme:       MatrixScheme,
		Speed:        SpeedNormal,
		TerminalMode: font.TerminalFull,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model.updateReveal()
	}
}
