package filesystem

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWalkDirectory(t *testing.T) {
	t.Run("Basic directory walking", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "walk-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		testFilesSetup := map[string]string{
			"small.txt":        "small",                                // 5 bytes
			"medium.txt":       strings.Repeat("a", 500),               // 500 bytes
			"large.txt":        strings.Repeat("b", 1500),              // 1500 bytes
			"subdir/small.go":  "package small",                        // 13 bytes
			"subdir/large.go":  "// " + strings.Repeat("c", 2000),      // ~2004 bytes
			".hidden/huge.txt": strings.Repeat("d", 5000),              // 5000 bytes
			"binary.bin":       string([]byte{0x00, 0x01, 0x02, 0x00}), // Binary file to test IsTextFile
		}

		for file, content := range testFilesSetup {
			path := filepath.Join(tempDir, filepath.FromSlash(file))
			dir := filepath.Dir(path)
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatalf("Failed to create directory %s: %v", dir, err)
			}
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				t.Fatalf("Failed to create file %s: %v", path, err)
			}
		}

		gitIgnore, err := NewGitIgnoreManager(tempDir)
		if err != nil {
			t.Fatalf("Failed to create GitIgnoreManager: %v", err)
		}

		testCases := []struct {
			expectedFiles    map[string]bool
			name             string
			filterPatterns   []string
			maxFileSize      int64
			useGitIgnore     bool
			showHidden       bool
			expectExcludeBin bool
		}{
			{
				name:           "No filters, default max size",
				filterPatterns: []string{},
				useGitIgnore:   false,
				showHidden:     false,
				maxFileSize:    math.MaxInt64,
				expectedFiles: map[string]bool{
					"small.txt":       true,
					"medium.txt":      true,
					"large.txt":       true,
					"subdir":          true,
					"subdir/small.go": true,
					"subdir/large.go": true,
				},
				expectExcludeBin: true,
			},
			{
				name:           "Limit 1kb, include all text",
				filterPatterns: []string{"*"},
				useGitIgnore:   false,
				showHidden:     false,
				maxFileSize:    1024,
				expectedFiles: map[string]bool{
					"small.txt":       true,
					"medium.txt":      true,
					"subdir":          true,
					"subdir/small.go": true,
				},
				expectExcludeBin: true,
			},
			{
				name:           "Limit 1kb, filter *.go",
				filterPatterns: []string{"*.go"},
				useGitIgnore:   false,
				showHidden:     false,
				maxFileSize:    1024,
				expectedFiles: map[string]bool{
					"subdir":          true,
					"subdir/small.go": true,
				},
				expectExcludeBin: true,
			},
			{
				name:           "Show hidden, limit 4kb",
				filterPatterns: []string{"*"},
				useGitIgnore:   false,
				showHidden:     true,
				maxFileSize:    4 * 1024,
				expectedFiles: map[string]bool{
					"small.txt":       true,
					"medium.txt":      true,
					"large.txt":       true,
					"subdir":          true,
					"subdir/small.go": true,
					"subdir/large.go": true,
					".hidden":         true,
				},
				expectExcludeBin: true,
			},
			{
				name:           "No limit, exclude binary",
				filterPatterns: []string{"*"},
				useGitIgnore:   false,
				showHidden:     false,
				maxFileSize:    math.MaxInt64,
				expectedFiles: map[string]bool{
					"small.txt":       true,
					"medium.txt":      true,
					"large.txt":       true,
					"subdir":          true,
					"subdir/small.go": true,
					"subdir/large.go": true,
				},
				expectExcludeBin: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				filter := NewFilterManager()
				for _, p := range tc.filterPatterns {
					filter.AddGlobPattern(p)
				}

				files, err := WalkDirectory(tempDir, gitIgnore, filter, tc.useGitIgnore, tc.showHidden, tc.maxFileSize)
				if err != nil {
					t.Fatalf("WalkDirectory failed: %v", err)
				}

				foundFilesMap := make(map[string]bool)
				var foundPaths []string
				for _, f := range files {
					foundFilesMap[f.Path] = true
					foundPaths = append(foundPaths, f.Path)
				}
				t.Logf("Found files (%d): %v", len(files), foundPaths)

				if len(files) != len(tc.expectedFiles) {
					t.Errorf("Expected %d files, got %d. Expected: %v, Got: %v", len(tc.expectedFiles), len(files), tc.expectedFiles, foundFilesMap)
				}

				for expectedPath := range tc.expectedFiles {
					if !foundFilesMap[expectedPath] {
						t.Errorf("Expected to find '%s' but didn't", expectedPath)
					}
				}
				for foundPath := range foundFilesMap {
					if !tc.expectedFiles[foundPath] {
						t.Errorf("Found unexpected file '%s'", foundPath)
					}
				}

				// Explicitly check binary exclusion if expected
				if tc.expectExcludeBin {
					if foundFilesMap["binary.bin"] {
						t.Errorf("Expected 'binary.bin' to be excluded, but it was found")
					}
				}
			})
		}
	})

	t.Run("Non-existent root directory", func(t *testing.T) {
		filter := NewFilterManager()
		gitIgnore, _ := NewGitIgnoreManager(".")

		_, err := WalkDirectory("/path/that/does/not/exist", gitIgnore, filter, false, false, math.MaxInt64)
		if err == nil {
			t.Error("Expected error for non-existent directory, got nil")
		}
	})
}
