package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

// GitIgnoreManager manages gitignore patterns for a repository.
type GitIgnoreManager struct {
	ignore *ignore.GitIgnore
	root   string
}

// normalizeLine trims whitespace and trailing slashes for a gitignore pattern.
func normalizeLine(line string) string {
	line = strings.TrimSpace(line)
	return strings.TrimSuffix(line, "/")
}

// NewGitIgnoreManager reads the .gitignore file (if present) and returns a manager.
func NewGitIgnoreManager(root string) (*GitIgnoreManager, error) {
	gitignorePath := filepath.Join(root, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		return &GitIgnoreManager{
			ignore: ignore.CompileIgnoreLines(),
			root:   root,
		}, nil
	}

	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read gitignore file: %w", err)
	}

	var patterns []string
	for _, line := range strings.Split(string(content), "\n") {
		line = normalizeLine(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}

	return &GitIgnoreManager{
		ignore: ignore.CompileIgnoreLines(patterns...),
		root:   root,
	}, nil
}

// IsIgnored returns true if the provided path is ignored.
func (g *GitIgnoreManager) IsIgnored(path string) bool {
	relPath, err := filepath.Rel(g.root, path)
	if err != nil {
		return false
	}
	relPath = filepath.ToSlash(relPath)
	relPath = strings.TrimSuffix(relPath, "/")
	return g.ignore.MatchesPath(relPath)
}
