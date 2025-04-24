package themes

import "github.com/charmbracelet/lipgloss"

// CatppuccinLatteTheme implements the Theme interface with Catppuccin Latte colors
type CatppuccinLatteTheme struct {
	colors ColorPalette
}

// NewCatppuccinLatteTheme creates a new instance of the Catppuccin Latte theme
func NewCatppuccinLatteTheme() Theme {
	return &CatppuccinLatteTheme{
		colors: ColorPalette{
			// Base colors from Catppuccin Latte
			Primary:   lipgloss.Color("#8839ef"), // mauve
			Secondary: lipgloss.Color("#ea76cb"), // pink
			Tertiary:  lipgloss.Color("#df8e1d"), // yellow

			// UI element colors
			Background: lipgloss.Color("#eff1f5"), // base
			Foreground: lipgloss.Color("#4c4f69"), // text
			Border:     lipgloss.Color("#7287fd"), // lavender

			// Semantic colors
			Success: lipgloss.Color("#40a02b"), // green
			Error:   lipgloss.Color("#d20f39"), // red
			Warning: lipgloss.Color("#fe640b"), // peach
			Info:    lipgloss.Color("#1e66f5"), // blue

			// File/directory colors
			Directory:  lipgloss.Color("#1e66f5"), // blue
			File:       lipgloss.Color("#4c4f69"), // text
			Selected:   lipgloss.Color("#40a02b"), // green
			Deselected: lipgloss.Color("#9ca0b0"), // overlay0

			// Text colors
			Text:      lipgloss.Color("#4c4f69"), // text
			Muted:     lipgloss.Color("#6c6f85"), // overlay2
			Highlight: lipgloss.Color("#ea76cb"), // pink

			HighlightBackground: lipgloss.Color("#ccd0da"), // surface0
		},
	}
}

// Colors returns the color palette
func (t *CatppuccinLatteTheme) Colors() ColorPalette {
	return t.colors
}

// Name returns the theme name
func (t *CatppuccinLatteTheme) Name() string {
	return "catppuccin-latte"
}
