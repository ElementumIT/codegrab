package filesystem

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/epilande/codegrab/internal/cache"
	"github.com/epilande/codegrab/internal/utils"
)

// BenchmarkWalkDirectory compares performance of concurrent vs sequential walking
func BenchmarkWalkDirectory(b *testing.B) {
	// Create a test directory structure with multiple files
	tmpDir, err := os.MkdirTemp("", "walk-benchmark-")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a realistic directory structure
	setupBenchmarkDir(b, tmpDir)

	// Initialize managers
	filterMgr := NewFilterManager()
	var gitIgnoreMgr *GitIgnoreManager = nil

	b.ResetTimer()

	// Benchmark the concurrent implementation
	for i := 0; i < b.N; i++ {
		_, err := WalkDirectory(tmpDir, gitIgnoreMgr, filterMgr, false, false, math.MaxInt64)
		if err != nil {
			b.Fatalf("WalkDirectory failed: %v", err)
		}
	}
}

// BenchmarkWalkDirectoryLarge tests performance on a larger directory structure
func BenchmarkWalkDirectoryLarge(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "walk-benchmark-large-")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a larger, more complex directory structure
	setupLargeBenchmarkDir(b, tmpDir)

	filterMgr := NewFilterManager()
	var gitIgnoreMgr *GitIgnoreManager = nil

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		files, err := WalkDirectory(tmpDir, gitIgnoreMgr, filterMgr, false, false, math.MaxInt64)
		if err != nil {
			b.Fatalf("WalkDirectory failed: %v", err)
		}
		// Ensure we're actually processing files
		if len(files) == 0 {
			b.Fatalf("No files found in benchmark directory")
		}
	}
}

// BenchmarkConcurrentVsSequential compares concurrent vs sequential processing
func BenchmarkConcurrentVsSequential(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "concurrent-vs-sequential-")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	setupLargeBenchmarkDir(b, tmpDir)

	filterMgr := NewFilterManager()
	var gitIgnoreMgr *GitIgnoreManager = nil

	b.Run("Concurrent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := WalkDirectory(tmpDir, gitIgnoreMgr, filterMgr, false, false, math.MaxInt64)
			if err != nil {
				b.Fatalf("Concurrent WalkDirectory failed: %v", err)
			}
		}
	})

	b.Run("Original", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := walkDirectoryOriginal(tmpDir, gitIgnoreMgr, filterMgr, false, false, math.MaxInt64)
			if err != nil {
				b.Fatalf("Original WalkDirectory failed: %v", err)
			}
		}
	})
}

// BenchmarkWorkerScaling tests performance with different worker counts
func BenchmarkWorkerScaling(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "worker-scaling-")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	setupLargeBenchmarkDir(b, tmpDir)

	filterMgr := NewFilterManager()
	var gitIgnoreMgr *GitIgnoreManager = nil

	workerCounts := []int{1, 2, 4, runtime.NumCPU(), runtime.NumCPU() * 2}

	for _, workers := range workerCounts {
		b.Run(fmt.Sprintf("Workers-%d", workers), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				walker := &ConcurrentWalker{
					root:         tmpDir,
					gitIgnore:    gitIgnoreMgr,
					filter:       filterMgr,
					useGitIgnore: false,
					showHidden:   false,
					maxFileSize:  math.MaxInt64,
					maxWorkers:   workers,
				}
				_, err := walker.walk()
				if err != nil {
					b.Fatalf("Walker failed with %d workers: %v", workers, err)
				}
			}
		})
	}
}

// setupBenchmarkDir creates a realistic directory structure for benchmarking
func setupBenchmarkDir(b *testing.B, baseDir string) {
	// Create directory structure similar to a real project
	dirs := []string{
		"src", "src/components", "src/utils", "src/tests",
		"lib", "lib/internal", "lib/external",
		"docs", "config", "scripts",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(baseDir, dir), 0755); err != nil {
			b.Fatalf("Failed to create dir %s: %v", dir, err)
		}
	}

	// Create files in each directory
	files := []struct {
		path    string
		content string
	}{
		{"src/main.go", "package main\n\nfunc main() {\n\tprintln(\"Hello, World!\")\n}"},
		{"src/components/button.go", "package components\n\ntype Button struct {\n\tText string\n}"},
		{"src/utils/helper.go", "package utils\n\nfunc Helper() string {\n\treturn \"helper\"\n}"},
		{"src/tests/main_test.go", "package main\n\nimport \"testing\"\n\nfunc TestMain(t *testing.T) {}"},
		{"lib/internal/logger.go", "package internal\n\ntype Logger struct {}"},
		{"lib/external/client.go", "package external\n\ntype Client struct {}"},
		{"docs/README.md", "# Project Documentation"},
		{"config/config.yaml", "version: 1.0\nname: test"},
		{"scripts/build.sh", "#!/bin/bash\necho 'Building...'"},
	}

	for _, file := range files {
		fullPath := filepath.Join(baseDir, file.path)
		if err := os.WriteFile(fullPath, []byte(file.content), 0644); err != nil {
			b.Fatalf("Failed to create file %s: %v", fullPath, err)
		}
	}
}

// setupLargeBenchmarkDir creates a larger directory structure
func setupLargeBenchmarkDir(b *testing.B, baseDir string) {
	// Create more extensive directory structure
	for i := 0; i < 10; i++ {
		for j := 0; j < 5; j++ {
			dir := filepath.Join(baseDir, fmt.Sprintf("dir%d", i), fmt.Sprintf("subdir%d", j))
			if err := os.MkdirAll(dir, 0755); err != nil {
				b.Fatalf("Failed to create dir %s: %v", dir, err)
			}

			// Create multiple files per directory
			for k := 0; k < 3; k++ {
				content := fmt.Sprintf("package dir%d\n\n// File %d in directory %d/%d\n", i, k, i, j)
				for l := 0; l < 50; l++ {
					content += fmt.Sprintf("// Line %d of content\n", l)
				}

				fileName := filepath.Join(dir, fmt.Sprintf("file%d.go", k))
				if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
					b.Fatalf("Failed to create file %s: %v", fileName, err)
				}
			}
		}
	}
}

// Add comprehensive profiling benchmarks
func BenchmarkWalkDirectoryImplementations(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "walk-implementations-")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	setupLargeBenchmarkDir(b, tmpDir)

	filterMgr := NewFilterManager()
	var gitIgnoreMgr *GitIgnoreManager = nil

	// Test all three implementations
	implementations := []struct {
		name string
		fn   func(string, *GitIgnoreManager, *FilterManager, bool, bool, int64) ([]FileItem, error)
	}{
		{"Concurrent", WalkDirectory},
 
		{"Original", walkDirectoryOriginal},
	}

	for _, impl := range implementations {
		b.Run(impl.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				files, err := impl.fn(tmpDir, gitIgnoreMgr, filterMgr, false, false, math.MaxInt64)
				if err != nil {
					b.Fatalf("%s implementation failed: %v", impl.name, err)
				}
				// Ensure we're actually processing files
				if len(files) == 0 {
					b.Fatalf("No files found in %s implementation", impl.name)
				}
			}
		})
	}
}

// walkDirectoryOriginal is the original sequential implementation for benchmark comparison only
func walkDirectoryOriginal(root string, gitIgnore *GitIgnoreManager, filter *FilterManager, useGitIgnore, showHidden bool, maxFileSize int64) ([]FileItem, error) {
	var files []FileItem

	if _, err := os.Stat(root); err != nil {
		return nil, fmt.Errorf("failed to access root directory: %w", err)
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: skipping %s: %v\n", path, err)
			return nil
		}
		if path == root {
			return nil
		}

		// Skip hidden directories/files
		if !showHidden && strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip gitignored paths
		if useGitIgnore && gitIgnore != nil && gitIgnore.IsIgnored(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %s: %w", path, err)
		}
		if relPath == "." {
			return nil
		}
		relPath = filepath.ToSlash(relPath)

		// Always include directories when using include patterns
		// This allows us to traverse into directories that might contain matching files
		if info.IsDir() {
			files = append(files, FileItem{
				Path:  relPath,
				IsDir: true,
				Level: strings.Count(relPath, "/"),
				Size:  info.Size(),
			})
			return nil
		}

		// Skip files larger than maxFileSize
		if info.Size() > maxFileSize {
			return nil
		}

		// Skip files not matching glob patterns
		if !filter.ShouldInclude(relPath) {
			return nil
		}

		fileCache := cache.GetGlobalFileCache()
		if ok, err := fileCache.GetTextFileStatus(path, utils.IsTextFile); err != nil || !ok {
			return nil
		}
		files = append(files, FileItem{
			Path:  relPath,
			IsDir: false,
			Level: strings.Count(relPath, "/"),
			Size:  info.Size(),
		})
		return nil
	})
	if err != nil {
		return files, fmt.Errorf("error walking directory: %w", err)
	}

	return files, nil
}