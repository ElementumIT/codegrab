package filesystem

import (
	"testing"
)

func TestFilterManager(t *testing.T) {
	t.Run("ShouldInclude with no patterns", func(t *testing.T) {
		fm := NewFilterManager()
		if !fm.ShouldInclude("file.txt") {
			t.Errorf("Expected file.txt to be included with no patterns")
		}
	})

	t.Run("ShouldInclude with positive pattern", func(t *testing.T) {
		fm := NewFilterManager()
		fm.AddGlobPattern("*.go")
		if !fm.ShouldInclude("main.go") {
			t.Errorf("Expected main.go to be included")
		}
		if fm.ShouldInclude("main.js") {
			t.Errorf("Expected main.js to be excluded")
		}
	})

	t.Run("ShouldInclude with negative pattern", func(t *testing.T) {
		fm := NewFilterManager()
		fm.AddGlobPattern("!*.go")
		if fm.ShouldInclude("main.go") {
			t.Errorf("Expected main.go to be excluded")
		}
		if !fm.ShouldInclude("main.js") {
			t.Errorf("Expected main.js to be included")
		}
	})

	t.Run("ShouldInclude with mixed patterns", func(t *testing.T) {
		fm := NewFilterManager()
		fm.AddGlobPattern("*.{js,ts}")
		fm.AddGlobPattern("!*.spec.js")
		if !fm.ShouldInclude("main.js") {
			t.Errorf("Expected main.js to be included")
		}
		if !fm.ShouldInclude("main.ts") {
			t.Errorf("Expected main.ts to be included")
		}
		if fm.ShouldInclude("main.spec.js") {
			t.Errorf("Expected main.spec.js to be excluded")
		}
		if !fm.ShouldInclude("main.spec.ts") {
			t.Errorf("Expected main.spec.ts to be included (not excluded)")
		}
	})
}

func TestMatchesPattern(t *testing.T) {
	testCases := []struct {
		pattern  string
		path     string
		base     string
		expected bool
	}{
		{"*.go", "main.go", "main.go", true},
		{"*.go", "src/main.go", "main.go", true},
		{"src/*.go", "src/main.go", "main.go", true},
		{"src/*.go", "pkg/main.go", "main.go", false},
		{"*.{js,ts}", "main.js", "main.js", true},
		{"*.{js,ts}", "main.ts", "main.ts", true},
		{"*.{js,ts}", "main.go", "main.go", false},
	}

	for _, tc := range testCases {
		t.Run(tc.pattern+"_"+tc.path, func(t *testing.T) {
			result := matchesPattern(tc.pattern, tc.path, tc.base)
			if result != tc.expected {
				t.Errorf("matchesPattern(%q, %q, %q) = %v, want %v",
					tc.pattern, tc.path, tc.base, result, tc.expected)
			}
		})
	}
}

func TestFilterManagerEdgeCases(t *testing.T) {
	t.Run("Multiple positive patterns", func(t *testing.T) {
		fm := NewFilterManager()
		fm.AddGlobPattern("*.go")
		fm.AddGlobPattern("*.md")

		if !fm.ShouldInclude("main.go") {
			t.Errorf("Expected main.go to be included")
		}
		if !fm.ShouldInclude("README.md") {
			t.Errorf("Expected README.md to be included")
		}
		if fm.ShouldInclude("styles.css") {
			t.Errorf("Expected styles.css to be excluded")
		}
	})

	t.Run("Negative pattern with no positive patterns", func(t *testing.T) {
		fm := NewFilterManager()
		fm.AddGlobPattern("!*.log")

		if !fm.ShouldInclude("main.go") {
			t.Errorf("Expected main.go to be included")
		}
		if fm.ShouldInclude("error.log") {
			t.Errorf("Expected error.log to be excluded")
		}
	})

	t.Run("Complex brace expansion", func(t *testing.T) {
		fm := NewFilterManager()
		fm.AddGlobPattern("*.{js,ts,jsx,tsx}")

		testCases := []struct {
			file     string
			expected bool
		}{
			{"app.js", true},
			{"component.jsx", true},
			{"service.ts", true},
			{"hook.tsx", true},
			{"styles.css", false},
			{"README.md", false},
		}

		for _, tc := range testCases {
			result := fm.ShouldInclude(tc.file)
			if result != tc.expected {
				t.Errorf("ShouldInclude(%q) = %v, want %v", tc.file, result, tc.expected)
			}
		}
	})
}
