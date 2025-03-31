package filesystem

import (
	"path/filepath"
	"strings"
)

type FilterManager struct {
	Patterns []string
}

func NewFilterManager() *FilterManager {
	return &FilterManager{
		Patterns: make([]string, 0),
	}
}

func (f *FilterManager) AddGlobPattern(pattern string) {
	f.Patterns = append(f.Patterns, pattern)
}

func (f *FilterManager) ShouldInclude(path string) bool {
	if len(f.Patterns) == 0 {
		return true
	}

	base := filepath.Base(path)
	hasPositivePattern := false

	// First check if path matches any negative patterns
	for _, pattern := range f.Patterns {
		isNegative := strings.HasPrefix(pattern, "!")
		if !isNegative {
			hasPositivePattern = true
			continue
		}

		// Remove the ! prefix
		pattern = pattern[1:]
		if matchesPattern(pattern, path, base) {
			return false
		}
	}

	// If there are no positive patterns, include everything that wasn't explicitly excluded
	if !hasPositivePattern {
		return true
	}

	// Check positive patterns
	for _, pattern := range f.Patterns {
		if strings.HasPrefix(pattern, "!") {
			continue
		}

		if matchesPattern(pattern, path, base) {
			return true
		}
	}

	return false
}

func matchesPattern(pattern, path, base string) bool {
	// Handle brace expansion
	if strings.Contains(pattern, "{") && strings.Contains(pattern, "}") {
		start := strings.Index(pattern, "{")
		end := strings.Index(pattern, "}")
		if start != -1 && end != -1 && end > start {
			prefix := pattern[:start]
			suffix := pattern[end+1:]
			options := strings.Split(pattern[start+1:end], ",")

			for _, opt := range options {
				opt = strings.TrimSpace(opt)
				expandedPattern := prefix + opt + suffix

				// Try matching against full path
				pathMatched, err := filepath.Match(expandedPattern, path)
				if err == nil && pathMatched {
					return true
				}

				// Try matching against base name
				baseMatched, err := filepath.Match(expandedPattern, base)
				if err == nil && baseMatched {
					return true
				}
			}
			return false
		}
	}

	pathMatched, err := filepath.Match(pattern, path)
	if err == nil && pathMatched {
		return true
	}

	baseMatched, err := filepath.Match(pattern, base)
	if err == nil && baseMatched {
		return true
	}

	return false
}
