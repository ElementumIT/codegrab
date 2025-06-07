package git

import (
	"os"
	"testing"
)

func TestIsGitURL(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		// Valid Git URLs
		{"https://github.com/user/repo.git", true},
		{"https://github.com/user/repo", true},
		{"https://gitlab.com/user/repo", true},
		{"https://bitbucket.org/user/repo", true},
		{"git@github.com:user/repo.git", true},
		{"ssh://git@github.com/user/repo.git", true},
		{"http://github.com/user/repo.git", true},

		// Invalid URLs
		{"./local/path", false},
		{"/absolute/path", false},
		{"../relative/path", false},
		{"not-a-url", false},
		{"https://example.com", false},
		{"https://github.com", false},
		{"", false},
		{"   ", false},
	}

	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			result := IsGitURL(test.url)
			if result != test.expected {
				t.Errorf("IsGitURL(%q) = %v, expected %v", test.url, result, test.expected)
			}
		})
	}
}

func TestGetRepoName(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{"https://github.com/user/repo.git", "repo"},
		{"https://github.com/user/repo", "repo"},
		{"https://gitlab.com/user/my-project", "my-project"},
		{"git@github.com:user/repo.git", "repo"},
		{"ssh://git@github.com/user/repo.git", "repo"},
		{"https://github.com/user/repo/", "repo"},
	}

	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			result := GetRepoName(test.url)
			if result != test.expected {
				t.Errorf("GetRepoName(%q) = %q, expected %q", test.url, result, test.expected)
			}
		})
	}
}

func TestCloneRepository_InvalidURL(t *testing.T) {
	_, cleanup, err := CloneRepository("not-a-git-url")
	if err == nil {
		if cleanup != nil {
			cleanup()
		}
		t.Error("CloneRepository should fail for invalid URL")
	}
}

// Note: We skip actual cloning tests as they require network access
// and would make tests slow and potentially flaky. In a real scenario,
// you might want integration tests with a test repository.
func TestCloneRepository_NetworkTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping network test in short mode")
	}

	// Only run if we have network access and git is available
	if _, err := os.Stat("/usr/bin/git"); err != nil {
		if _, err := os.Stat("/usr/local/bin/git"); err != nil {
			t.Skip("git not found, skipping clone test")
		}
	}

	// Test with a small public repository
	tempDir, cleanup, err := CloneRepository("https://github.com/octocat/Hello-World.git")
	if err != nil {
		t.Skipf("Could not clone test repository (network issue?): %v", err)
		return
	}
	defer cleanup()

	// Verify the directory was created
	if stat, err := os.Stat(tempDir); err != nil || !stat.IsDir() {
		t.Errorf("Cloned directory %s is not valid", tempDir)
	}

	// Verify cleanup works
	cleanup()
	if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
		t.Error("Cleanup did not remove temporary directory")
	}
}

