package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/epilande/codegrab/internal/cache"
	"github.com/epilande/codegrab/internal/utils"
)

// FileItem represents a file or directory found when walking the filesystem.
type FileItem struct {
	Path  string
	IsDir bool
	Level int
	Size  int64
}

// WalkDirectory traverses the root directory taking into account gitignore, hidden files, and max file size.
func WalkDirectory(root string, gitIgnore *GitIgnoreManager, filter *FilterManager, useGitIgnore, showHidden bool, maxFileSize int64) ([]FileItem, error) {
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

		// Skip hidden directories/files
		if !showHidden && strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip gitignored paths
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
				Size:  info.Size(),
			})
			return nil
		}

		// Skip files larger than maxFileSize
		if info.Size() > maxFileSize {
			return nil
		}

		// Skip files not matching glob patterns
		if !filter.ShouldInclude(relPath) {
			return nil
		}

		fileCache := cache.GetGlobalFileCache()
		if ok, err := fileCache.GetTextFileStatus(path, utils.IsTextFile); err != nil || !ok {
			return nil
		}
		files = append(files, FileItem{
			Path:  relPath,
			IsDir: false,
			Level: strings.Count(relPath, "/"),
			Size:  info.Size(),
		})
		return nil
	})
	if err != nil {
		return files, fmt.Errorf("error walking directory: %w", err)
	}

	return files, nil
}
