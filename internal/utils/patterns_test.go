package utils

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestNormalizeGlobPattern(t *testing.T) {
	testAbsRoot := filepath.FromSlash("/home/user/project")
	testAbsRootWin := `C:\Users\user\project`

	testCases := []struct {
		name            string
		inputPattern    string
		absRoot         string
		expectedPattern string
		expectedValid   bool
	}{
		// --- Basic Relative Patterns ---
		{
			name:            "Simple file pattern",
			inputPattern:    "*.go",
			absRoot:         testAbsRoot,
			expectedPattern: "*.go",
			expectedValid:   true,
		},
		{
			name:            "Relative path pattern",
			inputPattern:    "cmd/grab/*.go",
			absRoot:         testAbsRoot,
			expectedPattern: "cmd/grab/*.go",
			expectedValid:   true,
		},
		{
			name:            "Relative path pattern with Windows separators",
			inputPattern:    `cmd\grab\*.go`,
			absRoot:         testAbsRoot,
			expectedPattern: "cmd/grab/*.go",
			expectedValid:   true,
		},

		// --- Patterns with ./ prefix ---
		{
			name:            "./ prefix file pattern",
			inputPattern:    "./*.go",
			absRoot:         testAbsRoot,
			expectedPattern: "*.go",
			expectedValid:   true,
		},
		{
			name:            "./ prefix relative path",
			inputPattern:    "./cmd/grab/*.go",
			absRoot:         testAbsRoot,
			expectedPattern: "cmd/grab/*.go",
			expectedValid:   true,
		},
		{
			name:            "./ prefix only",
			inputPattern:    "./",
			absRoot:         testAbsRoot,
			expectedPattern: "",
			expectedValid:   false,
		},

		// --- Absolute Paths ---
		{
			name:            "Absolute path within root",
			inputPattern:    filepath.FromSlash("/home/user/project/cmd/main.go"),
			absRoot:         testAbsRoot,
			expectedPattern: "cmd/main.go",
			expectedValid:   true,
		},
		{
			name:            "Absolute dir pattern within root",
			inputPattern:    filepath.FromSlash("/home/user/project/internal/*"),
			absRoot:         testAbsRoot,
			expectedPattern: "internal/*",
			expectedValid:   true,
		},
		{
			name:            "Absolute path matching root",
			inputPattern:    testAbsRoot,
			absRoot:         testAbsRoot,
			expectedPattern: "",
			expectedValid:   false,
		},
		{
			name:            "Absolute path outside root",
			inputPattern:    filepath.FromSlash("/home/user/other/file.txt"),
			absRoot:         testAbsRoot,
			expectedPattern: "",
			expectedValid:   false,
		},
		{
			name:            "Absolute path above root",
			inputPattern:    filepath.FromSlash("/home/user/*.go"),
			absRoot:         testAbsRoot,
			expectedPattern: "",
			expectedValid:   false,
		},
		{
			name:            "Absolute path root level",
			inputPattern:    filepath.FromSlash("/*.go"),
			absRoot:         testAbsRoot,
			expectedPattern: "",
			expectedValid:   false,
		},
		{
			name:            "Absolute path Windows within root",
			inputPattern:    `C:\Users\user\project\cmd\main.go`,
			absRoot:         testAbsRootWin,
			expectedPattern: "cmd/main.go",
			expectedValid:   true,
		},
		{
			name:            "Absolute path Windows outside root",
			inputPattern:    `C:\Users\other\file.txt`,
			absRoot:         testAbsRootWin,
			expectedPattern: "",
			expectedValid:   false,
		},
		{
			name:            "Absolute path Windows different drive",
			inputPattern:    `D:\data\file.txt`,
			absRoot:         testAbsRootWin,
			expectedPattern: "",
			expectedValid:   false,
		},

		// --- Negated Patterns ---
		{
			name:            "Negated simple file pattern",
			inputPattern:    "!*.log",
			absRoot:         testAbsRoot,
			expectedPattern: "!*.log",
			expectedValid:   true,
		},
		{
			name:            "Negated relative path pattern",
			inputPattern:    "!cmd/grab/*.log",
			absRoot:         testAbsRoot,
			expectedPattern: "!cmd/grab/*.log",
			expectedValid:   true,
		},
		{
			name:            "Negated ./ prefix",
			inputPattern:    "!./tmp/*",
			absRoot:         testAbsRoot,
			expectedPattern: "!tmp/*",
			expectedValid:   true,
		},
		{
			name:            "Negated absolute path within root",
			inputPattern:    "!" + filepath.FromSlash("/home/user/project/secrets.txt"),
			absRoot:         testAbsRoot,
			expectedPattern: "!secrets.txt",
			expectedValid:   true,
		},
		{
			name:            "Negated absolute path outside root",
			inputPattern:    "!" + filepath.FromSlash("/home/user/other/config.yml"),
			absRoot:         testAbsRoot,
			expectedPattern: "",
			expectedValid:   false,
		},
		{
			name:            "Negated absolute path above root",
			inputPattern:    "!" + filepath.FromSlash("/home/user/*.bak"),
			absRoot:         testAbsRoot,
			expectedPattern: "",
			expectedValid:   false,
		},
		{
			name:            "Negated ./ only",
			inputPattern:    "!./",
			absRoot:         testAbsRoot,
			expectedPattern: "",
			expectedValid:   false,
		},

		// --- Escaped Negated Patterns ---
		{
			name:            "Escaped negated simple pattern",
			inputPattern:    `\!*.tmp`,
			absRoot:         testAbsRoot,
			expectedPattern: "!*.tmp",
			expectedValid:   true,
		},
		{
			name:            "Escaped negated relative path",
			inputPattern:    `\!build/*`,
			absRoot:         testAbsRoot,
			expectedPattern: "!build/*",
			expectedValid:   true,
		},
		{
			name:            "Escaped negated ./ prefix",
			inputPattern:    `\!./*.swp`,
			absRoot:         testAbsRoot,
			expectedPattern: "!*.swp",
			expectedValid:   true,
		},
		{
			name:            "Escaped negated absolute path within root",
			inputPattern:    `\!` + filepath.FromSlash("/home/user/project/vendor/*"),
			absRoot:         testAbsRoot,
			expectedPattern: "!vendor/*",
			expectedValid:   true,
		},
		{
			name:            "Escaped negated absolute path outside root",
			inputPattern:    `\!` + filepath.FromSlash("/etc/passwd"),
			absRoot:         testAbsRoot,
			expectedPattern: "",
			expectedValid:   false,
		},

		// --- Edge Cases ---
		{
			name:            "Empty input pattern",
			inputPattern:    "",
			absRoot:         testAbsRoot,
			expectedPattern: "",
			expectedValid:   false,
		},
		{
			name:            "Input pattern just !",
			inputPattern:    "!",
			absRoot:         testAbsRoot,
			expectedPattern: "",
			expectedValid:   false,
		},
		{
			name:            "Input pattern just \\!",
			inputPattern:    `\!`,
			absRoot:         testAbsRoot,
			expectedPattern: "",
			expectedValid:   false,
		},
	}

	for _, tc := range testCases {
		if runtime.GOOS != "windows" && tc.absRoot == testAbsRootWin {
			continue
		}

		t.Run(tc.name, func(t *testing.T) {
			normalized, isValid := NormalizeGlobPattern(tc.inputPattern, tc.absRoot)

			if isValid != tc.expectedValid {
				t.Errorf("Expected validity %v, but got %v", tc.expectedValid, isValid)
			}

			if tc.expectedValid && normalized != tc.expectedPattern {
				t.Errorf("Expected normalized pattern %q, but got %q", tc.expectedPattern, normalized)
			}

			if !tc.expectedValid && normalized != "" {
				t.Errorf("Expected empty pattern for invalid input, but got %q", normalized)
			}
		})
	}
}
