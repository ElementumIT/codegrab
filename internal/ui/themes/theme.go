package themes

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// ColorPalette defines the colors used throughout the application
type ColorPalette struct {
	// Base colors
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Tertiary  lipgloss.Color

	// UI element colors
	Background lipgloss.Color
	Foreground lipgloss.Color
	Border     lipgloss.Color

	// Semantic colors
	Success lipgloss.Color
	Error   lipgloss.Color
	Warning lipgloss.Color
	Info    lipgloss.Color

	// File/directory colors
	Directory  lipgloss.Color
	File       lipgloss.Color
	Selected   lipgloss.Color
	Deselected lipgloss.Color

	// Text colors
	Text      lipgloss.Color
	Muted     lipgloss.Color
	Highlight lipgloss.Color

	HighlightBackground lipgloss.Color
}

// Theme interface provides colors for the UI
type Theme interface {
	// Colors returns the color palette for the theme
	Colors() ColorPalette
	// Name returns the theme's name
	Name() string
}

var (
	// CurrentTheme holds the active theme
	CurrentTheme Theme

	// themeConstructors maps theme names to their constructors
	themeConstructors = map[string]func() Theme{}
)

// RegisterTheme registers a theme constructor
func RegisterTheme(name string, constructor func() Theme) {
	themeConstructors[name] = constructor
}

// Initialize registers theme constructors and sets the default theme
func Initialize() {
	RegisterTheme("catppuccin-latte", NewCatppuccinLatteTheme)
	RegisterTheme("catppuccin-frappe", NewCatppuccinFrappeTheme)
	RegisterTheme("catppuccin-macchiato", NewCatppuccinMacchiatoTheme)
	RegisterTheme("catppuccin-mocha", NewCatppuccinMochaTheme)
	RegisterTheme("dracula", NewDraculaTheme)
	RegisterTheme("nord", NewNordTheme)

	// Set default theme
	SetTheme("catppuccin-mocha")
}

// SetTheme changes the current theme by name
func SetTheme(name string) error {
	constructor, exists := themeConstructors[name]
	if !exists {
		return fmt.Errorf("theme %q not found", name)
	}

	CurrentTheme = constructor()
	return nil
}

// GetThemeNames returns a list of available theme names
func GetThemeNames() []string {
	names := make([]string, 0, len(themeConstructors))
	for name := range themeConstructors {
		names = append(names, name)
	}
	return names
}
