package model

import (
	"reflect"
	"testing"

	"github.com/epilande/codegrab/internal/filesystem"
)

func TestBuildDisplayNodes(t *testing.T) {
	m := Model{
		selected:     make(map[string]bool),
		deselected:   make(map[string]bool),
		collapsed:    make(map[string]bool),
		isDependency: make(map[string]bool),
		files: []filesystem.FileItem{
			{Path: "dir1", IsDir: true, Level: 0},
			{Path: "dir1/file1.txt", IsDir: false, Level: 1},
			{Path: "dir2", IsDir: true, Level: 0},
			{Path: "file2.txt", IsDir: false, Level: 0},
			{Path: "file3.txt", IsDir: false, Level: 0},
		},
	}

	// Mark some files as selected
	m.selected["dir1"] = true
	m.selected["file2.txt"] = true

	// Mark dir1 as collapsed
	m.collapsed["dir1"] = true

	// Build display nodes
	m.buildDisplayNodes()

	// Expected nodes (directories first, then files, all sorted alphabetically)
	expected := []FileNode{
		{Path: "dir1", Name: "dir1", IsDir: true, Level: 0, IsLast: false, Selected: true},
		{Path: "dir2", Name: "dir2", IsDir: true, Level: 0, IsLast: false, Selected: false},
		{Path: "file2.txt", Name: "file2.txt", IsDir: false, Level: 0, IsLast: false, Selected: true},
		{Path: "file3.txt", Name: "file3.txt", IsDir: false, Level: 0, IsLast: true, Selected: false},
	}

	// Check length
	if len(m.displayNodes) != len(expected) {
		t.Fatalf("Expected %d display nodes, got %d", len(expected), len(m.displayNodes))
	}

	// Check each node
	for i, expectedNode := range expected {
		actualNode := m.displayNodes[i]
		if !reflect.DeepEqual(actualNode, expectedNode) {
			t.Errorf("Node %d mismatch:\nExpected: %+v\nActual: %+v", i, expectedNode, actualNode)
		}
	}

	// Test with expanded directory
	m.collapsed["dir1"] = false
	m.buildDisplayNodes()

	// Expected nodes with dir1 expanded
	expectedExpanded := []FileNode{
		{Path: "dir1", Name: "dir1", IsDir: true, Level: 0, IsLast: false, Selected: true},
		{Path: "dir1/file1.txt", Name: "file1.txt", IsDir: false, Level: 1, IsLast: true, Selected: false},
		{Path: "dir2", Name: "dir2", IsDir: true, Level: 0, IsLast: false, Selected: false},
		{Path: "file2.txt", Name: "file2.txt", IsDir: false, Level: 0, IsLast: false, Selected: true},
		{Path: "file3.txt", Name: "file3.txt", IsDir: false, Level: 0, IsLast: true, Selected: false},
	}

	// Check length with dir1 expanded
	if len(m.displayNodes) != len(expectedExpanded) {
		t.Fatalf("Expected %d display nodes with dir1 expanded, got %d", len(expectedExpanded), len(m.displayNodes))
	}

	// Check each node with dir1 expanded
	for i, expectedNode := range expectedExpanded {
		actualNode := m.displayNodes[i]
		if !reflect.DeepEqual(actualNode, expectedNode) {
			t.Errorf("Node %d mismatch (expanded):\nExpected: %+v\nActual: %+v", i, expectedNode, actualNode)
		}
	}

	// Check that dir1/file1.txt is now included
	found := false
	for _, node := range m.displayNodes {
		if node.Path == "dir1/file1.txt" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected dir1/file1.txt to be included in display nodes when dir1 is expanded")
	}
}
