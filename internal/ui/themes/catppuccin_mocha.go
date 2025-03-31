package themes

import "github.com/charmbracelet/lipgloss"

// CatppuccinMochaTheme implements the Theme interface with Catppuccin Mocha colors
type CatppuccinMochaTheme struct {
	colors ColorPalette
}

// NewCatppuccinMochaTheme creates a new instance of the Catppuccin Mocha theme
func NewCatppuccinMochaTheme() Theme {
	return &CatppuccinMochaTheme{
		colors: ColorPalette{
			// Base colors from Catppuccin Mocha
			Primary:   lipgloss.Color("#cba6f7"), // mauve
			Secondary: lipgloss.Color("#f5c2e7"), // pink
			Tertiary:  lipgloss.Color("#f9e2af"), // yellow

			// UI element colors
			Background: lipgloss.Color("#1e1e2e"), // base
			Foreground: lipgloss.Color("#cdd6f4"), // text
			Border:     lipgloss.Color("#b4befe"), // lavender

			// Semantic colors
			Success: lipgloss.Color("#a6e3a1"), // green
			Error:   lipgloss.Color("#f38ba8"), // red
			Warning: lipgloss.Color("#fab387"), // peach
			Info:    lipgloss.Color("#89b4fa"), // blue

			// File/directory colors
			Directory:  lipgloss.Color("#89b4fa"), // blue
			File:       lipgloss.Color("#cdd6f4"), // text
			Selected:   lipgloss.Color("#a6e3a1"), // green
			Deselected: lipgloss.Color("#6c7086"), // overlay0

			// Text colors
			Text:      lipgloss.Color("#cdd6f4"), // text
			Muted:     lipgloss.Color("#9399b2"), // overlay2
			Highlight: lipgloss.Color("#f5c2e7"), // pink
		},
	}
}

// Colors returns the color palette
func (t *CatppuccinMochaTheme) Colors() ColorPalette {
	return t.colors
}

// Name returns the theme name
func (t *CatppuccinMochaTheme) Name() string {
	return "catppuccin-mocha"
}
