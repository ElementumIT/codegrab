package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/epilande/codegrab/internal/ui/themes"
)

// GetStyleHeader returns the header style using the current theme
func GetStyleHeader() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(themes.CurrentTheme.Colors().Primary).
		Padding(0, 1)
}

// GetStyleFormatIndicator returns the format indicator style using the current theme
func GetStyleFormatIndicator() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(themes.CurrentTheme.Colors().Highlight)
}

// GetStyleSuccess returns the success style using the current theme
func GetStyleSuccess() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(themes.CurrentTheme.Colors().Success).
		Bold(true)
}

// GetStyleError returns the error style using the current theme
func GetStyleError() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(themes.CurrentTheme.Colors().Error).
		Bold(true)
}

// GetStyleHelp returns the help style using the current theme
func GetStyleHelp() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(themes.CurrentTheme.Colors().Muted)
}

// GetStyleBorderedViewport returns the bordered viewport style using the current theme
func GetStyleBorderedViewport() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(themes.CurrentTheme.Colors().Border)
}

// GetStyleSearchCount returns the search count style using the current theme
func GetStyleSearchCount() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(themes.CurrentTheme.Colors().Muted)
}

// GetStyleWarning returns the warning style using the current theme
func GetStyleWarning() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(themes.CurrentTheme.Colors().Warning).
		Bold(false)
}

// GetStyleInfo returns the info style using the current theme
func GetStyleInfo() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(themes.CurrentTheme.Colors().Info).
		Bold(false)
}

// StyleCheckBox returns a styled checkbox based on the current theme
func StyleCheckBox(checked bool) string {
	if checked {
		return lipgloss.NewStyle().
			Foreground(themes.CurrentTheme.Colors().Success).
			Render("[x]")
	}
	return lipgloss.NewStyle().
		Foreground(themes.CurrentTheme.Colors().Text).
		Render("[ ]")
}

// StylePartialCheckBox returns a styled partial checkbox for directories with selected children
func StylePartialCheckBox() string {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(themes.CurrentTheme.Colors().Info).
		Render("[~]")
}

func StyleDirectoryName(name string) string {
	return lipgloss.NewStyle().
		Foreground(themes.CurrentTheme.Colors().Directory).
		Render(name)
}

// StyleFileLine styles a file line based on its properties and the current theme
func StyleFileLine(content string, isDir, isSelected, isDeselected, isCursor bool) string {
	colors := themes.CurrentTheme.Colors()
	style := lipgloss.NewStyle()

	if isDir {
		style = style.Foreground(colors.Directory)
		if isSelected {
			style = style.Bold(true).Foreground(colors.Border)
		}
	} else {
		if isSelected && !isDeselected {
			style = style.Bold(true).Foreground(colors.Selected)
		} else if isDeselected {
			style = style.Foreground(colors.Deselected)
		} else {
			style = style.Foreground(colors.Text)
		}
	}

	if isCursor {
		return style.Bold(true).Render(" ❯ " + content)
	}
	return style.Render("   " + content)
}

// NewSearchInput creates a new search input with styling from the current theme
func NewSearchInput() textinput.Model {
	colors := themes.CurrentTheme.Colors()

	ti := textinput.New()
	ti.Placeholder = "Search files..."
	ti.PromptStyle = ti.PromptStyle.
		Foreground(colors.Tertiary).
		Padding(0, 1)
	ti.TextStyle = ti.TextStyle.Foreground(colors.Tertiary)
	ti.Prompt = "🔍"
	ti.Width = 50
	return ti
}
