package themes

import "github.com/charmbracelet/lipgloss"

// NordTheme implements the Theme interface with Nord colors
type NordTheme struct {
	colors ColorPalette
}

// NewNordTheme creates a new instance of the Nord theme
func NewNordTheme() Theme {
	return &NordTheme{
		colors: ColorPalette{
			// Base colors from Nord
			Primary:   lipgloss.Color("#88c0d0"), // nord8
			Secondary: lipgloss.Color("#81a1c1"), // nord9
			Tertiary:  lipgloss.Color("#ebcb8b"), // nord13

			// UI element colors
			Background: lipgloss.Color("#2e3440"), // nord0
			Foreground: lipgloss.Color("#eceff4"), // nord6
			Border:     lipgloss.Color("#81a1c1"), // nord9

			// Semantic colors
			Success: lipgloss.Color("#a3be8c"), // nord14
			Error:   lipgloss.Color("#bf616a"), // nord11
			Warning: lipgloss.Color("#d08770"), // nord12
			Info:    lipgloss.Color("#5e81ac"), // nord10

			// File/directory colors
			Directory:  lipgloss.Color("#5e81ac"), // nord10
			File:       lipgloss.Color("#eceff4"), // nord6
			Selected:   lipgloss.Color("#a3be8c"), // nord14
			Deselected: lipgloss.Color("#4c566a"), // nord3

			// Text colors
			Text:      lipgloss.Color("#eceff4"), // nord6
			Muted:     lipgloss.Color("#7b88a1"), // between nord3 and nord4
			Highlight: lipgloss.Color("#88c0d0"), // nord8

			HighlightBackground: lipgloss.Color("#3b4252"), // nord1
		},
	}
}

// Colors returns the color palette
func (t *NordTheme) Colors() ColorPalette {
	return t.colors
}

// Name returns the theme name
func (t *NordTheme) Name() string {
	return "nord"
}
