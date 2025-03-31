package themes

import "github.com/charmbracelet/lipgloss"

// DraculaTheme implements the Theme interface with Dracula colors
type DraculaTheme struct {
	colors ColorPalette
}

// NewDraculaTheme creates a new instance of the Dracula theme
func NewDraculaTheme() Theme {
	return &DraculaTheme{
		colors: ColorPalette{
			// Dracula colors
			Primary:   lipgloss.Color("#bd93f9"), // purple
			Secondary: lipgloss.Color("#ff79c6"), // pink
			Tertiary:  lipgloss.Color("#f1fa8c"), // yellow

			Background: lipgloss.Color("#282a36"), // background
			Foreground: lipgloss.Color("#f8f8f2"), // foreground
			Border:     lipgloss.Color("#bd93f9"), // purple

			Success: lipgloss.Color("#50fa7b"), // green
			Error:   lipgloss.Color("#ff5555"), // red
			Warning: lipgloss.Color("#ffb86c"), // orange
			Info:    lipgloss.Color("#8be9fd"), // cyan

			Directory:  lipgloss.Color("#8be9fd"), // cyan
			File:       lipgloss.Color("#f8f8f2"), // foreground
			Selected:   lipgloss.Color("#50fa7b"), // green
			Deselected: lipgloss.Color("#6272a4"), // comment

			Text:      lipgloss.Color("#f8f8f2"), // foreground
			Muted:     lipgloss.Color("#6272a4"), // comment
			Highlight: lipgloss.Color("#ff79c6"), // pink
		},
	}
}

// Colors returns the color palette
func (t *DraculaTheme) Colors() ColorPalette {
	return t.colors
}

// Name returns the theme name
func (t *DraculaTheme) Name() string {
	return "dracula"
}
