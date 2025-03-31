package themes

import "github.com/charmbracelet/lipgloss"

// CatppuccinMacchiatoTheme implements the Theme interface with Catppuccin Macchiato colors
type CatppuccinMacchiatoTheme struct {
	colors ColorPalette
}

// NewCatppuccinMacchiatoTheme creates a new instance of the Catppuccin Macchiato theme
func NewCatppuccinMacchiatoTheme() Theme {
	return &CatppuccinMacchiatoTheme{
		colors: ColorPalette{
			// Base colors from Catppuccin Macchiato
			Primary:   lipgloss.Color("#c6a0f6"), // mauve
			Secondary: lipgloss.Color("#f5bde6"), // pink
			Tertiary:  lipgloss.Color("#eed49f"), // yellow

			// UI element colors
			Background: lipgloss.Color("#24273a"), // base
			Foreground: lipgloss.Color("#cad3f5"), // text
			Border:     lipgloss.Color("#b7bdf8"), // lavender

			// Semantic colors
			Success: lipgloss.Color("#a6da95"), // green
			Error:   lipgloss.Color("#ed8796"), // red
			Warning: lipgloss.Color("#f5a97f"), // peach
			Info:    lipgloss.Color("#8aadf4"), // blue

			// File/directory colors
			Directory:  lipgloss.Color("#8aadf4"), // blue
			File:       lipgloss.Color("#cad3f5"), // text
			Selected:   lipgloss.Color("#a6da95"), // green
			Deselected: lipgloss.Color("#6e738d"), // overlay0

			// Text colors
			Text:      lipgloss.Color("#cad3f5"), // text
			Muted:     lipgloss.Color("#939ab7"), // overlay2
			Highlight: lipgloss.Color("#f5bde6"), // pink
		},
	}
}

// Colors returns the color palette
func (t *CatppuccinMacchiatoTheme) Colors() ColorPalette {
	return t.colors
}

// Name returns the theme name
func (t *CatppuccinMacchiatoTheme) Name() string {
	return "catppuccin-macchiato"
}
