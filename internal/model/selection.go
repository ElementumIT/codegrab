package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/epilande/codegrab/internal/utils"
)

func (m *Model) filterSelections() {
	for p := range m.selected {
		if (!m.showHidden && utils.IsHiddenPath(p)) ||
			(m.useGitIgnore && m.gitIgnoreMgr.IsIgnored(filepath.Join(m.rootPath, p))) {
			delete(m.selected, p)
		}
	}
	for p := range m.deselected {
		if (!m.showHidden && utils.IsHiddenPath(p)) ||
			(m.useGitIgnore && m.gitIgnoreMgr.IsIgnored(filepath.Join(m.rootPath, p))) {
			delete(m.deselected, p)
		}
	}
}

// Helper to find the nearest selected parent directory
func findParentDirectory(path, root string, selected map[string]bool) string {
	current := filepath.Dir(path)
	for current != root && current != "." && current != path {
		if selected[current] {
			return current
		}
		path = current
		current = filepath.Dir(path)
	}
	return ""
}

func (m *Model) toggleSelection(path string, isDir bool) {
	if path == "" {
		return
	}

	fullPath := filepath.Join(m.rootPath, path)
	if _, err := os.Stat(fullPath); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: cannot access %s: %v\n", path, err)
		return
	}

	if isDir {
		if m.selected[path] {
			delete(m.selected, path)
			for _, f := range m.files {
				if strings.HasPrefix(f.Path, path+"/") {
					delete(m.selected, f.Path)
					m.deselected[f.Path] = true
				}
			}
		} else {
			m.selected[path] = true

			// If we're in search mode, only select files that are in the search results
			if m.isSearching && len(m.searchResults) > 0 {
				searchResultPaths := make(map[string]bool)
				for _, node := range m.searchResults {
					if !node.IsDir {
						searchResultPaths[node.Path] = true
					}
				}

				// Process all files in the directory
				for _, f := range m.files {
					if strings.HasPrefix(f.Path, path+"/") && !f.IsDir {
						if searchResultPaths[f.Path] {
							// File is in search results - select it
							m.selected[f.Path] = true
							delete(m.deselected, f.Path)
						} else {
							// File is NOT in search results - deselect it
							delete(m.selected, f.Path)
							m.deselected[f.Path] = true
						}
					}
				}
			} else {
				for _, f := range m.files {
					if strings.HasPrefix(f.Path, path+"/") {
						m.selected[f.Path] = true
						delete(m.deselected, f.Path)
					}
				}
			}
		}
	} else {
		if m.selected[path] {
			delete(m.selected, path)
			// If under a selected directory, mark as deselected
			if parentDir := findParentDirectory(path, m.rootPath, m.selected); parentDir != "" {
				m.deselected[path] = true
			}
		} else {
			m.selected[path] = true
			delete(m.deselected, path)
		}
	}
}

func (m *Model) expandAllDirectories() {
	m.collapsed = make(map[string]bool)
	m.buildDisplayNodes()
}

func (m *Model) collapseAllDirectories() {
	for _, node := range m.files {
		if node.IsDir {
			m.collapsed[node.Path] = true
		}
	}
	m.buildDisplayNodes()
}
