package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/epilande/codegrab/internal/cache"
	"github.com/epilande/codegrab/internal/utils"
)

// FileItem represents a file or directory found when walking the filesystem.
type FileItem struct {
	Path  string
	IsDir bool
	Level int
	Size  int64
}

// pathItem represents a discovered path to be processed
type pathItem struct {
	path string
	info os.FileInfo
}

// ConcurrentWalker handles concurrent file system traversal
type ConcurrentWalker struct {
	root         string
	gitIgnore    *GitIgnoreManager
	filter       *FilterManager
	useGitIgnore bool
	showHidden   bool
	maxFileSize  int64
	maxWorkers   int
}

// WalkDirectory traverses the root directory taking into account gitignore, hidden files, and max file size.
func WalkDirectory(root string, gitIgnore *GitIgnoreManager, filter *FilterManager, useGitIgnore, showHidden bool, maxFileSize int64) ([]FileItem, error) {
	walker := &ConcurrentWalker{
		root:         root,
		gitIgnore:    gitIgnore,
		filter:       filter,
		useGitIgnore: useGitIgnore,
		showHidden:   showHidden,
		maxFileSize:  maxFileSize,
		maxWorkers:   runtime.NumCPU(),
	}
	return walker.walk()
}


// walk performs the concurrent directory traversal
func (w *ConcurrentWalker) walk() ([]FileItem, error) {
	if _, err := os.Stat(w.root); err != nil {
		return nil, fmt.Errorf("failed to access root directory: %w", err)
	}

	// Create channels with buffered capacity
	pathQueue := make(chan pathItem, 1000)
	resultQueue := make(chan FileItem, 1000)
	errorChan := make(chan error, w.maxWorkers)

	var wg sync.WaitGroup
	var files []FileItem
	var firstError error

	// Start worker goroutines
	for i := 0; i < w.maxWorkers; i++ {
		wg.Add(1)
		go w.processWorker(pathQueue, resultQueue, errorChan, &wg)
	}

	// Start result collector goroutine
	done := make(chan struct{})
	go func() {
		defer close(done)
		for item := range resultQueue {
			files = append(files, item)
		}
	}()

	// Start error collector goroutine
	errorDone := make(chan struct{})
	go func() {
		defer close(errorDone)
		for err := range errorChan {
			if firstError == nil {
				firstError = err
			}
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		}
	}()

	// Discover paths using sequential filepath.Walk (directory traversal is I/O bound)
	walkErr := filepath.Walk(w.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			select {
			case errorChan <- fmt.Errorf("skipping %s: %w", path, err):
			default:
			}
			return nil
		}

		if path == w.root {
			return nil
		}

		// Early filtering for directories to skip entire subtrees
		if info.IsDir() {
			// Skip hidden directories
			if !w.showHidden && strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			// Skip gitignored directories
			if w.useGitIgnore && w.gitIgnore != nil && w.gitIgnore.IsIgnored(path) {
				return filepath.SkipDir
			}
		}

		// Queue path for processing
		select {
		case pathQueue <- pathItem{path: path, info: info}:
		case <-done: // Stop if result collector finished (unlikely but safe)
			return fmt.Errorf("result collection finished early")
		}

		return nil
	})

	// Close path queue to signal workers to finish
	close(pathQueue)

	// Wait for all workers to complete
	wg.Wait()
	close(resultQueue)
	close(errorChan)

	// Wait for collectors to finish
	<-done
	<-errorDone

	if walkErr != nil {
		return files, fmt.Errorf("error walking directory: %w", walkErr)
	}

	return files, firstError
}

// processWorker processes paths from the queue
func (w *ConcurrentWalker) processWorker(pathQueue <-chan pathItem, resultQueue chan<- FileItem, errorChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for item := range pathQueue {
		fileItem, err := w.processPath(item.path, item.info)
		if err != nil {
			select {
			case errorChan <- err:
			default:
			}
			continue
		}

		if fileItem != nil {
			select {
			case resultQueue <- *fileItem:
			default:
			}
		}
	}
}

// processPath handles the processing of a single path
func (w *ConcurrentWalker) processPath(path string, info os.FileInfo) (*FileItem, error) {
	// Skip hidden files/directories
	if !w.showHidden && strings.HasPrefix(info.Name(), ".") {
		return nil, nil
	}

	// Skip gitignored paths
	if w.useGitIgnore && w.gitIgnore != nil && w.gitIgnore.IsIgnored(path) {
		return nil, nil
	}

	relPath, err := filepath.Rel(w.root, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get relative path for %s: %w", path, err)
	}
	if relPath == "." {
		return nil, nil
	}
	relPath = filepath.ToSlash(relPath)

	// Handle directories
	if info.IsDir() {
		return &FileItem{
			Path:  relPath,
			IsDir: true,
			Level: strings.Count(relPath, "/"),
			Size:  info.Size(),
		}, nil
	}

	// Handle files
	// Skip files larger than maxFileSize
	if info.Size() > w.maxFileSize {
		return nil, nil
	}

	// Skip files not matching glob patterns
	if !w.filter.ShouldInclude(relPath) {
		return nil, nil
	}

	// Check if file is text (this is the most expensive operation)
	fileCache := cache.GetGlobalFileCache()
	if ok, err := fileCache.GetTextFileStatus(path, utils.IsTextFile); err != nil || !ok {
		return nil, nil
	}

	return &FileItem{
		Path:  relPath,
		IsDir: false,
		Level: strings.Count(relPath, "/"),
		Size:  info.Size(),
	}, nil
}
