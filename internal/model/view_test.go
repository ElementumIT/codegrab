package model

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/epilande/codegrab/internal/filesystem"
	"github.com/epilande/codegrab/internal/generator"
	"github.com/epilande/codegrab/internal/generator/formats"
	"github.com/epilande/codegrab/internal/ui/themes"
)

func init() {
	// Initialize themes for testing
	themes.Initialize()
}

func TestGetTotalFileCount(t *testing.T) {
	m := Model{
		files: []filesystem.FileItem{
			{Path: "dir1", IsDir: true},
			{Path: "dir1/file1.txt", IsDir: false},
			{Path: "dir1/file2.txt", IsDir: false},
			{Path: "dir2", IsDir: true},
			{Path: "dir2/file3.txt", IsDir: false},
			{Path: "file4.txt", IsDir: false},
		},
	}

	count := m.getTotalFileCount()
	if count != 4 {
		t.Errorf("Expected total file count to be 4, got %d", count)
	}
}

func TestGetSelectedFileCount(t *testing.T) {
	m := Model{
		files: []filesystem.FileItem{
			{Path: "dir1", IsDir: true},
			{Path: "dir1/file1.txt", IsDir: false},
			{Path: "dir1/file2.txt", IsDir: false},
			{Path: "dir2", IsDir: true},
			{Path: "dir2/file3.txt", IsDir: false},
			{Path: "file4.txt", IsDir: false},
		},
		selected:   make(map[string]bool),
		deselected: make(map[string]bool),
	}

	// Test with no selections
	count := m.getSelectedFileCount()
	if count != 0 {
		t.Errorf("Expected selected file count to be 0 with no selections, got %d", count)
	}

	// Test with individual file selections
	m.selected["file4.txt"] = true
	m.selected["dir1/file1.txt"] = true
	count = m.getSelectedFileCount()
	if count != 2 {
		t.Errorf("Expected selected file count to be 2 with individual selections, got %d", count)
	}

	// Test with directory selection
	m.selected = make(map[string]bool)
	m.selected["dir1"] = true
	count = m.getSelectedFileCount()
	if count != 2 {
		t.Errorf("Expected selected file count to be 2 with dir1 selected, got %d", count)
	}

	// Test with deselected files
	m.deselected["dir1/file1.txt"] = true
	count = m.getSelectedFileCount()
	if count != 1 {
		t.Errorf("Expected selected file count to be 1 with one file deselected, got %d", count)
	}

	// Test with search mode
	m.isSearching = true
	m.searchResults = []FileNode{
		{Path: "dir1", Name: "dir1", IsDir: true},
		{Path: "dir1/file2.txt", Name: "file2.txt", IsDir: false},
	}
	count = m.getSelectedFileCount()
	if count != 1 {
		t.Errorf("Expected selected file count to be 1 in search mode, got %d", count)
	}
}

func TestEnsureCursorVisible(t *testing.T) {
	m := Model{
		viewport: viewport.Model{
			Height:  10,
			YOffset: 5,
		},
	}

	// Test cursor above viewport
	m.cursor = 3
	m.ensureCursorVisible()
	if m.viewport.YOffset != 3 {
		t.Errorf("Expected YOffset to be 3 when cursor is above viewport, got %d", m.viewport.YOffset)
	}

	// Test cursor within viewport
	m.viewport.YOffset = 5
	m.cursor = 8
	m.ensureCursorVisible()
	if m.viewport.YOffset != 5 {
		t.Errorf("Expected YOffset to remain 5 when cursor is within viewport, got %d", m.viewport.YOffset)
	}

	// Test cursor below viewport
	m.viewport.YOffset = 5
	m.cursor = 20
	m.ensureCursorVisible()
	if m.viewport.YOffset != 11 {
		t.Errorf("Expected YOffset to be 11 when cursor is below viewport, got %d", m.viewport.YOffset)
	}
}

func TestRefreshViewportContent(t *testing.T) {
	m := Model{
		viewport: viewport.Model{
			Height: 24, // Set viewport height to ensure content renders
			Width:  80,
		},
		selected:   make(map[string]bool),
		deselected: make(map[string]bool),
		displayNodes: []FileNode{
			{Path: "dir1", Name: "dir1", IsDir: true, Level: 0, IsLast: false},
			{Path: "dir1/file1.txt", Name: "file1.txt", IsDir: false, Level: 1, IsLast: true},
			{Path: "file2.txt", Name: "file2.txt", IsDir: false, Level: 0, IsLast: true},
		},
		collapsed: make(map[string]bool),
		cursor:    1,  // Cursor on dir1/file1.txt
		width:     80, // Ensure width matches viewport
	}

	// Mark dir1 as selected
	m.selected["dir1"] = true

	// Refresh viewport content
	m.refreshViewportContent()

	// Check that content was set
	content := m.viewport.View()
	if content == "" {
		t.Errorf("Expected viewport content to be non-empty")
	}

	// Split content into lines
	lines := strings.Split(content, "\n")
	if len(lines) < 3 { // Expect 3 nodes
		t.Fatalf("Expected 3 lines in viewport content, got %d", len(lines))
	}

	// Check directory line rendering
	if !strings.Contains(lines[0], "├──") {
		t.Errorf("Expected directory line to contain branch character, got: %s", lines[0])
	}

	// Check cursor highlighting
	if !strings.Contains(lines[1], "❯") {
		t.Errorf("Expected cursor line to contain highlight marker '❯', got: %s", lines[1])
	}

	// Check checkbox state
	if !strings.Contains(lines[0], "[x]") {
		t.Errorf("Expected selected directory checkbox, got: %s", lines[0])
	}
}

func TestView(t *testing.T) {
	// Create a minimal model for testing View()
	gen := generator.NewGenerator(".", nil, nil, "", false)
	format := formats.GetFormat("markdown")
	gen.SetFormat(format)

	m := Model{
		viewport:    viewport.Model{},
		generator:   gen,
		searchInput: textinput.Model{},
		width:       80,
		height:      24,
	}

	// Test normal view
	view := m.View()
	if !strings.Contains(view, "Code Grab") {
		t.Errorf("Expected view to contain 'Code Grab' header")
	}
	if !strings.Contains(view, "Press '?' for help") {
		t.Errorf("Expected view to contain help text")
	}

	// Test help view
	m.showHelp = true
	m.showHelpScreen()
	view = m.View()
	if !strings.Contains(view, "Help Menu") {
		t.Errorf("Expected help view to contain 'Help Menu' header")
	}
	if !strings.Contains(view, "Exit: esc") {
		t.Errorf("Expected help view to contain 'Exit: esc' footer")
	}

	// Test search view
	m.showHelp = false
	m.isSearching = true
	view = m.View()
	if !strings.Contains(view, "Next: ctrl+n") {
		t.Errorf("Expected search view to contain search help text")
	}

	// Test success message
	m.isSearching = false
	m.successMsg = "Test success message"
	view = m.View()
	if !strings.Contains(view, "Test success message") {
		t.Errorf("Expected view to contain success message")
	}

	// Test error message
	m.successMsg = ""
	m.err = fmt.Errorf("Test error message")
	view = m.View()
	if !strings.Contains(view, "Test error message") {
		t.Errorf("Expected view to contain error message")
	}
}
