package model

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
)

func TestFuzzyMatch(t *testing.T) {
	testCases := []struct {
		query    string
		target   string
		expected bool
	}{
		{"", "anything", true},           // Empty query matches anything
		{"abc", "abc", true},             // Exact match
		{"abc", "abcdef", true},          // Prefix match
		{"bcd", "abcdef", true},          // Substring match
		{"adf", "abcdef", true},          // Non-consecutive letters match
		{"xyz", "abcdef", false},         // No match
		{"ABC", "abc", true},             // Case insensitive
		{"abc", "ABC", true},             // Case insensitive
		{"md", "main.go", false},         // No match
		{"mg", "main.go", true},          // First letters of words
		{"mgo", "main.go", true},         // First letters plus more
		{"test", "test_file.go", true},   // Prefix match
		{"tst", "test_file.go", true},    // Skip letters
		{"tfile", "test_file.go", true},  // Skip underscore
		{"tfil", "test_file.go", true},   // Skip underscore and letters
		{"tstfil", "test_file.go", true}, // Skip underscore and more letters
		{" main ", "main.go", true},      // Leading/trailing spaces
		{"Äc", "aÄc.go", true},           // Unicode characters
	}

	for _, tc := range testCases {
		t.Run(tc.query+"_"+tc.target, func(t *testing.T) {
			result := fuzzyMatch(tc.query, tc.target)
			if result != tc.expected {
				t.Errorf("fuzzyMatch(%q, %q) = %v, want %v", tc.query, tc.target, result, tc.expected)
			}
		})
	}
}

func TestUpdateSearchResults(t *testing.T) {
	m := Model{
		searchInput: textinput.Model{},
		displayNodes: []FileNode{
			{Path: "dir1", Name: "dir1", IsDir: true, Level: 0},
			{Path: "dir1/file1.txt", Name: "file1.txt", IsDir: false, Level: 1},
			{Path: "dir1/test_file.go", Name: "test_file.go", IsDir: false, Level: 1},
			{Path: "dir2", Name: "dir2", IsDir: true, Level: 0},
			{Path: "dir2/main.go", Name: "main.go", IsDir: false, Level: 1},
			{Path: "file3.txt", Name: "file3.txt", IsDir: false, Level: 0},
		},
		collapsed: make(map[string]bool),
	}

	// Test with empty query
	m.searchInput.SetValue("")
	m.updateSearchResults()
	if len(m.searchResults) != 0 {
		t.Errorf("Expected 0 search results with empty query, got %d", len(m.searchResults))
	}

	// Test with query that matches a file
	m.searchInput.SetValue("main")
	m.updateSearchResults()
	if len(m.searchResults) != 2 {
		t.Fatalf("Expected 2 search results with 'main' query, got %d", len(m.searchResults))
	}
	// Should include dir2 (parent) and main.go
	foundDir2 := false
	foundMainGo := false
	for _, node := range m.searchResults {
		if node.Path == "dir2" {
			foundDir2 = true
		}
		if node.Path == "dir2/main.go" {
			foundMainGo = true
		}
	}
	if !foundDir2 {
		t.Errorf("Expected dir2 to be in search results")
	}
	if !foundMainGo {
		t.Errorf("Expected dir2/main.go to be in search results")
	}

	m.searchInput.SetValue("file")
	m.updateSearchResults()
	if len(m.searchResults) != 4 {
		t.Fatalf("Expected 4 search results with 'file' query, got %d", len(m.searchResults))
	}
	// Should include dir1, dir1/file1.txt, dir1/test_file.go, and file3.txt
	foundDir1 := false
	foundFile1 := false
	foundTestFile := false
	foundFile3 := false
	for _, node := range m.searchResults {
		if node.Path == "dir1" {
			foundDir1 = true
		}
		if node.Path == "dir1/file1.txt" {
			foundFile1 = true
		}
		if node.Path == "dir1/test_file.go" {
			foundTestFile = true
		}
		if node.Path == "file3.txt" {
			foundFile3 = true
		}
	}
	if !foundDir1 {
		t.Errorf("Expected dir1 to be in search results")
	}
	if !foundFile1 {
		t.Errorf("Expected dir1/file1.txt to be in search results")
	}
	if !foundTestFile {
		t.Errorf("Expected dir1/test_file.go to be in search results")
	}
	if !foundFile3 {
		t.Errorf("Expected file3.txt to be in search results")
	}

	// Test with fuzzy query
	m.searchInput.SetValue("tgo")
	m.updateSearchResults()
	if len(m.searchResults) != 2 {
		t.Fatalf("Expected 2 search results with 'tgo' query, got %d", len(m.searchResults))
	}
	// Should include dir1 (parent) and test_file.go
	foundDir1 = false
	foundTestFile = false
	for _, node := range m.searchResults {
		if node.Path == "dir1" {
			foundDir1 = true
		}
		if node.Path == "dir1/test_file.go" {
			foundTestFile = true
		}
	}
	if !foundDir1 {
		t.Errorf("Expected dir1 to be in search results")
	}
	if !foundTestFile {
		t.Errorf("Expected dir1/test_file.go to be in search results")
	}
}

func TestIsInSearchResults(t *testing.T) {
	m := Model{
		isSearching: true,
		searchResults: []FileNode{
			{Path: "dir1", Name: "dir1", IsDir: true},
			{Path: "dir1/file1.txt", Name: "file1.txt", IsDir: false},
			{Path: "file2.txt", Name: "file2.txt", IsDir: false},
		},
	}

	// Test with a file that is in search results
	if !m.isInSearchResults("dir1/file1.txt") {
		t.Errorf("Expected dir1/file1.txt to be in search results")
	}

	// Test with a file that is not in search results
	if m.isInSearchResults("dir1/file3.txt") {
		t.Errorf("Expected dir1/file3.txt not to be in search results")
	}

	// Test with a directory (should return false since we only check files)
	if m.isInSearchResults("dir1") {
		t.Errorf("Expected dir1 not to be in search results (directories should be ignored)")
	}

	// Test when search is not active
	m.isSearching = false
	if m.isInSearchResults("dir1/file1.txt") {
		t.Errorf("Expected isInSearchResults to return false when search is not active")
	}

	// Test with empty search results
	m.isSearching = true
	m.searchResults = nil
	if m.isInSearchResults("dir1/file1.txt") {
		t.Errorf("Expected isInSearchResults to return false with empty search results")
	}
}
