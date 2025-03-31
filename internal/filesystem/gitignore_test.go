package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGitIgnoreManager(t *testing.T) {
	t.Run("New with no gitignore file", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "gitignore-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		manager, err := NewGitIgnoreManager(tempDir)
		if err != nil {
			t.Fatalf("Failed to create GitIgnoreManager: %v", err)
		}

		if manager.IsIgnored(filepath.Join(tempDir, "file.txt")) {
			t.Error("Expected file.txt not to be ignored when no gitignore exists")
		}
	})

	t.Run("New with gitignore file", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "gitignore-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		gitignoreContent := "*.log\nnode_modules/\n# Comment line\n"
		err = os.WriteFile(filepath.Join(tempDir, ".gitignore"), []byte(gitignoreContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create .gitignore file: %v", err)
		}

		manager, err := NewGitIgnoreManager(tempDir)
		if err != nil {
			t.Fatalf("Failed to create GitIgnoreManager: %v", err)
		}

		testCases := []struct {
			path     string
			expected bool
		}{
			{filepath.Join(tempDir, "app.log"), true},
			{filepath.Join(tempDir, "node_modules", "package.json"), true},
			{filepath.Join(tempDir, "src", "main.go"), false},
		}

		for _, tc := range testCases {
			t.Run(tc.path, func(t *testing.T) {
				if manager.IsIgnored(tc.path) != tc.expected {
					t.Errorf("IsIgnored(%q) = %v, want %v", tc.path, !tc.expected, tc.expected)
				}
			})
		}
	})

	t.Run("normalizeLine function", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected string
		}{
			{"node_modules/", "node_modules"},
			{"  *.log  ", "*.log"},
			{"dist/", "dist"},
			{"src/temp", "src/temp"},
		}

		for _, tc := range testCases {
			result := normalizeLine(tc.input)
			if result != tc.expected {
				t.Errorf("normalizeLine(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		}
	})
}
