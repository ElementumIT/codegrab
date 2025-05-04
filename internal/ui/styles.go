package ui

import (
	"fmt"
	"strings"

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

// StyleFileLine styles a file line based on its properties and the current theme
func StyleFileLine(
	rawCheckbox string,
	treePrefix string,
	icon string,
	iconColor string,
	name string,
	rawSuffix string,
	isDir bool,
	isSelected bool,
	isCursor bool,
	isPartialDir bool,
	viewportWidth int,
) string {
	colors := themes.CurrentTheme.Colors()

	checkboxStyle := lipgloss.NewStyle()
	prefixStyle := lipgloss.NewStyle().Foreground(colors.Text)
	nameStyle := lipgloss.NewStyle()
	suffixStyle := lipgloss.NewStyle().Foreground(colors.Muted)
	iconStyle := lipgloss.NewStyle()

	switch rawCheckbox {
	case "[x]":
		checkboxStyle = checkboxStyle.Foreground(colors.Success)
	case "[~]":
		checkboxStyle = checkboxStyle.Foreground(colors.Info)
	default:
		checkboxStyle = checkboxStyle.Foreground(colors.Muted)
	}

	if iconColor != "" {
		iconStyle = iconStyle.Foreground(lipgloss.Color(iconColor))
	} else {
		if isDir {
			iconStyle = iconStyle.Foreground(colors.Directory)
		} else {
			iconStyle = iconStyle.Foreground(colors.File)
		}
	}

	shouldBold := isSelected || isPartialDir
	if isDir {
		nameStyle = nameStyle.Foreground(colors.Directory)
	} else {
		if isSelected {
			nameStyle = nameStyle.Foreground(colors.Selected)
		} else {
			nameStyle = nameStyle.Foreground(colors.Text)
		}
	}
	if shouldBold {
		nameStyle = nameStyle.Bold(true)
		checkboxStyle = checkboxStyle.Bold(true)
	}

	renderedCheckbox := checkboxStyle.PaddingRight(1).Render(rawCheckbox)
	renderedPrefix := prefixStyle.Render(treePrefix)
	renderedIcon := ""
	if icon != "" {
		renderedIcon = iconStyle.Render(icon + " ")
	}
	renderedName := nameStyle.Render(name)
	renderedSuffix := suffixStyle.Render(rawSuffix)

	if isCursor {
		cursorBaseStyle := lipgloss.NewStyle().Background(colors.HighlightBackground).Bold(true)
		cursorIndicator := cursorBaseStyle.Foreground(colors.Text).Render(" ❯ ")

		cursorCheckboxStyle := checkboxStyle
		if rawCheckbox == "[ ]" {
			cursorCheckboxStyle = cursorCheckboxStyle.Foreground(colors.Text)
		}

		cursorIconStyle := iconStyle

		renderedCheckbox = cursorBaseStyle.Inherit(cursorCheckboxStyle).PaddingRight(1).Render(rawCheckbox)
		renderedPrefix = cursorBaseStyle.Inherit(prefixStyle).Render(treePrefix)
		if icon != "" {
			renderedIcon = cursorBaseStyle.Inherit(cursorIconStyle).Bold(false).Render(icon + " ")
		} else {
			renderedIcon = ""
		}
		renderedName = cursorBaseStyle.Inherit(nameStyle).Render(name)
		renderedSuffix = cursorBaseStyle.Inherit(suffixStyle).Render(rawSuffix)

		cursorLineContent := fmt.Sprintf("%s%s%s%s%s",
			renderedCheckbox,
			renderedPrefix,
			renderedIcon,
			renderedName,
			renderedSuffix,
		)

		fullContentWidth := lipgloss.Width(cursorIndicator + cursorLineContent)
		paddingWidth := viewportWidth - fullContentWidth
		if paddingWidth < 0 {
			paddingWidth = 0
		}
		padding := cursorBaseStyle.Render(strings.Repeat(" ", paddingWidth))

		return cursorIndicator + cursorLineContent + padding
	} else {
		lineContent := fmt.Sprintf("%s%s%s%s%s",
			renderedCheckbox,
			renderedPrefix,
			renderedIcon,
			renderedName,
			renderedSuffix,
		)

		return "   " + lineContent
	}
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
