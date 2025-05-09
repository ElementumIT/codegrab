package themes

import "github.com/charmbracelet/lipgloss"

// RosePineDawnTheme implements the Theme interface with Rosé Pine Dawn colors
type RosePineDawnTheme struct {
	colors ColorPalette
}

// NewRosePineTheme creates a new instance of the Rosé Pine Dawn theme
func NewRosePineDawnTheme() Theme {
	return &RosePineDawnTheme{
		colors: ColorPalette{
			// Base colors from Rosé Pine Dawn
			Primary:   lipgloss.Color("#907aa9"), // iris
			Secondary: lipgloss.Color("#d7827e"), // rose
			Tertiary:  lipgloss.Color("#ea9d34"), // gold

			// UI element colors
			Background: lipgloss.Color("#faf4ed"), // base
			Foreground: lipgloss.Color("#fffaf3"), // surface
			Border:     lipgloss.Color("#cecacd"), // highlight high

			// Semantic colors
			Success: lipgloss.Color("#286983"), // pine
			Error:   lipgloss.Color("#b4637a"), // love
			Warning: lipgloss.Color("#ea9d34"), // gold
			Info:    lipgloss.Color("#56949f"), // foam

			// File/directory colors
			Directory:  lipgloss.Color("#56949f"), // foam
			File:       lipgloss.Color("#d7827e"), // rose - more visible than surface
			Selected:   lipgloss.Color("#d7827e"), // rose
			Deselected: lipgloss.Color("#9893a5"), // muted

			// Text colors
			Text:      lipgloss.Color("#575279"), // text
			Muted:     lipgloss.Color("#9893a5"), // muted
			Highlight: lipgloss.Color("#d7827e"), // rose

			HighlightBackground: lipgloss.Color("#dfdad9"), // highlight med
		},
	}
}

// Colors returns the color palette
func (t *RosePineDawnTheme) Colors() ColorPalette {
	return t.colors
}

// Name returns the theme name
func (t *RosePineDawnTheme) Name() string {
	return "rose-pine-dawn"
}
