package tui

import "github.com/charmbracelet/lipgloss"

// ColorScheme defines a color scheme for the reveal animation
type ColorScheme struct {
	Name       string
	Scrambled  lipgloss.Color // Color for scrambled glyphs
	Revealing  lipgloss.Color // Color for glyphs mid-reveal
	Revealed   lipgloss.Color // Color for fully revealed characters
	Background lipgloss.Color // Optional background color
}

// Predefined color schemes
var (
	// MatrixScheme - Classic green Matrix aesthetic
	MatrixScheme = ColorScheme{
		Name:       "matrix",
		Scrambled:  lipgloss.Color("#00FF00"), // Bright green
		Revealing:  lipgloss.Color("#00CC00"), // Medium green
		Revealed:   lipgloss.Color("#00FF00"), // Bright green
		Background: lipgloss.Color("#000000"), // Black
	}

	// CyberScheme - Cyberpunk neon blue
	CyberScheme = ColorScheme{
		Name:       "cyber",
		Scrambled:  lipgloss.Color("#00FFFF"), // Cyan
		Revealing:  lipgloss.Color("#0088FF"), // Blue
		Revealed:   lipgloss.Color("#00DDFF"), // Bright cyan
		Background: lipgloss.Color("#000000"), // Black
	}

	// FireScheme - Hot fire colors
	FireScheme = ColorScheme{
		Name:       "fire",
		Scrambled:  lipgloss.Color("#FF4400"), // Red-orange
		Revealing:  lipgloss.Color("#FFAA00"), // Orange-yellow
		Revealed:   lipgloss.Color("#FFFF00"), // Yellow
		Background: lipgloss.Color("#000000"), // Black
	}

	// VaporScheme - Vaporwave aesthetic
	VaporScheme = ColorScheme{
		Name:       "vapor",
		Scrambled:  lipgloss.Color("#FF6EC7"), // Hot pink
		Revealing:  lipgloss.Color("#C774E8"), // Purple
		Revealed:   lipgloss.Color("#00FFFF"), // Cyan
		Background: lipgloss.Color("#000000"), // Black
	}

	// MonoScheme - Monochrome for professional use
	MonoScheme = ColorScheme{
		Name:       "mono",
		Scrambled:  lipgloss.Color("#888888"), // Gray
		Revealing:  lipgloss.Color("#AAAAAA"), // Light gray
		Revealed:   lipgloss.Color("#FFFFFF"), // White
		Background: lipgloss.Color("#000000"), // Black
	}

	// NordScheme - Nord color palette (cool and professional)
	NordScheme = ColorScheme{
		Name:       "nord",
		Scrambled:  lipgloss.Color("#81A1C1"), // Nord frost 3
		Revealing:  lipgloss.Color("#88C0D0"), // Nord frost 2
		Revealed:   lipgloss.Color("#8FBCBB"), // Nord frost 1
		Background: lipgloss.Color("#2E3440"), // Nord polar night 0
	}

	// GruvboxScheme - Gruvbox color palette (warm and earthy)
	GruvboxScheme = ColorScheme{
		Name:       "gruvbox",
		Scrambled:  lipgloss.Color("#FE8019"), // Orange
		Revealing:  lipgloss.Color("#FABD2F"), // Yellow
		Revealed:   lipgloss.Color("#B8BB26"), // Green
		Background: lipgloss.Color("#282828"), // Dark background
	}
)

// AllSchemes returns all available color schemes
func AllSchemes() []ColorScheme {
	return []ColorScheme{
		MatrixScheme,
		CyberScheme,
		FireScheme,
		VaporScheme,
		MonoScheme,
		NordScheme,
		GruvboxScheme,
	}
}

// GetScheme returns a color scheme by name
func GetScheme(name string) ColorScheme {
	switch name {
	case "matrix":
		return MatrixScheme
	case "cyber":
		return CyberScheme
	case "fire":
		return FireScheme
	case "vapor":
		return VaporScheme
	case "mono":
		return MonoScheme
	case "nord":
		return NordScheme
	case "gruvbox":
		return GruvboxScheme
	default:
		return MatrixScheme
	}
}

// SchemeNames returns all available scheme names
func SchemeNames() []string {
	schemes := AllSchemes()
	names := make([]string, len(schemes))
	for i, s := range schemes {
		names[i] = s.Name
	}
	return names
}
