package themes

import "github.com/charmbracelet/lipgloss"

// RosePineTheme implements the Theme interface with Rosé Pine colors
type RosePineTheme struct {
	colors ColorPalette
}

// NewRosePineTheme creates a new instance of the Rosé Pine theme
func NewRosePineTheme() Theme {
	return &RosePineTheme{
		colors: ColorPalette{
			// Base colors from Rosé Pine
			Primary:   lipgloss.Color("#c4a7e7"), // iris
			Secondary: lipgloss.Color("#ebbcba"), // rose
			Tertiary:  lipgloss.Color("#f6c177"), // gold

			// UI element colors
			Background: lipgloss.Color("#191724"), // base
			Foreground: lipgloss.Color("#1f1d2e"), // surface
			Border:     lipgloss.Color("#524f67"), // highlight high

			// Semantic colors
			Success: lipgloss.Color("#31748f"), // pine
			Error:   lipgloss.Color("#eb6f92"), // love
			Warning: lipgloss.Color("#f6c177"), // gold
			Info:    lipgloss.Color("#9ccfd8"), // foam

			// File/directory colors
			Directory:  lipgloss.Color("#9ccfd8"), // foam
			File:       lipgloss.Color("#1f1d2e"), // surface
			Selected:   lipgloss.Color("#ebbcba"), // rose
			Deselected: lipgloss.Color("#908caa"), // muted

			// Text colors
			Text:      lipgloss.Color("#e0def4"), // text
			Muted:     lipgloss.Color("#908caa"), // muted
			Highlight: lipgloss.Color("#ebbcba"), // rose

			HighlightBackground: lipgloss.Color("#403d52"), // highlight med
		},
	}
}

// Colors returns the color palette
func (t *RosePineTheme) Colors() ColorPalette {
	return t.colors
}

// Name returns the theme name
func (t *RosePineTheme) Name() string {
	return "rose-pine"
}
