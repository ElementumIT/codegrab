package model

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/epilande/codegrab/internal/filesystem"
)

func TestNewModel(t *testing.T) {
	// Setup
	tempDir, err := os.MkdirTemp("", "model-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filterMgr := filesystem.NewFilterManager()
	config := Config{
		RootPath:    tempDir,
		FilterMgr:   filterMgr,
		OutputPath:  "output.md",
		Format:      "markdown",
		UseTempFile: false,
	}

	// Test
	model := NewModel(config)

	// Assertions
	if model.rootPath != tempDir {
		t.Errorf("Expected rootPath to be %q, got %q", tempDir, model.rootPath)
	}
	if model.filterMgr != filterMgr {
		t.Errorf("Expected filterMgr to be the provided instance")
	}
	if model.generator == nil {
		t.Errorf("Expected generator to be initialized")
	}
	if model.generator.GetFormatName() != "markdown" {
		t.Errorf("Expected format to be markdown, got %s", model.generator.GetFormatName())
	}
	if len(model.selected) != 0 {
		t.Errorf("Expected selected map to be empty, got %d items", len(model.selected))
	}
	if len(model.deselected) != 0 {
		t.Errorf("Expected deselected map to be empty, got %d items", len(model.deselected))
	}
	if len(model.collapsed) != 0 {
		t.Errorf("Expected collapsed map to be empty, got %d items", len(model.collapsed))
	}
	if !model.useGitIgnore {
		t.Errorf("Expected useGitIgnore to be true by default")
	}
	if model.showHidden {
		t.Errorf("Expected showHidden to be false by default")
	}
}

func TestToggleCollapse(t *testing.T) {
	m := Model{
		collapsed: make(map[string]bool),
	}

	// Test toggling a directory that's not collapsed
	m.toggleCollapse("dir1")
	if !m.collapsed["dir1"] {
		t.Errorf("Expected dir1 to be collapsed after toggle")
	}

	// Test toggling a directory that's already collapsed
	m.toggleCollapse("dir1")
	if m.collapsed["dir1"] {
		t.Errorf("Expected dir1 to be expanded after second toggle")
	}
}

func TestExpandCollapseAllDirectories(t *testing.T) {
	m := Model{
		collapsed: make(map[string]bool),
		files: []filesystem.FileItem{
			{Path: "dir1", IsDir: true},
			{Path: "dir2", IsDir: true},
			{Path: "file1.txt", IsDir: false},
		},
	}

	// Test collapseAllDirectories
	m.collapseAllDirectories()
	if !m.collapsed["dir1"] || !m.collapsed["dir2"] {
		t.Errorf("Expected all directories to be collapsed")
	}

	// Test expandAllDirectories
	m.expandAllDirectories()
	if len(m.collapsed) != 0 {
		t.Errorf("Expected collapsed map to be empty after expanding all")
	}
}

func TestFilterSelections(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "filter-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a .gitignore file
	gitignoreContent := "*.log\nnode_modules/\n"
	err = os.WriteFile(filepath.Join(tempDir, ".gitignore"), []byte(gitignoreContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .gitignore file: %v", err)
	}

	// Create test files
	testFiles := []string{
		"file1.txt",
		"file2.log",
		".hidden/file3.txt",
		"node_modules/package.json",
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

	gitIgnoreMgr, _ := filesystem.NewGitIgnoreManager(tempDir)
	filterMgr := filesystem.NewFilterManager()

	m := Model{
		rootPath:     tempDir,
		gitIgnoreMgr: gitIgnoreMgr,
		filterMgr:    filterMgr,
		useGitIgnore: true,
		showHidden:   false,
		selected: map[string]bool{
			"file1.txt":                 true,
			"file2.log":                 true,
			".hidden/file3.txt":         true,
			"node_modules/package.json": true,
		},
		deselected: map[string]bool{},
	}

	// Test filterSelections
	m.filterSelections()

	// Should keep file1.txt
	if !m.selected["file1.txt"] {
		t.Errorf("Expected file1.txt to remain selected")
	}

	// Should remove file2.log (gitignored)
	if m.selected["file2.log"] {
		t.Errorf("Expected file2.log to be removed from selected (gitignored)")
	}

	// Should remove .hidden/file3.txt (hidden)
	if m.selected[".hidden/file3.txt"] {
		t.Errorf("Expected .hidden/file3.txt to be removed from selected (hidden)")
	}

	// Should remove node_modules/package.json (gitignored)
	if m.selected["node_modules/package.json"] {
		t.Errorf("Expected node_modules/package.json to be removed from selected (gitignored)")
	}

	// Test with gitignore disabled
	m.selected = map[string]bool{
		"file1.txt":                 true,
		"file2.log":                 true,
		".hidden/file3.txt":         true,
		"node_modules/package.json": true,
	}
	m.useGitIgnore = false
	m.filterSelections()

	// Should keep file1.txt
	if !m.selected["file1.txt"] {
		t.Errorf("Expected file1.txt to remain selected")
	}

	// Should keep file2.log (gitignore disabled)
	if !m.selected["file2.log"] {
		t.Errorf("Expected file2.log to remain selected (gitignore disabled)")
	}

	// Should remove .hidden/file3.txt (hidden)
	if m.selected[".hidden/file3.txt"] {
		t.Errorf("Expected .hidden/file3.txt to be removed from selected (hidden)")
	}

	// Should keep node_modules/package.json (gitignore disabled)
	if !m.selected["node_modules/package.json"] {
		t.Errorf("Expected node_modules/package.json to remain selected (gitignore disabled)")
	}

	// Test with hidden files enabled
	m.selected = map[string]bool{
		"file1.txt":                 true,
		"file2.log":                 true,
		".hidden/file3.txt":         true,
		"node_modules/package.json": true,
	}
	m.useGitIgnore = true
	m.showHidden = true
	m.filterSelections()

	// Should keep file1.txt
	if !m.selected["file1.txt"] {
		t.Errorf("Expected file1.txt to remain selected")
	}

	// Should remove file2.log (gitignored)
	if m.selected["file2.log"] {
		t.Errorf("Expected file2.log to be removed from selected (gitignored)")
	}

	// Should keep .hidden/file3.txt (hidden files enabled)
	if !m.selected[".hidden/file3.txt"] {
		t.Errorf("Expected .hidden/file3.txt to remain selected (hidden files enabled)")
	}

	// Should remove node_modules/package.json (gitignored)
	if m.selected["node_modules/package.json"] {
		t.Errorf("Expected node_modules/package.json to be removed from selected (gitignored)")
	}
}

func TestToggleSelection(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "selection-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test directory structure
	testFiles := []string{
		"dir1/file1.txt",
		"dir1/file2.txt",
		"dir2/file3.txt",
	}
	var fileItems []filesystem.FileItem
	fileItems = append(fileItems, filesystem.FileItem{Path: "dir1", IsDir: true, Level: 0})
	fileItems = append(fileItems, filesystem.FileItem{Path: "dir2", IsDir: true, Level: 0})

	for _, file := range testFiles {
		path := filepath.Join(tempDir, filepath.FromSlash(file))
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		if err := os.WriteFile(path, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", path, err)
		}
		relPath, _ := filepath.Rel(tempDir, path)
		fileItems = append(fileItems, filesystem.FileItem{
			Path:  filepath.ToSlash(relPath),
			IsDir: false,
			Level: strings.Count(filepath.ToSlash(relPath), "/"),
		})
	}

	config := Config{
		RootPath:    tempDir,
		FilterMgr:   filesystem.NewFilterManager(),
		MaxFileSize: math.MaxInt64,
	}
	m := NewModel(config)
	m.files = fileItems

	m.selected = make(map[string]bool)
	m.deselected = make(map[string]bool)
	m.isDependency = make(map[string]bool)

	// Test selecting a directory
	m.toggleSelection("dir1", true)
	if !m.selected["dir1"] {
		t.Errorf("Expected dir1 to be selected")
	}
	if !m.selected["dir1/file1.txt"] {
		t.Errorf("Expected dir1/file1.txt to be selected")
	}
	if !m.selected["dir1/file2.txt"] {
		t.Errorf("Expected dir1/file2.txt to be selected")
	}

	// Test deselecting a directory
	m.toggleSelection("dir1", true)
	if m.selected["dir1"] {
		t.Errorf("Expected dir1 to be deselected")
	}
	if m.selected["dir1/file1.txt"] {
		t.Errorf("Expected dir1/file1.txt to be deselected")
	}
	if m.selected["dir1/file2.txt"] {
		t.Errorf("Expected dir1/file2.txt to be deselected")
	}
	if !m.deselected["dir1/file1.txt"] {
		t.Errorf("Expected dir1/file1.txt to be marked as deselected")
	}
	if !m.deselected["dir1/file2.txt"] {
		t.Errorf("Expected dir1/file2.txt to be marked as deselected")
	}

	// Test selecting a file
	m.selected = make(map[string]bool)
	m.deselected = make(map[string]bool)
	m.toggleSelection("dir1/file1.txt", false)
	if !m.selected["dir1/file1.txt"] {
		t.Errorf("Expected dir1/file1.txt to be selected")
	}

	// Test deselecting a file
	m.toggleSelection("dir1/file1.txt", false)
	if m.selected["dir1/file1.txt"] {
		t.Errorf("Expected dir1/file1.txt to be deselected")
	}

	// Test selecting a file under a selected directory
	m.selected = make(map[string]bool)
	m.deselected = make(map[string]bool)
	m.toggleSelection("dir1", true)
	m.toggleSelection("dir1/file1.txt", false)
	if m.selected["dir1/file1.txt"] {
		t.Errorf("Expected dir1/file1.txt to be deselected")
	}
	if !m.deselected["dir1/file1.txt"] {
		t.Errorf("Expected dir1/file1.txt to be marked as deselected")
	}
}
