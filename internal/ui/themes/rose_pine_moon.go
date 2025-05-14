package themes

import "github.com/charmbracelet/lipgloss"

// RosePineMoonTheme implements the Theme interface with Rosé Pine Moon colors
type RosePineMoonTheme struct {
	colors ColorPalette
}

// NewRosePineTheme creates a new instance of the Rosé Pine Moon theme
func NewRosePineMoonTheme() Theme {
	return &RosePineMoonTheme{
		colors: ColorPalette{
			// Base colors from Rosé Pine Moon
			Primary:   lipgloss.Color("#c4a7e7"), // iris
			Secondary: lipgloss.Color("#ea9a97"), // rose
			Tertiary:  lipgloss.Color("#f6c177"), // gold

			// UI element colors
			Background: lipgloss.Color("#232136"), // base
			Foreground: lipgloss.Color("#2a273f"), // surface
			Border:     lipgloss.Color("#56526e"), // highlight high

			// Semantic colors
			Success: lipgloss.Color("#3e8fb0"), // pine
			Error:   lipgloss.Color("#eb6f92"), // love
			Warning: lipgloss.Color("#f6c177"), // gold
			Info:    lipgloss.Color("#9ccfd8"), // foam

			// File/directory colors
			Directory:  lipgloss.Color("#9ccfd8"), // foam
			File:       lipgloss.Color("#2a273f"), // surface
			Selected:   lipgloss.Color("#ea9a97"), // rose
			Deselected: lipgloss.Color("#6e6a86"), // muted

			// Text colors
			Text:      lipgloss.Color("#e0def4"), // text
			Muted:     lipgloss.Color("#6e6a86"), // muted
			Highlight: lipgloss.Color("#ea9a97"), // rose

			HighlightBackground: lipgloss.Color("#44415a"), // highlight med
		},
	}
}

// Colors returns the color palette
func (t *RosePineMoonTheme) Colors() ColorPalette {
	return t.colors
}

// Name returns the theme name
func (t *RosePineMoonTheme) Name() string {
	return "rose-pine-moon"
}
