package generator

import (
	"testing"
)

func TestFormatInterface(t *testing.T) {
	var _ Format = &mockFormat{}
}

func TestTemplateData(t *testing.T) {
	data := TemplateData{
		Structure: "root/\n└── file.txt\n",
		Files: []FileData{
			{
				Path:     "file.txt",
				Content:  "File content",
				Language: "text",
			},
		},
	}

	if data.Structure != "root/\n└── file.txt\n" {
		t.Errorf("Expected Structure to be %q, got %q", "root/\n└── file.txt\n", data.Structure)
	}

	if len(data.Files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(data.Files))
	}

	file := data.Files[0]
	if file.Path != "file.txt" {
		t.Errorf("Expected Path to be %q, got %q", "file.txt", file.Path)
	}
	if file.Content != "File content" {
		t.Errorf("Expected Content to be %q, got %q", "File content", file.Content)
	}
	if file.Language != "text" {
		t.Errorf("Expected Language to be %q, got %q", "text", file.Language)
	}
}
