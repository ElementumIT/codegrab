package themes

import "github.com/charmbracelet/lipgloss"

// CatppuccinFrappeTheme implements the Theme interface with Catppuccin Frappé colors
type CatppuccinFrappeTheme struct {
	colors ColorPalette
}

// NewCatppuccinFrappeTheme creates a new instance of the Catppuccin Frappé theme
func NewCatppuccinFrappeTheme() Theme {
	return &CatppuccinFrappeTheme{
		colors: ColorPalette{
			// Base colors from Catppuccin Frappé
			Primary:   lipgloss.Color("#ca9ee6"), // mauve
			Secondary: lipgloss.Color("#f4b8e4"), // pink
			Tertiary:  lipgloss.Color("#e5c890"), // yellow

			// UI element colors
			Background: lipgloss.Color("#303446"), // base
			Foreground: lipgloss.Color("#c6d0f5"), // text
			Border:     lipgloss.Color("#babbf1"), // lavender

			// Semantic colors
			Success: lipgloss.Color("#a6d189"), // green
			Error:   lipgloss.Color("#e78284"), // red
			Warning: lipgloss.Color("#ef9f76"), // peach
			Info:    lipgloss.Color("#8caaee"), // blue

			// File/directory colors
			Directory:  lipgloss.Color("#8caaee"), // blue
			File:       lipgloss.Color("#f4b8e4"), // pink - more visible than text
			Selected:   lipgloss.Color("#a6d189"), // green
			Deselected: lipgloss.Color("#737994"), // overlay0

			// Text colors
			Text:      lipgloss.Color("#c6d0f5"), // text
			Muted:     lipgloss.Color("#949cbb"), // overlay2
			Highlight: lipgloss.Color("#f4b8e4"), // pink

			HighlightBackground: lipgloss.Color("#414559"), // surface0
		},
	}
}

// Colors returns the color palette
func (t *CatppuccinFrappeTheme) Colors() ColorPalette {
	return t.colors
}

// Name returns the theme name
func (t *CatppuccinFrappeTheme) Name() string {
	return "catppuccin-frappe"
}
