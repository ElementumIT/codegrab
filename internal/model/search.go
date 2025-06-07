package model

import (
	"path/filepath"
	"strings"
	"unicode"
)

// fuzzyMatch checks if query fuzzy matches the target string
func fuzzyMatch(query, target string) bool {
	// Normalize query and target strings
	queryRunes := []rune(strings.ToLower(strings.TrimSpace(query)))
	targetRunes := []rune(strings.ToLower(target))

	if len(queryRunes) == 0 {
		return true
	}

	queryIdx := 0
	for _, char := range targetRunes {
		if queryIdx < len(queryRunes) && unicode.ToLower(char) == queryRunes[queryIdx] {
			queryIdx++
		}
	}
	return queryIdx == len(queryRunes)
}

// updateSearchResults filters displayNodes based on search query and preserves folder structure
func (m *Model) updateSearchResults() {
	m.searchResults = nil
	query := strings.TrimSpace(m.searchInput.Value())

	// If query is empty, return without setting search results
	if query == "" {
		return
	}

	// First find all matching file paths
	matchedFiles := make(map[string]bool)
	for _, node := range m.displayNodes {
		if !node.IsDir && fuzzyMatch(query, node.Path) {
			matchedFiles[node.Path] = true
		}
	}

	// Then identify all parent directories of matched files
	matchedPaths := make(map[string]bool)
	for path := range matchedFiles {
		matchedPaths[path] = true

		// Add all parent directories to the matched paths
		dir := filepath.Dir(path)
		for dir != "." && dir != "/" {
			matchedPaths[dir] = true
			dir = filepath.Dir(dir)
		}
	}

	// Then add all nodes that are either matches or parents of matches
	// Keep the original order to maintain hierarchy
	for _, node := range m.displayNodes {
		if matchedPaths[node.Path] {
			// If it's a directory that contains matches, make sure it's expanded
			if node.IsDir {
				m.collapsed[node.Path] = false
			}
			m.searchResults = append(m.searchResults, node)
		}
	}
}

// isInSearchResults checks if a file path is present in the current search results
func (m *Model) isInSearchResults(path string) bool {
	if !m.isSearching || len(m.searchResults) == 0 {
		return false
	}

	for _, node := range m.searchResults {
		if !node.IsDir && node.Path == path {
			return true
		}
	}
	return false
}
