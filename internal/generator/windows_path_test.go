package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestBuildTreeWithWindowsPaths tests the tree building with Windows-style paths
func TestBuildTreeWithWindowsPaths(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "windows-path-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files using OS-specific paths
	testFiles := []struct {
		path    string
		content string
	}{
		{"file1.txt", "Content of file1"},
		{"file2.go", "package main\n\nfunc main() {}"},
		{"subdir/file3.txt", "Content of file3"},
		{"subdir/nested/file4.js", "console.log('Hello');"},
	}

	for _, tf := range testFiles {
		path := filepath.Join(tempDir, filepath.FromSlash(tf.path))
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		if err := os.WriteFile(path, []byte(tf.content), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", path, err)
		}
	}

	gen := NewGenerator(tempDir, nil, nil, "", false)

	// Test with Windows-style paths in SelectedFiles
	gen.SelectedFiles = make(map[string]bool)
	
	// Simulate what might come from Windows - paths with backslashes
	windowsPaths := []string{
		"file1.txt",
		"file2.go",
		"subdir\\file3.txt",       // Windows-style
		"subdir\\nested\\file4.js", // Windows-style nested
	}

	// Add these to selected files
	for _, path := range windowsPaths {
		gen.SelectedFiles[path] = true
	}

	// Build the tree
	root := gen.buildTree()

	// Verify the tree structure
	if root.Name != filepath.Base(tempDir) {
		t.Errorf("Expected root name to be %q, got %q", filepath.Base(tempDir), root.Name)
	}
	if !root.IsDir {
		t.Errorf("Expected root to be a directory")
	}

	// Check that we have the right number of children at root level
	// We should have: file1.txt, file2.go, and subdir/
	expectedRootChildren := 3 // file1.txt, file2.go, subdir
	if len(root.Children) != expectedRootChildren {
		t.Errorf("Expected root to have %d children, got %d", expectedRootChildren, len(root.Children))
		t.Logf("Root children:")
		for i, child := range root.Children {
			t.Logf("  %d: %s (IsDir: %v)", i, child.Name, child.IsDir)
		}
	}

	// Find the subdir
	var subdirNode *Node
	for _, child := range root.Children {
		if child.Name == "subdir" && child.IsDir {
			subdirNode = child
			break
		}
	}

	if subdirNode == nil {
		t.Fatalf("Expected to find 'subdir' directory in root children")
	}

	// Check subdir contents - should have file3.txt and nested/
	expectedSubdirChildren := 2 // file3.txt and nested/
	if len(subdirNode.Children) != expectedSubdirChildren {
		t.Errorf("Expected subdir to have %d children, got %d", expectedSubdirChildren, len(subdirNode.Children))
		t.Logf("Subdir children:")
		for i, child := range subdirNode.Children {
			t.Logf("  %d: %s (IsDir: %v)", i, child.Name, child.IsDir)
		}
	}

	// Find the nested dir
	var nestedNode *Node
	for _, child := range subdirNode.Children {
		if child.Name == "nested" && child.IsDir {
			nestedNode = child
			break
		}
	}

	if nestedNode == nil {
		t.Fatalf("Expected to find 'nested' directory in subdir children")
	}

	// Check nested contents - should have file4.js
	if len(nestedNode.Children) != 1 {
		t.Errorf("Expected nested to have 1 child, got %d", len(nestedNode.Children))
	}

	if nestedNode.Children[0].Name != "file4.js" || nestedNode.Children[0].IsDir {
		t.Errorf("Expected nested child to be 'file4.js' file, got %s (IsDir: %v)", nestedNode.Children[0].Name, nestedNode.Children[0].IsDir)
	}

	// Test tree rendering to see if it's flat
	var builder strings.Builder
	renderTree(root, "", true, &builder, root.Name, make(map[string]bool))
	
	treeOutput := builder.String()
	t.Logf("Tree output:\n%s", treeOutput)
	
	// The tree should have proper nesting, not be flat
	// Check for the presence of tree symbols that indicate nesting
	if !strings.Contains(treeOutput, "├──") && !strings.Contains(treeOutput, "└──") {
		t.Errorf("Tree output appears to be flat - missing tree structure symbols")
	}
	
	// Check for proper nesting symbols
	if !strings.Contains(treeOutput, "│   ") {
		t.Errorf("Tree output is missing nested indentation symbols")
	}
}