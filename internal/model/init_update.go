package model

import (
	"fmt"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/epilande/codegrab/internal/filesystem"
	"github.com/epilande/codegrab/internal/generator/formats"
)

type filesLoadedMsg struct {
	err   error
	files []filesystem.FileItem
}

type outputGeneratedMsg struct {
	err        error
	path       string
	format     string
	tokenCount int
}

type clipboardCopiedMsg struct {
	err        error
	format     string
	tokenCount int
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		files, err := filesystem.WalkDirectory(m.rootPath, m.gitIgnoreMgr, m.filterMgr, m.useGitIgnore, m.showHidden)
		if err != nil {
			return filesLoadedMsg{files: nil, err: fmt.Errorf("failed to walk directory: %w", err)}
		}
		return filesLoadedMsg{files: files, err: nil}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.successMsg = ""
	m.isGrabbing = false

	switch msg := msg.(type) {
	case filesLoadedMsg:
		m.err = msg.err
		m.files = msg.files
		for _, f := range m.files {
			if f.IsDir {
				m.collapsed[f.Path] = true
			}
		}
		m.buildDisplayNodes()
		m.refreshViewportContent()
		return m, nil

	case outputGeneratedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.successMsg = ""
		} else {
			m.isGrabbing = true
			m.err = nil
			m.successMsg = fmt.Sprintf("✅ %s generated: %s (%d tokens)", msg.format, msg.path, msg.tokenCount)
		}
		m.refreshViewportContent()
		return m, nil

	case clipboardCopiedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.successMsg = ""
		} else {
			m.isGrabbing = true
			m.err = nil
			m.successMsg = fmt.Sprintf("✅ %s copied to clipboard! (%d tokens)", msg.format, msg.tokenCount)
		}
		m.refreshViewportContent()
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		headerHeight, footerHeight := 2, 2
		m.viewport.Width = m.width - 2
		m.viewport.Height = m.height - headerHeight - footerHeight
		m.refreshViewportContent()
		return m, nil

	case tea.KeyMsg:
		if m.showHelp {
			switch msg.String() {
			case "q", "esc", "?":
				m.showHelp = false
				m.refreshViewportContent()
				return m, nil
			case "j", "down":
				m.viewport.LineDown(1)
				return m, nil
			case "k", "up":
				m.viewport.LineUp(1)
				return m, nil
			case "g", "home":
				m.viewport.GotoTop()
				return m, nil
			case "G", "end":
				m.viewport.GotoBottom()
				return m, nil
			default:
				return m, nil
			}
		}

		if m.isSearching {
			var cmd tea.Cmd

			switch msg.String() {
			case "esc":
				m.isSearching = false
				m.searchInput.Blur()
				m.searchInput.SetValue("")
				m.searchResults = nil
				m.collapseAllDirectories()
				m.refreshViewportContent()
				return m, nil
			case "tab", "enter":
				if len(m.searchResults) > 0 && m.cursor < len(m.searchResults) {
					node := m.searchResults[m.cursor]
					m.toggleSelection(node.Path, node.IsDir)
					m.buildDisplayNodes()
					m.updateSearchResults()
					m.ensureCursorVisible()
					m.refreshViewportContent()
				}
				return m, nil
			case "ctrl+n", "down":
				if len(m.searchResults) > 0 {
					m.cursor = (m.cursor + 1) % len(m.searchResults)
					m.ensureCursorVisible()
					m.refreshViewportContent()
				}
				return m, nil
			case "ctrl+p", "up":
				if len(m.searchResults) > 0 {
					m.cursor--
					if m.cursor < 0 {
						m.cursor = len(m.searchResults) - 1
					}
					m.ensureCursorVisible()
					m.refreshViewportContent()
				}
				return m, nil
			}

			m.searchInput, cmd = m.searchInput.Update(msg)

			m.updateSearchResults()
			m.refreshViewportContent()
			return m, cmd
		}

		switch key := msg.String(); key {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "j", "down", "ctrl+n":
			if m.cursor < len(m.displayNodes)-1 {
				m.cursor++
				m.ensureCursorVisible()
				m.refreshViewportContent()
			}
		case "k", "up", "ctrl+p":
			if m.cursor > 0 {
				m.cursor--
				m.ensureCursorVisible()
				m.refreshViewportContent()
			}
		case "H", "home":
			m.cursor = 0
			m.ensureCursorVisible()
			m.refreshViewportContent()
		case "L", "end":
			m.cursor = len(m.displayNodes) - 1
			m.ensureCursorVisible()
			m.refreshViewportContent()
		case "h", "left":
			if m.cursor < len(m.displayNodes) {
				node := m.displayNodes[m.cursor]
				if node.IsDir && !m.collapsed[node.Path] {
					m.toggleCollapse(node.Path)
					m.buildDisplayNodes()
					if m.cursor >= len(m.displayNodes) {
						m.cursor = len(m.displayNodes) - 1
					}
					m.ensureCursorVisible()
					m.refreshViewportContent()
				}
			}
		case "l", "right":
			if m.cursor < len(m.displayNodes) {
				node := m.displayNodes[m.cursor]
				if node.IsDir && m.collapsed[node.Path] {
					m.toggleCollapse(node.Path)
					m.buildDisplayNodes()
					m.ensureCursorVisible()
					m.refreshViewportContent()
				}
			}
		case "e":
			if len(m.collapsed) > 0 {
				m.expandAllDirectories()
			} else {
				m.collapseAllDirectories()
			}
			m.ensureCursorVisible()
			m.refreshViewportContent()
		case " ", "tab":
			if m.cursor < len(m.displayNodes) {
				node := m.displayNodes[m.cursor]
				m.toggleSelection(node.Path, node.IsDir)
				m.buildDisplayNodes()
				m.ensureCursorVisible()
				m.refreshViewportContent()
			}
		case "/":
			m.isSearching = true
			m.searchInput.Focus()
			m.searchInput.SetValue("")
			m.searchResults = nil
			m.expandAllDirectories()
			m.refreshViewportContent()
			return m, nil
		case "?":
			m.showHelp = !m.showHelp
			if m.showHelp {
				m.showHelpScreen()
			} else {
				m.refreshViewportContent()
			}
		case "esc":
			m.showHelp = false
			m.refreshViewportContent()
		case "y":
			return m, m.copyOutputToClipboard()
		case "g":
			return m, m.generateOutput()
		case "i":
			m.useGitIgnore = !m.useGitIgnore
			m.generator.UseGitIgnore = m.useGitIgnore
			return m, m.reloadFiles()
		case ".":
			m.showHidden = !m.showHidden
			m.generator.ShowHidden = m.showHidden
			return m, m.reloadFiles()
		case "F":
			formatNames := formats.GetFormatNames()
			if len(formatNames) == 0 {
				break
			}

			currentFormatName := m.generator.GetFormatName()

			nextIndex := 0
			for i, name := range formatNames {
				if name == currentFormatName {
					nextIndex = (i + 1) % len(formatNames)
					break
				}
			}

			nextFormat := formats.GetFormat(formatNames[nextIndex])
			m.generator.SetFormat(nextFormat)

			m.successMsg = fmt.Sprintf("Format changed: %s", m.generator.GetFormatName())
			m.refreshViewportContent()
		}
	}
	return m, nil
}

func (m *Model) reloadFiles() tea.Cmd {
	return func() tea.Msg {
		m.filterSelections()
		files, err := filesystem.WalkDirectory(m.rootPath, m.gitIgnoreMgr, m.filterMgr, m.useGitIgnore, m.showHidden)
		if err != nil {
			return filesLoadedMsg{files: nil, err: fmt.Errorf("failed to reload files: %w", err)}
		}
		return filesLoadedMsg{files: files, err: nil}
	}
}

func (m *Model) toggleCollapse(path string) {
	m.collapsed[path] = !m.collapsed[path]
}

// generateOutput creates the output in the current format
func (m Model) generateOutput() tea.Cmd {
	return func() tea.Msg {
		m.generator.SelectedFiles = m.selected
		m.generator.DeselectedFiles = m.deselected
		outPath, tokenCount, err := m.generator.Generate()
		if err != nil {
			return outputGeneratedMsg{
				err:        fmt.Errorf("failed to generate output: %w", err),
				path:       "",
				tokenCount: 0,
				format:     m.generator.GetFormatName(),
			}
		}
		return outputGeneratedMsg{
			err:        nil,
			path:       outPath,
			tokenCount: tokenCount,
			format:     m.generator.GetFormatName(),
		}
	}
}

// copyOutputToClipboard copies the generated content to the clipboard
func (m Model) copyOutputToClipboard() tea.Cmd {
	return func() tea.Msg {
		m.generator.SelectedFiles = m.selected
		m.generator.DeselectedFiles = m.deselected

		content, tokenCount, err := m.generator.GenerateString()
		if err != nil {
			return clipboardCopiedMsg{
				err:        err,
				tokenCount: 0,
				format:     m.generator.GetFormatName(),
			}
		}

		if err = clipboard.WriteAll(content); err != nil {
			return clipboardCopiedMsg{
				err:        fmt.Errorf("failed to copy to clipboard: %w", err),
				tokenCount: 0,
				format:     m.generator.GetFormatName(),
			}
		}

		return clipboardCopiedMsg{
			err:        nil,
			tokenCount: tokenCount,
			format:     m.generator.GetFormatName(),
		}
	}
}
