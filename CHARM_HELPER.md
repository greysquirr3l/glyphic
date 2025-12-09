# Charmbracelet Package Reference for Glyphic

A quick reference guide for the Charm ecosystem packages used in the glyphic project.

---

## Package Overview

| Package | Purpose | Import |
|---------|---------|--------|
| Lip Gloss | Terminal styling (colors, borders, layout) | `github.com/charmbracelet/lipgloss` |
| Bubble Tea | TUI framework (Elm architecture) | `github.com/charmbracelet/bubbletea` |
| Bubbles | Pre-built TUI components | `github.com/charmbracelet/bubbles` |
| Harmonica | Spring-based animations | `github.com/charmbracelet/harmonica` |
| Log | Styled logging | `github.com/charmbracelet/log` |
| Glamour | Markdown rendering | `github.com/charmbracelet/glamour` |
| Huh | Form/prompt library | `github.com/charmbracelet/huh` |

---

## Lip Gloss — Styling

### Basic Style Creation

```go
import "github.com/charmbracelet/lipgloss"

style := lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Background(lipgloss.Color("#7D56F4")).
    Padding(1, 2).
    Margin(1, 0).
    Width(40)

output := style.Render("Hello, World")
```

### Color Types

```go
// ANSI 16 colors (0-15)
lipgloss.Color("5")      // Magenta
lipgloss.Color("9")      // Bright Red

// ANSI 256 colors (0-255)
lipgloss.Color("86")     // Aqua
lipgloss.Color("201")    // Hot Pink

// True Color (hex)
lipgloss.Color("#FF5733")
lipgloss.Color("#04B575")

// Adaptive (light/dark backgrounds)
lipgloss.AdaptiveColor{
    Light: "#333333",
    Dark:  "#EEEEEE",
}

// Complete (specify all profiles)
lipgloss.CompleteColor{
    TrueColor: "#FF5733",
    ANSI256:   "202",
    ANSI:      "9",
}
```

### Text Formatting

```go
style := lipgloss.NewStyle().
    Bold(true).
    Italic(true).
    Underline(true).
    Strikethrough(true).
    Faint(true).
    Blink(true).
    Reverse(true)
```

### Layout & Spacing

```go
style := lipgloss.NewStyle().
    Width(40).
    Height(10).
    Padding(1, 2, 1, 2).     // top, right, bottom, left
    PaddingTop(1).
    PaddingRight(2).
    PaddingBottom(1).
    PaddingLeft(2).
    Margin(1, 2).
    MarginTop(1).
    Align(lipgloss.Center).  // Left, Center, Right
    AlignVertical(lipgloss.Center)
```

### Borders

```go
style := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("63")).
    BorderBackground(lipgloss.Color("0")).
    BorderTop(true).
    BorderRight(true).
    BorderBottom(true).
    BorderLeft(true)

// Border types
lipgloss.NormalBorder()   // ┌─┐│└─┘
lipgloss.RoundedBorder()  // ╭─╮│╰─╯
lipgloss.ThickBorder()    // ┏━┓┃┗━┛
lipgloss.DoubleBorder()   // ╔═╗║╚═╝
lipgloss.HiddenBorder()   // spaces (preserves layout)
lipgloss.BlockBorder()    // █▀█ █▄█
```

### Compositing & Layout

```go
// Join strings horizontally
row := lipgloss.JoinHorizontal(lipgloss.Top, left, center, right)
row := lipgloss.JoinHorizontal(lipgloss.Center, a, b, c)
row := lipgloss.JoinHorizontal(lipgloss.Bottom, x, y, z)

// Join strings vertically
col := lipgloss.JoinVertical(lipgloss.Left, top, middle, bottom)
col := lipgloss.JoinVertical(lipgloss.Center, a, b, c)
col := lipgloss.JoinVertical(lipgloss.Right, x, y, z)

// Place content in a box
box := lipgloss.Place(
    width, height,
    lipgloss.Center, lipgloss.Center,  // horizontal, vertical position
    content,
    lipgloss.WithWhitespaceChars("·"),
    lipgloss.WithWhitespaceForeground(lipgloss.Color("240")),
)

// Measure strings (ANSI-aware)
w := lipgloss.Width(styledString)
h := lipgloss.Height(styledString)
```

### Style Inheritance & Copying

```go
baseStyle := lipgloss.NewStyle().
    Foreground(lipgloss.Color("white")).
    Background(lipgloss.Color("black"))

// Copy and modify
derivedStyle := baseStyle.
    Bold(true).
    Foreground(lipgloss.Color("green"))

// Unset properties
resetStyle := style.
    UnsetBold().
    UnsetForeground()
```

---

## Lip Gloss Sub-packages

### Tables (`lipgloss/table`)

```go
import "github.com/charmbracelet/lipgloss/table"

t := table.New().
    Border(lipgloss.NormalBorder()).
    BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
    Headers("NAME", "SIZE", "STATUS").
    Row("wordlist-1", "7,776", "✓").
    Row("wordlist-2", "2,048", "✓").
    Row("wordlist-3", "1,296", "✓").
    StyleFunc(func(row, col int) lipgloss.Style {
        if row == table.HeaderRow {
            return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
        }
        if row%2 == 0 {
            return lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
        }
        return lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
    })

fmt.Println(t)
```

### Lists (`lipgloss/list`)

```go
import "github.com/charmbracelet/lipgloss/list"

l := list.New(
    "EFF Large Wordlist",
    "BIP-39 English",
    "SecureDrop Wordlist",
).
    Enumerator(list.Bullet).        // Bullet, Arabic, Roman, Alphabet, Tree
    EnumeratorStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
    ItemStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("212")))

fmt.Println(l)
```

### Trees (`lipgloss/tree`)

```go
import "github.com/charmbracelet/lipgloss/tree"

t := tree.New().
    Root("glyphic").
    Child(
        tree.New().Root("wordlists").
            Child("eff-large.txt").
            Child("bip39.txt"),
    ).
    Child(
        tree.New().Root("exclusions").
            Child("profanity.txt").
            Child("confusing.txt"),
    )

fmt.Println(t)
```

---

## Bubble Tea — TUI Framework

### Basic Structure

```go
import tea "github.com/charmbracelet/bubbletea"

type model struct {
    // Your state
}

func (m model) Init() tea.Cmd {
    // Return initial command (or nil)
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        }
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    }
    return m, nil
}

func (m model) View() string {
    return "Hello, Bubble Tea!"
}

func main() {
    p := tea.NewProgram(model{})
    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### Common Messages

```go
// Key press
case tea.KeyMsg:
    switch msg.String() {
    case "enter":
    case "esc":
    case "up", "k":
    case "down", "j":
    case "ctrl+c":
    }
    // Or check specific keys
    switch msg.Type {
    case tea.KeyEnter:
    case tea.KeyEsc:
    case tea.KeyCtrlC:
    }

// Window resize
case tea.WindowSizeMsg:
    width := msg.Width
    height := msg.Height

// Mouse (if enabled)
case tea.MouseMsg:
    x, y := msg.X, msg.Y
    button := msg.Button
```

### Commands

```go
// Quit
return m, tea.Quit

// Batch multiple commands
return m, tea.Batch(cmd1, cmd2, cmd3)

// Sequence commands
return m, tea.Sequence(cmd1, cmd2, cmd3)

// Tick (for animations)
func tick() tea.Cmd {
    return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

// Custom command
func doSomething() tea.Cmd {
    return func() tea.Msg {
        // Do work...
        return resultMsg{data: result}
    }
}
```

### Program Options

```go
p := tea.NewProgram(
    model{},
    tea.WithAltScreen(),           // Full-screen mode
    tea.WithMouseCellMotion(),     // Enable mouse
    tea.WithoutSignalHandler(),    // Custom signal handling
    tea.WithInput(customReader),   // Custom input
    tea.WithOutput(customWriter),  // Custom output
)
```

---

## Bubbles — Components

### Spinner

```go
import "github.com/charmbracelet/bubbles/spinner"

s := spinner.New()
s.Spinner = spinner.Dot  // Dot, Line, MiniDot, Jump, Pulse, Points, Globe, Moon, Monkey, Meter, Hamburger
s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

// In Update:
s, cmd = s.Update(msg)

// In View:
s.View()
```

### Progress Bar

```go
import "github.com/charmbracelet/bubbles/progress"

p := progress.New(
    progress.WithDefaultGradient(),     // Blue to pink gradient
    progress.WithWidth(40),
    progress.WithoutPercentage(),
)

// Or solid color
p := progress.New(
    progress.WithSolidFill("#FF5733"),
)

// Update
p.SetPercent(0.5)  // 50%

// View
p.View()
```

### Text Input

```go
import "github.com/charmbracelet/bubbles/textinput"

ti := textinput.New()
ti.Placeholder = "Enter password..."
ti.CharLimit = 64
ti.Width = 30
ti.EchoMode = textinput.EchoPassword  // EchoNormal, EchoPassword, EchoNone
ti.EchoCharacter = '•'
ti.Focus()

// Styling
ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("white"))

// In Update:
ti, cmd = ti.Update(msg)

// Get value:
value := ti.Value()
```

### Viewport (Scrollable Content)

```go
import "github.com/charmbracelet/bubbles/viewport"

vp := viewport.New(80, 20)
vp.SetContent(longContent)
vp.Style = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

// In Update:
vp, cmd = vp.Update(msg)

// View:
vp.View()

// Scroll programmatically:
vp.GotoTop()
vp.GotoBottom()
vp.LineDown(1)
vp.LineUp(1)
```

### Help

```go
import "github.com/charmbracelet/bubbles/help"
import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
    Up    key.Binding
    Down  key.Binding
    Quit  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
    return []key.Binding{k.Up, k.Down, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
    return [][]key.Binding{{k.Up, k.Down}, {k.Quit}}
}

var keys = keyMap{
    Up:   key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
    Down: key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
    Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
}

h := help.New()
h.View(keys)
```

---

## Harmonica — Animation

```go
import "github.com/charmbracelet/harmonica"

// Create a spring (for smooth animations)
spring := harmonica.NewSpring(harmonica.FPS(60), 6.0, 1.0)
// Parameters: FPS, frequency (Hz), damping ratio

// In animation loop:
position, velocity = spring.Update(position, velocity, targetPosition)

// Presets
harmonica.FPS(60)        // 60 frames per second
harmonica.FPS(120)       // 120 fps for smoother animation
```

---

## Log — Styled Logging

```go
import "github.com/charmbracelet/log"

// Basic usage
log.Debug("Debug message")
log.Info("Info message", "key", "value")
log.Warn("Warning", "count", 42)
log.Error("Error occurred", "err", err)
log.Fatal("Fatal error")  // Exits

// Set level
log.SetLevel(log.DebugLevel)  // DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel

// Custom logger
logger := log.NewWithOptions(os.Stderr, log.Options{
    ReportCaller:    true,
    ReportTimestamp: true,
    TimeFormat:      time.Kitchen,
    Prefix:          "glyphic",
})

// Styled logger
logger.SetStyles(log.DefaultStyles())

// With fields
logger.With("component", "wordlist").Info("Loading...")
```

---

## Glyphic-Specific Styles

### Color Schemes for Decode Animation

```go
package tui

import "github.com/charmbracelet/lipgloss"

type ColorScheme struct {
    Name      string
    Scrambled lipgloss.Style
    Decoding  lipgloss.Style
    Revealed  lipgloss.Style
    Border    lipgloss.Style
}

var ColorSchemes = map[string]ColorScheme{
    "matrix": {
        Name:      "Matrix",
        Scrambled: lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
        Decoding:  lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true),
        Revealed:  lipgloss.NewStyle().Foreground(lipgloss.Color("46")),  // Bright green
        Border:    lipgloss.NewStyle().Foreground(lipgloss.Color("34")),
    },
    "cyber": {
        Name:      "Cyber",
        Scrambled: lipgloss.NewStyle().Foreground(lipgloss.Color("238")),
        Decoding:  lipgloss.NewStyle().Foreground(lipgloss.Color("87")).Bold(true),
        Revealed:  lipgloss.NewStyle().Foreground(lipgloss.Color("51")),  // Cyan
        Border:    lipgloss.NewStyle().Foreground(lipgloss.Color("39")),
    },
    "fire": {
        Name:      "Fire",
        Scrambled: lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
        Decoding:  lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true),
        Revealed:  lipgloss.NewStyle().Foreground(lipgloss.Color("208")), // Orange
        Border:    lipgloss.NewStyle().Foreground(lipgloss.Color("196")),
    },
    "vapor": {
        Name:      "Vaporwave",
        Scrambled: lipgloss.NewStyle().Foreground(lipgloss.Color("60")),
        Decoding:  lipgloss.NewStyle().Foreground(lipgloss.Color("213")).Bold(true),
        Revealed:  lipgloss.NewStyle().Foreground(lipgloss.Color("219")), // Pink
        Border:    lipgloss.NewStyle().Foreground(lipgloss.Color("141")),
    },
    "mono": {
        Name:      "Monochrome",
        Scrambled: lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
        Decoding:  lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true),
        Revealed:  lipgloss.NewStyle().Foreground(lipgloss.Color("255")),
        Border:    lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
    },
}
```

### Application Styles

```go
package tui

import "github.com/charmbracelet/lipgloss"

var (
    // Title/Header
    TitleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("205")).
        Background(lipgloss.Color("235")).
        Padding(0, 2).
        MarginBottom(1)
    
    // Status messages
    SuccessStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("42")).
        Bold(true)
    
    ErrorStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("196")).
        Bold(true)
    
    WarningStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("214"))
    
    // Info/subtle text
    SubtleStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("245"))
    
    // Entropy display
    EntropyHighStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("42")).
        Bold(true)
    
    EntropyMedStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("214"))
    
    EntropyLowStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("196"))
    
    // Password display box
    PasswordBoxStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("62")).
        Padding(1, 2).
        MarginTop(1).
        MarginBottom(1)
    
    // Progress/loading
    SpinnerStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("205"))
)
```

---

## Dependencies for Glyphic

```go
// go.mod
module github.com/youruser/glyphic

go 1.25

require (
    github.com/charmbracelet/bubbletea v1.2.0
    github.com/charmbracelet/lipgloss v1.0.0
    github.com/charmbracelet/bubbles v0.20.0
    github.com/charmbracelet/log v0.4.0
    github.com/charmbracelet/harmonica v0.2.0
    github.com/spf13/cobra v1.9.0
    golang.org/x/image v0.20.0
)
```

---

## Quick Tips

1. **Always use `lipgloss.Width()` for measuring** — it's ANSI-escape-aware
2. **Copy styles before modifying** — styles are value types, but explicit copies are clearer
3. **Use `AdaptiveColor` for light/dark support** — terminals vary
4. **Test with `TERM=dumb`** — ensures graceful degradation
5. **Use `tea.Batch()` for multiple commands** — keeps code clean
6. **Animations: 30-60 FPS is usually enough** — balance smoothness vs CPU
