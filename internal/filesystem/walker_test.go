package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWalkDirectory(t *testing.T) {
	t.Run("Basic directory walking", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "walk-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		testFiles := []string{
			"file1.txt",
			"file2.go",
			"subdir/file3.txt",
			"subdir/file4.go",
			".hidden/file5.txt",
		}

		for _, file := range testFiles {
			path := filepath.Join(tempDir, filepath.FromSlash(file))
			dir := filepath.Dir(path)
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatalf("Failed to create directory %s: %v", dir, err)
			}
			if err := os.WriteFile(path, []byte("test content"), 0644); err != nil {
				t.Fatalf("Failed to create file %s: %v", path, err)
			}
		}

		filter := NewFilterManager()
		filter.AddGlobPattern("*.go")

		gitIgnore, err := NewGitIgnoreManager(tempDir)
		if err != nil {
			t.Fatalf("Failed to create GitIgnoreManager: %v", err)
		}

		testCases := []struct {
			name         string
			checkFile    string
			expectedLen  int
			useGitIgnore bool
			showHidden   bool
			shouldExist  bool
		}{
			{
				name:         "Filter Go files only",
				useGitIgnore: false,
				showHidden:   false,
				expectedLen:  3,
				checkFile:    "file2.go",
				shouldExist:  true,
			},
			{
				name:         "Show hidden files",
				useGitIgnore: false,
				showHidden:   true,
				expectedLen:  4,
				checkFile:    ".hidden",
				shouldExist:  true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				files, err := WalkDirectory(tempDir, gitIgnore, filter, tc.useGitIgnore, tc.showHidden)
				if err != nil {
					t.Fatalf("WalkDirectory failed: %v", err)
				}

				if len(files) != tc.expectedLen {
					t.Errorf("Expected %d files, got %d", tc.expectedLen, len(files))
				}

				found := false
				for _, file := range files {
					if file.Path == tc.checkFile {
						found = true
						break
					}
				}

				if found != tc.shouldExist {
					if tc.shouldExist {
						t.Errorf("Expected to find %s but didn't", tc.checkFile)
					} else {
						t.Errorf("Expected not to find %s but did", tc.checkFile)
					}
				}
			})
		}
	})

	t.Run("Non-existent root directory", func(t *testing.T) {
		filter := NewFilterManager()
		gitIgnore, _ := NewGitIgnoreManager(".")

		_, err := WalkDirectory("/path/that/does/not/exist", gitIgnore, filter, false, false)
		if err == nil {
			t.Error("Expected error for non-existent directory, got nil")
		}
	})
}
