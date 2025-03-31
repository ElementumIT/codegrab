package model

import (
	"reflect"
	"testing"

	"github.com/epilande/codegrab/internal/filesystem"
)

func TestBuildDisplayNodes(t *testing.T) {
	m := Model{
		selected:   make(map[string]bool),
		deselected: make(map[string]bool),
		collapsed:  make(map[string]bool),
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

	// With dir1 expanded, we should have 5 nodes
	if len(m.displayNodes) != 5 {
		t.Fatalf("Expected 5 display nodes with dir1 expanded, got %d", len(m.displayNodes))
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

func TestAddNodeAndChildren(t *testing.T) {
	m := Model{
		selected:   make(map[string]bool),
		deselected: make(map[string]bool),
		collapsed:  make(map[string]bool),
		files: []filesystem.FileItem{
			{Path: "dir1", IsDir: true, Level: 0},
			{Path: "dir1/file1.txt", IsDir: false, Level: 1},
			{Path: "dir1/subdir", IsDir: true, Level: 1},
			{Path: "dir1/subdir/file2.txt", IsDir: false, Level: 2},
			{Path: "file3.txt", IsDir: false, Level: 0},
		},
	}

	// Test with collapsed directory
	m.collapsed["dir1"] = true
	m.displayNodes = nil
	m.addNodeAndChildren(filesystem.FileItem{Path: "dir1", IsDir: true, Level: 0}, 0, true)

	if len(m.displayNodes) != 1 {
		t.Fatalf("Expected 1 display node with collapsed dir1, got %d", len(m.displayNodes))
	}
	if m.displayNodes[0].Path != "dir1" {
		t.Errorf("Expected first node to be dir1, got %s", m.displayNodes[0].Path)
	}

	// Test with expanded directory
	m.collapsed["dir1"] = false
	m.displayNodes = nil
	m.addNodeAndChildren(filesystem.FileItem{Path: "dir1", IsDir: true, Level: 0}, 0, true)

	if len(m.displayNodes) != 4 {
		t.Fatalf("Expected 4 display nodes with expanded dir1, got %d", len(m.displayNodes))
	}
	if m.displayNodes[0].Path != "dir1" {
		t.Errorf("Expected first node to be dir1, got %s", m.displayNodes[0].Path)
	}
	if m.displayNodes[1].Path != "dir1/subdir" {
		t.Errorf("Expected second node to be dir1/subdir, got %s", m.displayNodes[1].Path)
	}
	if m.displayNodes[3].Path != "dir1/file1.txt" {
		t.Errorf("Expected fourth node to be dir1/file1.txt, got %s", m.displayNodes[3].Path)
	}

	// Test with nested expanded directories
	m.collapsed["dir1"] = false
	m.collapsed["dir1/subdir"] = false
	m.displayNodes = nil
	m.addNodeAndChildren(filesystem.FileItem{Path: "dir1", IsDir: true, Level: 0}, 0, true)

	if len(m.displayNodes) != 4 {
		t.Fatalf("Expected 4 display nodes with expanded nested dirs, got %d", len(m.displayNodes))
	}
	if m.displayNodes[0].Path != "dir1" {
		t.Errorf("Expected first node to be dir1, got %s", m.displayNodes[0].Path)
	}
	if m.displayNodes[1].Path != "dir1/subdir" {
		t.Errorf("Expected second node to be dir1/subdir, got %s", m.displayNodes[1].Path)
	}
	if m.displayNodes[2].Path != "dir1/subdir/file2.txt" {
		t.Errorf("Expected third node to be dir1/subdir/file2.txt, got %s", m.displayNodes[2].Path)
	}
	if m.displayNodes[3].Path != "dir1/file1.txt" {
		t.Errorf("Expected fourth node to be dir1/file1.txt, got %s", m.displayNodes[3].Path)
	}
}
