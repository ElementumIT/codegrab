package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/epilande/codegrab/internal/utils"
)

// FileItem represents a file or directory found when walking the filesystem.
type FileItem struct {
	Path  string
	IsDir bool
	Level int
}

// WalkDirectory traverses the root directory taking into account gitignore and hidden files.
func WalkDirectory(root string, gitIgnore *GitIgnoreManager, filter *FilterManager, useGitIgnore, showHidden bool) ([]FileItem, error) {
	var files []FileItem

	if _, err := os.Stat(root); err != nil {
		return nil, fmt.Errorf("failed to access root directory: %w", err)
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: skipping %s: %v\n", path, err)
			return nil
		}
		if path == root {
			return nil
		}

		if !showHidden && strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if useGitIgnore && gitIgnore != nil && gitIgnore.IsIgnored(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %s: %w", path, err)
		}
		if relPath == "." {
			return nil
		}
		relPath = filepath.ToSlash(relPath)

		// Always include directories when using include patterns
		// This allows us to traverse into directories that might contain matching files
		if info.IsDir() {
			files = append(files, FileItem{
				Path:  relPath,
				IsDir: true,
				Level: strings.Count(relPath, "/"),
			})
			return nil
		}

		if !filter.ShouldInclude(relPath) {
			return nil
		}

		if ok, err := utils.IsTextFile(path); err != nil || !ok {
			return nil
		}

		files = append(files, FileItem{
			Path:  relPath,
			IsDir: false,
			Level: strings.Count(relPath, "/"),
		})
		return nil
	})
	if err != nil {
		return files, fmt.Errorf("error walking directory: %w", err)
	}

	return files, nil
}
