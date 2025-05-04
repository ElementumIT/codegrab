package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/epilande/codegrab/internal/ui"
	"github.com/epilande/codegrab/internal/utils"
)

func (m Model) View() string {
	if m.showHelp {
		header := ui.GetStyleHeader().Render("❔ Help Menu")
		content := ui.GetStyleBorderedViewport().Render(m.viewport.View())
		footer := ui.GetStyleHelp().Render("Exit: esc")
		return header + "\n" + content + "\n" + footer
	}

	header := m.renderHeader()
	footer := m.renderFooter()
	content := ui.GetStyleBorderedViewport().Render(m.viewport.View())
	return header + "\n" + content + "\n" + footer
}

func (m Model) renderHeader() string {
	headerIcon := "✋"
	if m.isGrabbing {
		headerIcon = "✊"
	}
	leftContent := ui.GetStyleHeader().Render(headerIcon + " Code Grab")
	rightContent := ""
	totalFiles := m.getTotalFileCount()
	selectedCount := m.getSelectedFileCount()

	formatExt := strings.TrimPrefix(m.generator.GetFormat().Extension(), ".")
	formatIndicator := ui.GetStyleFormatIndicator().Render(formatExt)

	if m.isSearching {
		leftContent = m.searchInput.View()
		matchCount := 0
		for _, node := range m.searchResults {
			if !node.IsDir {
				matchCount++
			}
		}
		rightContent = ui.GetStyleSearchCount().Render(fmt.Sprintf("%d/%d [%d] %s",
			matchCount,
			totalFiles,
			selectedCount,
			formatIndicator,
		))
	} else {
		rightContent = ui.GetStyleSearchCount().Render(fmt.Sprintf("%d [%d] %s",
			totalFiles,
			selectedCount,
			formatIndicator,
		))
	}

	leftWidth := lipgloss.Width(leftContent)
	rightWidth := lipgloss.Width(rightContent)
	spacing := m.width - leftWidth - rightWidth - 1

	if spacing < 1 {
		spacing = 1
	}

	header := lipgloss.JoinHorizontal(
		lipgloss.Center,
		leftContent,
		strings.Repeat(" ", spacing),
		rightContent,
	)

	return header
}

func (m Model) renderFooter() string {
	if m.showHelp {
		return ""
	}

	var leftParts []string
	var rightParts []string

	// Left side: Status/Error/Help
	if m.isSearching {
		searchHelp := "Next: ctrl+n | Prev: ctrl+p | Select: tab | Exit: esc"
		leftParts = append(leftParts, ui.GetStyleHelp().Render(searchHelp))
	} else if m.err != nil {
		leftParts = append(leftParts, ui.GetStyleError().Render(m.err.Error()))
	} else if m.successMsg != "" {
		leftParts = append(leftParts, ui.GetStyleSuccess().Render(m.successMsg))
	} else {
		helpText := "Press '?' for help | Select: space | Generate: g | Copy: y"
		leftParts = append(leftParts, ui.GetStyleHelp().Render(helpText))
	}

	// Right side: Warn/Redaction status
	redactionStatus := ""
	if m.warningMsg != "" {
		rightParts = append(rightParts, ui.GetStyleWarning().Render(m.warningMsg))
	} else if m.redactSecrets {
		redactionStatus = ui.GetStyleInfo().Render("🛡️ Redacting")
	} else {
		redactionStatus = ui.GetStyleWarning().Render("⚠️ NOT Redacting")
	}
	rightParts = append(rightParts, redactionStatus)

	depStatus := ""
	if m.resolveDeps {
		depStatus = ui.GetStyleInfo().Render(" | 🔗 Deps")
	}
	rightParts = append(rightParts, depStatus)

	leftContent := lipgloss.JoinHorizontal(lipgloss.Top, leftParts...)
	rightContent := lipgloss.JoinHorizontal(lipgloss.Top, rightParts...)

	leftWidth := lipgloss.Width(leftContent)
	rightWidth := lipgloss.Width(rightContent)
	availableWidth := m.width - 2

	spacing := availableWidth - leftWidth - rightWidth
	if spacing < 1 {
		spacing = 1
	}

	footer := lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		leftContent,
		lipgloss.NewStyle().Width(spacing).Render(""),
		rightContent,
	)
	return lipgloss.NewStyle().Padding(0, 1).Render(footer)
}

// refreshViewportContent regenerates the lines for our displayNodes, highlights
// the cursor, and sets that as the viewport content.
func (m *Model) refreshViewportContent() {
	var nodes []FileNode
	if m.isSearching && len(m.searchResults) > 0 {
		nodes = m.searchResults
	} else {
		nodes = m.displayNodes
	}

	dirsWithSelectedChildren := make(map[string]bool)
	for path := range m.selected {
		dir := filepath.Dir(path)
		for dir != "." && dir != "/" {
			dirsWithSelectedChildren[dir] = true
			dir = filepath.Dir(dir)
		}
	}

	// Calculate directory selected counts
	dirSelectedCounts := make(map[string]int)
	for _, file := range m.files {
		if !file.IsDir && m.selected[file.Path] && !m.deselected[file.Path] {
			dir := filepath.Dir(file.Path)
			for dir != "." && dir != "/" {
				dirSelectedCounts[dir]++
				dir = filepath.Dir(dir)
			}
		}
	}

	var lines []string
	parentIsLast := make(map[int]bool)

	for i, node := range nodes {
		var prefixBuilder strings.Builder
		for l := 0; l < node.Level; l++ {
			if parentIsLast[l] {
				prefixBuilder.WriteString("    ")
			} else {
				prefixBuilder.WriteString("│   ")
			}
		}

		treeBranch := "├── "
		if node.IsLast {
			treeBranch = "└── "
			parentIsLast[node.Level] = true
		} else {
			parentIsLast[node.Level] = false
		}
		for l := node.Level + 1; l < len(parentIsLast); l++ {
			parentIsLast[l] = false
		}

		treePrefix := prefixBuilder.String() + treeBranch

		icon := node.Icon
		iconColor := node.IconColor
		if node.IsDir && icon == "" {
			if m.collapsed[node.Path] {
				icon = ""
			} else {
				icon = ""
			}
			// Fallback to use theme directory color
			iconColor = ""
		}

		isPartialDir := !m.selected[node.Path] && dirsWithSelectedChildren[node.Path] && node.IsDir
		rawCheckbox := "[ ]"
		if node.IsDir {
			if m.selected[node.Path] {
				rawCheckbox = "[x]"
			} else if isPartialDir {
				rawCheckbox = "[~]"
			}
		} else {
			if node.Selected {
				rawCheckbox = "[x]"
			}
		}

		rawName := node.Name
		if node.IsDir {
			rawName += "/"
		}

		rawSuffix := ""
		if node.IsDir {
			if count := dirSelectedCounts[node.Path]; count > 0 {
				rawSuffix = fmt.Sprintf(" [%d]", count)
			}
		} else {
			if node.IsDependency {
				rawSuffix += " [dep]"
			}
			if m.showTokenCount {
				if ok, err := utils.IsTextFile(node.Path); ok && err == nil {
					if contentBytes, err := os.ReadFile(node.Path); err == nil {
						content := string(contentBytes)
						tokensEstimate := utils.EstimateTokens(content)
						rawSuffix += fmt.Sprintf(" [%d tokens]", tokensEstimate)
					}
				}
			}
		}
		isCursorLine := i == m.cursor

		rendered := ui.StyleFileLine(
			rawCheckbox,
			treePrefix,
			icon,
			iconColor,
			rawName,
			rawSuffix,
			node.IsDir,
			m.selected[node.Path],
			isCursorLine,
			isPartialDir,
			m.viewport.Width,
		)
		lines = append(lines, rendered)
	}

	m.viewport.SetContent(strings.Join(lines, "\n"))
}

// ensureCursorVisible adjusts viewport.YOffset so the selected line is visible.
func (m *Model) ensureCursorVisible() {
	if m.cursor < m.viewport.YOffset {
		m.viewport.YOffset = m.cursor
	} else if m.cursor >= m.viewport.YOffset+m.viewport.Height {
		m.viewport.YOffset = m.cursor - m.viewport.Height + 1
		if m.viewport.YOffset < 0 {
			m.viewport.YOffset = 0
		}
	}
}

func (m *Model) showHelpScreen() {
	helpContent := ui.GetStyleHelp().Render(ui.HelpText + "\n\nPress '?' to close this help menu.")
	m.viewport.SetContent(helpContent)
	m.viewport.GotoTop() // Reset scroll position to top when help is opened
}

func (m *Model) getTotalFileCount() int {
	count := 0
	for _, file := range m.files {
		if !file.IsDir {
			count++
		}
	}
	return count
}

// getSelectedFileCount calculates the effective number of selected files.
// This considers explicitly selected files, files within selected directories
// (respecting search results), and excludes deselected files.
func (m *Model) getSelectedFileCount() int {
	effectiveSelection := make(map[string]bool)

	searchResultPaths := make(map[string]bool)
	useSearchResults := m.isSearching && len(m.searchResults) > 0
	if useSearchResults {
		for _, node := range m.searchResults {
			searchResultPaths[node.Path] = true
		}
	}

	for _, item := range m.files {
		if effectiveSelection[item.Path] {
			continue
		}

		if m.deselected[item.Path] {
			continue
		}

		if useSearchResults && !searchResultPaths[item.Path] {
			isChildOfSearchResult := false
			parent := filepath.Dir(item.Path)
			for parent != "." && parent != "/" {
				if searchResultPaths[parent] {
					isChildOfSearchResult = true
					break
				}
				parent = filepath.Dir(parent)
			}
			if !isChildOfSearchResult {
				continue
			}
		}

		if m.selected[item.Path] {
			if !item.IsDir {
				effectiveSelection[item.Path] = true
			} else {
				prefix := item.Path + string(os.PathSeparator)
				for _, child := range m.files {
					if !child.IsDir && strings.HasPrefix(child.Path, prefix) && !m.deselected[child.Path] {
						if useSearchResults && !searchResultPaths[child.Path] {
							continue
						}
						effectiveSelection[child.Path] = true
					}
				}
			}
			continue
		}

		if !item.IsDir {
			parent := filepath.Dir(item.Path)
			isInSelectedDir := false
			for parent != "." && parent != "/" {
				if m.selected[parent] {
					isInSelectedDir = true
					break
				}
				parent = filepath.Dir(parent)
			}
			if isInSelectedDir {
				effectiveSelection[item.Path] = true
			}
		}
	}

	return len(effectiveSelection)
}
