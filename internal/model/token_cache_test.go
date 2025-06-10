package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTokenCache_Basic(t *testing.T) {
	cache := NewTokenCache()
	defer cache.Close()

	// Create a temporary test file
	tmpFile, err := ioutil.TempFile("", "test_tokens_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write test content
	testContent := "hello world test content for token estimation"
	if _, err := tmpFile.WriteString(testContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	tokens, cached := cache.GetTokens(tmpFile.Name())
	if cached {
		t.Errorf("Expected not cached on first call, got cached=true")
	}
	if tokens != 0 {
		t.Errorf("Expected 0 tokens on first call, got %d", tokens)
	}

	time.Sleep(100 * time.Millisecond)

	tokens, cached = cache.GetTokens(tmpFile.Name())
	if !cached {
		t.Errorf("Expected cached on second call, got cached=false")
	}
	if tokens <= 0 {
		t.Errorf("Expected positive token count, got %d", tokens)
	}

	formatted := cache.GetTokensFormatted(tmpFile.Name())
	if formatted == "" {
		t.Errorf("Expected non-empty formatted string")
	}
}

func TestTokenCache_InvalidFile(t *testing.T) {
	cache := NewTokenCache()
	defer cache.Close()

	nonExistentFile := "/tmp/does_not_exist_test_file.txt"

	tokens, cached := cache.GetTokens(nonExistentFile)
	if cached {
		t.Errorf("Expected not cached on first call for invalid file")
	}

	time.Sleep(100 * time.Millisecond)

	tokens, cached = cache.GetTokens(nonExistentFile)
	if !cached {
		t.Errorf("Expected cached on second call for invalid file")
	}
	if tokens != 0 {
		t.Errorf("Expected 0 tokens for invalid file, got %d", tokens)
	}
}

func TestTokenCache_ClearAndInvalidate(t *testing.T) {
	cache := NewTokenCache()
	defer cache.Close()

	// Create a temporary test file
	tmpFile, err := ioutil.TempFile("", "test_clear_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString("test content")
	tmpFile.Close()

	cache.GetTokens(tmpFile.Name())
	time.Sleep(100 * time.Millisecond)

	_, cached := cache.GetTokens(tmpFile.Name())
	if !cached {
		t.Errorf("Expected file to be cached")
	}

	cache.InvalidateFile(tmpFile.Name())
	_, cached = cache.GetTokens(tmpFile.Name())
	if cached {
		t.Errorf("Expected file to be invalidated")
	}

	cache.GetTokens(tmpFile.Name())
	time.Sleep(100 * time.Millisecond)

	cache.ClearCache()
	_, cached = cache.GetTokens(tmpFile.Name())
	if cached {
		t.Errorf("Expected all cache to be cleared")
	}
}

func TestTokenCache_Stats(t *testing.T) {
	cache := NewTokenCache()
	defer cache.Close()

	cached, _ := cache.Stats()
	if cached != 0 {
		t.Errorf("Expected 0 cached files, got %d", cached)
	}

	tmpDir, err := ioutil.TempDir("", "test_stats_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	for i := 0; i < 3; i++ {
		tmpFile := filepath.Join(tmpDir, fmt.Sprintf("test%d.txt", i))
		if err := ioutil.WriteFile(tmpFile, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		cache.GetTokens(tmpFile)
	}

	time.Sleep(200 * time.Millisecond)
	cached, _ = cache.Stats()
	if cached != 3 {
		t.Errorf("Expected 3 cached files, got %d", cached)
	}
}

func BenchmarkTokenCache_SyncVsAsync(b *testing.B) {
	tmpDir, err := ioutil.TempDir("", "bench_tokens_")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFiles := make([]string, 10)
	for i := 0; i < 10; i++ {
		filePath := filepath.Join(tmpDir, fmt.Sprintf("test%d.txt", i))
		content := make([]byte, 1000)
		for j := range content {
			content[j] = byte('a' + (j % 26))
		}
		if err := ioutil.WriteFile(filePath, content, 0644); err != nil {
			b.Fatalf("Failed to create test file: %v", err)
		}
		testFiles[i] = filePath
	}

	b.Run("Synchronous", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, file := range testFiles {
				if contentBytes, err := os.ReadFile(file); err == nil {
					_ = len(string(contentBytes)) / 4
				}
			}
		}
	})

	b.Run("AsyncCache", func(b *testing.B) {
		cache := NewTokenCache()
		defer cache.Close()

		for _, file := range testFiles {
			cache.GetTokens(file)
		}
		time.Sleep(200 * time.Millisecond)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, file := range testFiles {
				cache.GetTokens(file)
			}
		}
	})
}