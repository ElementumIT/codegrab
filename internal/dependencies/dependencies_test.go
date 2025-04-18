package dependencies

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestEnv(t *testing.T, files map[string]string) (string, func()) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "deps-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	if err := os.Chmod(tempDir, 0755); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to chmod temp dir: %v", err)
	}

	for path, content := range files {
		fullPath := filepath.Join(tempDir, filepath.FromSlash(path))
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			os.RemoveAll(tempDir)
			t.Fatalf("Failed to create dir %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			os.RemoveAll(tempDir)
			t.Fatalf("Failed to write file %s: %v", fullPath, err)
		}
	}

	return tempDir, func() { os.RemoveAll(tempDir) }
}
