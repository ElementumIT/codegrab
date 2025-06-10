package cache

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestFileCache_BasicOperations(t *testing.T) {
	// Create a temporary file for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := "Hello, World!"

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cache := NewFileCache(1024*1024, 100) // 1MB, 100 items

	// Test cache miss - should load from disk
	content, err := cache.Get(testFile)
	if err != nil {
		t.Fatalf("Failed to get file content: %v", err)
	}
	if content != testContent {
		t.Errorf("Expected content %q, got %q", testContent, content)
	}

	// Test cache hit - should return cached content
	content2, err := cache.Get(testFile)
	if err != nil {
		t.Fatalf("Failed to get cached file content: %v", err)
	}
	if content2 != testContent {
		t.Errorf("Expected cached content %q, got %q", testContent, content2)
	}

	// Verify file is in cache
	if !cache.Contains(testFile) {
		t.Error("File should be in cache")
	}
}

func TestFileCache_ModificationDetection(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	// Create initial file
	initialContent := "Initial content"
	err := os.WriteFile(testFile, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cache := NewFileCache(1024*1024, 100)

	// Load into cache
	content, err := cache.Get(testFile)
	if err != nil {
		t.Fatalf("Failed to get file content: %v", err)
	}
	if content != initialContent {
		t.Errorf("Expected content %q, got %q", initialContent, content)
	}

	// Wait a moment to ensure different modification time
	time.Sleep(10 * time.Millisecond)

	// Modify the file
	modifiedContent := "Modified content"
	err = os.WriteFile(testFile, []byte(modifiedContent), 0644)
	if err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}

	// Should detect modification and reload
	content, err = cache.Get(testFile)
	if err != nil {
		t.Fatalf("Failed to get modified file content: %v", err)
	}
	if content != modifiedContent {
		t.Errorf("Expected modified content %q, got %q", modifiedContent, content)
	}
}

func TestFileCache_TextFileStatus(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	err := os.WriteFile(testFile, []byte("Hello"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cache := NewFileCache(1024*1024, 100)

	// Mock function that checks if file is text
	checkFunc := func(path string) (bool, error) {
		return strings.HasSuffix(path, ".txt"), nil
	}

	// First call should execute checkFunc
	isText, err := cache.GetTextFileStatus(testFile, checkFunc)
	if err != nil {
		t.Fatalf("Failed to get text file status: %v", err)
	}
	if !isText {
		t.Error("Expected file to be detected as text")
	}

	// Second call should return cached result
	isText2, err := cache.GetTextFileStatus(testFile, checkFunc)
	if err != nil {
		t.Fatalf("Failed to get cached text file status: %v", err)
	}
	if !isText2 {
		t.Error("Expected cached result to be text")
	}
}

func TestFileCache_LRUEviction(t *testing.T) {
	tmpDir := t.TempDir()
	cache := NewFileCache(1024, 2) // Small cache: 1KB, 2 items max

	// Create test files
	files := make([]string, 3)
	for i := 0; i < 3; i++ {
		files[i] = filepath.Join(tmpDir, "test"+string(rune('0'+i))+".txt")
		content := strings.Repeat("X", 200) // 200 bytes each
		err := os.WriteFile(files[i], []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %d: %v", i, err)
		}
	}

	// Load first two files
	_, err := cache.Get(files[0])
	if err != nil {
		t.Fatalf("Failed to load file 0: %v", err)
	}

	_, err = cache.Get(files[1])
	if err != nil {
		t.Fatalf("Failed to load file 1: %v", err)
	}

	// Both should be in cache
	if !cache.Contains(files[0]) {
		t.Error("File 0 should be in cache")
	}
	if !cache.Contains(files[1]) {
		t.Error("File 1 should be in cache")
	}

	// Add a small delay to ensure different access times
	time.Sleep(10 * time.Millisecond)

	// Access file 0 to make it more recently used
	_, err = cache.Get(files[0])
	if err != nil {
		t.Fatalf("Failed to access file 0: %v", err)
	}

	// Load third file - should evict file 1 (LRU)
	_, err = cache.Get(files[2])
	if err != nil {
		t.Fatalf("Failed to load file 2: %v", err)
	}

	// File 0 and 2 should be in cache, file 1 should be evicted
	if !cache.Contains(files[0]) {
		t.Error("File 0 should still be in cache (recently used)")
	}
	if cache.Contains(files[1]) {
		t.Error("File 1 should be evicted (LRU)")
	}
	if !cache.Contains(files[2]) {
		t.Error("File 2 should be in cache (newly added)")
	}
}

func TestFileCache_SizeBasedEviction(t *testing.T) {
	tmpDir := t.TempDir()
	cache := NewFileCache(500, 100) // 500 bytes max, 100 items max

	// Create files that will exceed size limit
	files := make([]string, 3)
	for i := 0; i < 3; i++ {
		files[i] = filepath.Join(tmpDir, "test"+string(rune('0'+i))+".txt")
		content := strings.Repeat("X", 200) // 200 bytes each
		err := os.WriteFile(files[i], []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %d: %v", i, err)
		}
	}

	// Load first two files (400 bytes total)
	_, err := cache.Get(files[0])
	if err != nil {
		t.Fatalf("Failed to load file 0: %v", err)
	}

	_, err = cache.Get(files[1])
	if err != nil {
		t.Fatalf("Failed to load file 1: %v", err)
	}

	// Load third file (200 more bytes) - should trigger eviction
	_, err = cache.Get(files[2])
	if err != nil {
		t.Fatalf("Failed to load file 2: %v", err)
	}

	// Check cache stats
	items, sizeBytes, _ := cache.Stats()
	if items > 2 {
		t.Errorf("Expected at most 2 items in cache, got %d", items)
	}
	if sizeBytes > 500 {
		t.Errorf("Expected cache size <= 500 bytes, got %d", sizeBytes)
	}
}

func TestFileCache_Clear(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	err := os.WriteFile(testFile, []byte("Hello"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cache := NewFileCache(1024*1024, 100)

	// Load file into cache
	_, err = cache.Get(testFile)
	if err != nil {
		t.Fatalf("Failed to load file: %v", err)
	}

	// Verify it's in cache
	if !cache.Contains(testFile) {
		t.Error("File should be in cache")
	}

	// Clear cache
	cache.Clear()

	// Verify it's no longer in cache
	if cache.Contains(testFile) {
		t.Error("File should not be in cache after clear")
	}

	// Verify stats
	items, sizeBytes, _ := cache.Stats()
	if items != 0 {
		t.Errorf("Expected 0 items after clear, got %d", items)
	}
	if sizeBytes != 0 {
		t.Errorf("Expected 0 bytes after clear, got %d", sizeBytes)
	}
}

func TestFileCache_NonExistentFile(t *testing.T) {
	cache := NewFileCache(1024*1024, 100)

	_, err := cache.Get("/nonexistent/file.txt")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestDefaultFileCache(t *testing.T) {
	cache := DefaultFileCache()

	items, sizeBytes, _ := cache.Stats()
	if items != 0 {
		t.Errorf("Expected 0 items in new cache, got %d", items)
	}
	if sizeBytes != 0 {
		t.Errorf("Expected 0 bytes in new cache, got %d", sizeBytes)
	}
}

func TestCacheContentAfterTextFileCheck(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := "Hello, World! This is test content."

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cache := NewFileCache(1024*1024, 100)

	checkFunc := func(path string) (bool, error) {
		return true, nil
	}

	isText, err := cache.GetTextFileStatus(testFile, checkFunc)
	if err != nil {
		t.Fatalf("GetTextFileStatus failed: %v", err)
	}
	if !isText {
		t.Error("Expected file to be detected as text")
	}

	content, err := cache.Get(testFile)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if content != testContent {
		t.Errorf("Expected content %q, got %q", testContent, content)
	}

	items, sizeBytes, _ := cache.Stats()
	if items != 1 {
		t.Errorf("Expected 1 cache item, got %d", items)
	}
	expectedSize := int64(len(testContent))
	if sizeBytes != expectedSize {
		t.Errorf("Expected cache size %d bytes, got %d", expectedSize, sizeBytes)
	}
}

func TestMultipleTextFileChecksAndReads(t *testing.T) {
	tmpDir := t.TempDir()

	testFiles := []struct {
		name    string
		content string
	}{
		{"file1.txt", "Content of file 1"},
		{"file2.go", "package main\nfunc main() {}"},
		{"file3.md", "# Markdown content"},
	}

	for _, tf := range testFiles {
		filePath := filepath.Join(tmpDir, tf.name)
		err := os.WriteFile(filePath, []byte(tf.content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", tf.name, err)
		}
	}

	cache := NewFileCache(1024*1024, 100)
	checkFunc := func(path string) (bool, error) {
		return true, nil
	}

	for _, tf := range testFiles {
		filePath := filepath.Join(tmpDir, tf.name)

		isText, err := cache.GetTextFileStatus(filePath, checkFunc)
		if err != nil {
			t.Fatalf("GetTextFileStatus failed for %s: %v", tf.name, err)
		}
		if !isText {
			t.Errorf("Expected %s to be detected as text", tf.name)
		}
	}

	for _, tf := range testFiles {
		filePath := filepath.Join(tmpDir, tf.name)

		content, err := cache.Get(filePath)
		if err != nil {
			t.Fatalf("Get failed for %s: %v", tf.name, err)
		}
		if content != tf.content {
			t.Errorf("Content mismatch for %s: expected %q, got %q", tf.name, tf.content, content)
		}
	}

	items, sizeBytes, _ := cache.Stats()
	if items != len(testFiles) {
		t.Errorf("Expected %d cache items, got %d", len(testFiles), items)
	}

	expectedTotalSize := int64(0)
	for _, tf := range testFiles {
		expectedTotalSize += int64(len(tf.content))
	}

	if sizeBytes != expectedTotalSize {
		t.Errorf("Expected total cache size %d bytes, got %d", expectedTotalSize, sizeBytes)
	}
}

