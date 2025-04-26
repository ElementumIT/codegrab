package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// NormalizeGlobPattern takes a glob pattern and the absolute project root path,
// then returns a normalized pattern for matching against relative file paths,
// and a boolean indicating if the pattern is valid and should be used.
func NormalizeGlobPattern(pattern, absRoot string) (string, bool) {
	originalInputPattern := pattern
	isNegative := strings.HasPrefix(pattern, "!") || strings.HasPrefix(pattern, "\\!")

	if strings.HasPrefix(pattern, "\\!") {
		pattern = "!" + strings.TrimPrefix(pattern, "\\!")
	}

	pathPart := pattern
	if isNegative {
		pathPart = pattern[1:]
	}

	var normalizedPath string

	if filepath.IsAbs(pathPart) {
		relPattern, err := filepath.Rel(absRoot, pathPart)
		if err == nil && relPattern != "." && !strings.HasPrefix(relPattern, ".."+string(filepath.Separator)) && relPattern != ".." {
			normalizedPath = relPattern
		} else {
			fmt.Fprintf(os.Stderr, "Warning: Absolute glob pattern %q points outside the target directory %q or is invalid. Ignoring.\n", originalInputPattern, absRoot)
			return "", false
		}
	} else if strings.HasPrefix(pathPart, "./") {
		normalizedPath = strings.TrimPrefix(pathPart, "./")
	} else {
		normalizedPath = pathPart
	}

	normalizedPath = strings.ReplaceAll(normalizedPath, "\\", "/")

	if normalizedPath == "" {
		if originalInputPattern != "./" && originalInputPattern != "!./" && originalInputPattern != "\\!./" {
			fmt.Fprintf(os.Stderr, "Warning: Glob pattern %q resulted in an empty path after normalization. Ignoring.\n", originalInputPattern)
		}
		return "", false
	}

	finalNormalizedPattern := normalizedPath
	if isNegative {
		finalNormalizedPattern = "!" + normalizedPath
	}

	return finalNormalizedPattern, true
}
