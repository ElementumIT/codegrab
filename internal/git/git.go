package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// GitURLPatterns contains regex patterns to detect Git URLs
var GitURLPatterns = []*regexp.Regexp{
	regexp.MustCompile(`^https?://.*\.git/?$`),                    // https://github.com/user/repo.git
	regexp.MustCompile(`^https?://github\.com/[^/]+/[^/]+/?$`),    // https://github.com/user/repo
	regexp.MustCompile(`^https?://gitlab\.com/[^/]+/[^/]+/?$`),    // https://gitlab.com/user/repo
	regexp.MustCompile(`^https?://bitbucket\.org/[^/]+/[^/]+/?$`), // https://bitbucket.org/user/repo
	regexp.MustCompile(`^git@.*:.*\.git$`),                        // git@github.com:user/repo.git
	regexp.MustCompile(`^ssh://git@.*/.*/.*\.git$`),               // ssh://git@github.com/user/repo.git
}

// IsGitURL checks if the given string appears to be a Git repository URL
func IsGitURL(url string) bool {
	url = strings.TrimSpace(url)
	for _, pattern := range GitURLPatterns {
		if pattern.MatchString(url) {
			return true
		}
	}
	return false
}

// CloneRepository clones a Git repository to a temporary directory
// Returns the path to the cloned directory and a cleanup function
func CloneRepository(url string) (string, func(), error) {
	if !IsGitURL(url) {
		return "", nil, fmt.Errorf("invalid Git URL: %s", url)
	}

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "codegrab-clone-*")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	// Clone the repository with shallow depth
	cmd := exec.Command("git", "clone", "--depth=1", url, tempDir)
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("failed to clone repository %s: %w", url, err)
	}

	// Verify the cloned directory exists and is valid
	if stat, err := os.Stat(tempDir); err != nil || !stat.IsDir() {
		cleanup()
		return "", nil, fmt.Errorf("cloned directory is invalid: %s", tempDir)
	}

	return tempDir, cleanup, nil
}

// GetRepoName extracts a repository name from a Git URL for display purposes
func GetRepoName(url string) string {
	url = strings.TrimSpace(url)
	url = strings.TrimSuffix(url, "/")
	url = strings.TrimSuffix(url, ".git")

	// Handle different URL formats
	if strings.Contains(url, "/") {
		parts := strings.Split(url, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
	}

	// Handle git@ format
	if strings.Contains(url, ":") && !strings.Contains(url, "://") {
		parts := strings.Split(url, ":")
		if len(parts) > 1 {
			pathParts := strings.Split(parts[1], "/")
			if len(pathParts) > 0 {
				return pathParts[len(pathParts)-1]
			}
		}
	}

	return filepath.Base(url)
}

