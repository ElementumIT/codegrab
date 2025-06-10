package model

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/epilande/codegrab/internal/filesystem"
)

// TestTokenCacheIntegration tests the token cache working with the full model
func TestTokenCacheIntegration(t *testing.T) {
	// Create a temporary directory with test files
	tmpDir, err := ioutil.TempDir("", "token_integration_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFiles := map[string]string{
		"main.go":    "package main\n\nfunc main() {\n\tprintln(\"hello world\")\n}",
		"utils.go":   "package main\n\nfunc helper() string {\n\treturn \"utility function\"\n}",
		"README.md":  "# Test Project\n\nThis is a test project for token cache integration.",
	}

	for filename, content := range testFiles {
		filePath := tmpDir + "/" + filename
		if err := ioutil.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	filterMgr := filesystem.NewFilterManager()

	config := Config{
		FilterMgr:      filterMgr,
		RootPath:       tmpDir,
		OutputPath:     "/tmp/test_integration.md",
		Format:         "markdown",
		MaxDepth:       10,
		MaxFileSize:    1024 * 1024,
		UseTempFile:    false,
		SkipRedaction:  true,
		ResolveDeps:    false,
		ShowIcons:      false,
		ShowTokenCount: true,
	}

	model := NewModel(config)
	defer func() {
		if model.tokenCache != nil {
			model.tokenCache.Close()
		}
	}()

	if model.tokenCache == nil {
		t.Fatal("Token cache should be initialized")
	}

	files, err := filesystem.WalkDirectory(tmpDir, model.gitIgnoreMgr, model.filterMgr, true, false, model.maxFileSize)
	if err != nil {
		t.Fatalf("Failed to walk directory: %v", err)
	}
	model.files = files

	model.buildDisplayNodes()

	time.Sleep(200 * time.Millisecond)

	for filename := range testFiles {
		_ = model.tokenCache.GetTokensFormatted(filename)
	}

	model.refreshViewportContent()

	if model.viewport.View() == "" {
		t.Error("Viewport should have content after refresh")
	}
}

// TestModelShowTokenCountFlag tests that the showTokenCount flag works correctly
func TestModelShowTokenCountFlag(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "token_flag_test_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filterMgr := filesystem.NewFilterManager()

	configNoTokens := Config{
		FilterMgr:      filterMgr,
		RootPath:       tmpDir,
		OutputPath:     "/tmp/test.md",
		Format:         "markdown",
		MaxDepth:       10,
		MaxFileSize:    1024 * 1024,
		ShowTokenCount: false,
	}

	modelNoTokens := NewModel(configNoTokens)
	defer modelNoTokens.tokenCache.Close()

	if modelNoTokens.showTokenCount {
		t.Error("showTokenCount should be false when disabled in config")
	}

	configWithTokens := Config{
		FilterMgr:      filterMgr,
		RootPath:       tmpDir,
		OutputPath:     "/tmp/test.md",
		Format:         "markdown",
		MaxDepth:       10,
		MaxFileSize:    1024 * 1024,
		ShowTokenCount: true,
	}

	modelWithTokens := NewModel(configWithTokens)
	defer modelWithTokens.tokenCache.Close()

	if !modelWithTokens.showTokenCount {
		t.Error("showTokenCount should be true when enabled in config")
	}
	if modelNoTokens.tokenCache == nil {
		t.Error("Token cache should be initialized even when showTokenCount is false")
	}
	if modelWithTokens.tokenCache == nil {
		t.Error("Token cache should be initialized when showTokenCount is true")
	}
}