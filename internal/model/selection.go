package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/epilande/codegrab/internal/dependencies"
	"github.com/epilande/codegrab/internal/utils"
)

type QueuedDep struct {
	Path  string
	Depth int
}

func (m *Model) filterSelections() {
	for p := range m.selected {
		if (!m.showHidden && utils.IsHiddenPath(p)) ||
			(m.useGitIgnore && m.gitIgnoreMgr.IsIgnored(filepath.Join(m.rootPath, p))) {
			delete(m.selected, p)
			delete(m.isDependency, p)
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
	// Stop if current is empty, root, or the path itself
	for current != "" && current != "." && current != "/" && current != root && current != path {
		rel, err := filepath.Rel(root, current)
		if err != nil {
			rel = current
		}
		rel = filepath.ToSlash(rel)
		if selected[rel] {
			return rel
		}
		path = current
		current = filepath.Dir(path)
	}
	return ""
}

// getDirectDependencies resolves direct dependencies for a given file path.
// It reads the file content and uses the appropriate resolver.
// Returns a slice of dependency paths relative to the project root, or an error.
func (m *Model) getDirectDependencies(filePath string) ([]string, error) {
	resolver := dependencies.GetResolver(filePath)
	if resolver == nil {
		return nil, nil // No resolver for this file type
	}

	fullPath := filepath.Join(m.rootPath, filePath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file %s for dependency resolution: %w", filePath, err)
	}

	deps, err := resolver.Resolve(content, filePath, m.rootPath, m.projectModuleName)
	if err != nil {
		return nil, fmt.Errorf("error resolving dependencies for %s: %w", filePath, err)
	}

	normalizedDeps := make([]string, 0, len(deps))
	for _, depPath := range deps {
		cleanDepPath := filepath.ToSlash(filepath.Clean(depPath))
		if cleanDepPath != "." && !filepath.IsAbs(cleanDepPath) {
			normalizedDeps = append(normalizedDeps, cleanDepPath)
		}
	}

	return normalizedDeps, nil
}

func (m *Model) toggleSelection(path string, isDir bool) tea.Cmd {
	if path == "" {
		return nil
	}

	fullPath := filepath.Join(m.rootPath, path)
	if _, err := os.Stat(fullPath); err != nil {
		return nil
	}

	maxDepth := 0
	if m.resolveDeps {
		maxDepth = m.maxDepth
	}

	depQueue := []QueuedDep{}
	depProcessed := make(map[string]bool)

	if isDir {
		// Selecting/Deselecting Directory
		if m.selected[path] {
			// Deselecting directory
			delete(m.selected, path)
			delete(m.isDependency, path)

			// Keep track of files deselected by this action
			filesDeselectedInDir := []string{}

			// Deselect children implicitly
			for _, f := range m.files {
				if strings.HasPrefix(f.Path, path+"/") {
					if m.selected[f.Path] {
						if !f.IsDir {
							filesDeselectedInDir = append(filesDeselectedInDir, f.Path)
						}
						delete(m.selected, f.Path)
						delete(m.isDependency, f.Path)
						m.deselected[f.Path] = true
					}
				}
			}
			m.deselected[path] = true

			if m.resolveDeps {
				for _, filePath := range filesDeselectedInDir {
					deps, _ := m.getDirectDependencies(filePath)
					if deps != nil {
						for _, depPath := range deps {
							if m.isDependency[depPath] {
								if parentDir := findParentDirectory(filepath.Join(m.rootPath, depPath), m.rootPath, m.selected); parentDir == "" {
									delete(m.selected, depPath)
									delete(m.isDependency, depPath)
									m.deselected[depPath] = true
								}
							}
						}
					}
				}
			}
		} else {
			// Selecting directory
			m.selected[path] = true
			delete(m.deselected, path)

			searchResultPaths := make(map[string]bool)
			useSearchResults := m.isSearching && len(m.searchResults) > 0
			if useSearchResults {
				for _, node := range m.searchResults {
					if !node.IsDir {
						searchResultPaths[node.Path] = true
					}
				}
			}

			for _, f := range m.files {
				if strings.HasPrefix(f.Path, path+"/") {
					if useSearchResults && !f.IsDir && !searchResultPaths[f.Path] {
						continue
					}

					fFullPath := filepath.Join(m.rootPath, f.Path)
					if (m.useGitIgnore && m.gitIgnoreMgr.IsIgnored(fFullPath)) ||
						(!m.showHidden && utils.IsHiddenPath(f.Path)) ||
						!m.filterMgr.ShouldInclude(f.Path) {
						continue
					}

					newlySelected := !m.selected[f.Path]
					m.selected[f.Path] = true
					delete(m.deselected, f.Path)
					delete(m.isDependency, f.Path)

					if newlySelected && !f.IsDir && m.resolveDeps && maxDepth > 0 {
						if !depProcessed[f.Path] {
							depQueue = append(depQueue, QueuedDep{Path: f.Path, Depth: 0})
							depProcessed[f.Path] = true
						}
					}
				}
			}
		}
	} else {
		// Selecting/Deselecting a File
		if m.selected[path] {
			// Deselecting file
			delete(m.selected, path)
			wasDependency := m.isDependency[path]
			delete(m.isDependency, path)

			if parentDir := findParentDirectory(fullPath, m.rootPath, m.selected); parentDir != "" {
				m.deselected[path] = true
			} else {
				delete(m.deselected, path)
			}

			if m.resolveDeps && wasDependency {
				deps, _ := m.getDirectDependencies(path)
				if deps != nil {
					for _, depPath := range deps {
						if m.isDependency[depPath] {
							if parentDir := findParentDirectory(filepath.Join(m.rootPath, depPath), m.rootPath, m.selected); parentDir == "" {
								delete(m.selected, depPath)
								delete(m.isDependency, depPath)
								m.deselected[depPath] = true
							}
						}
					}
				}
			}
		} else {
			// Selecting file
			m.selected[path] = true
			delete(m.deselected, path)
			delete(m.isDependency, path)

			if m.resolveDeps && maxDepth > 0 {
				if !depProcessed[path] {
					depQueue = append(depQueue, QueuedDep{Path: path, Depth: 0})
					depProcessed[path] = true
				}
			}
		}
	}

	// This runs only if m.resolveDeps is true and items were added to depQueue
	i := 0
	for i < len(depQueue) {
		currentItem := depQueue[i]
		i++

		filePath := currentItem.Path
		currentDepth := currentItem.Depth

		if currentDepth >= maxDepth {
			continue
		}

		directDeps, err := m.getDirectDependencies(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
			continue
		}

		for _, depPath := range directDeps {
			depFullPath := filepath.Join(m.rootPath, depPath)
			info, statErr := os.Stat(depFullPath)
			if statErr != nil || info.IsDir() ||
				(m.useGitIgnore && m.gitIgnoreMgr.IsIgnored(depFullPath)) ||
				(!m.showHidden && utils.IsHiddenPath(depPath)) ||
				!m.filterMgr.ShouldInclude(depPath) {
				continue
			}

			if !m.selected[depPath] {
				m.selected[depPath] = true
				m.isDependency[depPath] = true
				delete(m.deselected, depPath)

				if !depProcessed[depPath] {
					depQueue = append(depQueue, QueuedDep{Path: depPath, Depth: currentDepth + 1})
					depProcessed[depPath] = true
				}
			} else {
				delete(m.deselected, depPath)
			}
		}
	}

	return nil
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
