package cache

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type CacheEntry struct {
	Content    string
	ModTime    time.Time
	Size       int64
	IsTextFile *bool
	LoadTime   time.Time
}

type FileCache struct {
	mu       sync.RWMutex
	entries  map[string]*CacheEntry
	usage    map[string]time.Time
	maxSize  int64
	maxItems int
	curSize  int64
}

func NewFileCache(maxSizeBytes int64, maxItems int) *FileCache {
	return &FileCache{
		entries:  make(map[string]*CacheEntry),
		usage:    make(map[string]time.Time),
		maxSize:  maxSizeBytes,
		maxItems: maxItems,
		curSize:  0,
	}
}

func DefaultFileCache() *FileCache {
	return NewFileCache(100*1024*1024, 10000)
}

func (fc *FileCache) Get(filePath string) (string, error) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	if entry, exists := fc.entries[filePath]; exists {
		if stat, err := os.Stat(filePath); err == nil {
			if stat.ModTime().Equal(entry.ModTime) && stat.Size() == entry.Size {
				if entry.Content != "" {
					fc.usage[filePath] = time.Now()
					return entry.Content, nil
				}
				return fc.loadContentForExistingEntry(filePath, entry)
			}
		}
		fc.removeEntry(filePath)
	}

	return fc.loadAndCache(filePath)
}

func (fc *FileCache) GetTextFileStatus(filePath string, checkFunc func(string) (bool, error)) (bool, error) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	if entry, exists := fc.entries[filePath]; exists {
		if entry.IsTextFile != nil {
			if stat, err := os.Stat(filePath); err == nil {
				if stat.ModTime().Equal(entry.ModTime) && stat.Size() == entry.Size {
					fc.usage[filePath] = time.Now()
					return *entry.IsTextFile, nil
				}
			}
			fc.removeEntry(filePath)
		}
	}

	isText, err := checkFunc(filePath)
	if err != nil {
		return false, err
	}

	stat, statErr := os.Stat(filePath)
	if statErr != nil {
		return isText, nil
	}

	entry, exists := fc.entries[filePath]
	if !exists {
		entry = &CacheEntry{
			ModTime:  stat.ModTime(),
			Size:     stat.Size(),
			LoadTime: time.Now(),
		}
		fc.entries[filePath] = entry
	}
	entry.IsTextFile = &isText
	fc.usage[filePath] = time.Now()

	return isText, nil
}

func (fc *FileCache) Preload(filePath string) error {
	_, err := fc.Get(filePath)
	return err
}

func (fc *FileCache) Contains(filePath string) bool {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	entry, exists := fc.entries[filePath]
	if !exists {
		return false
	}

	if stat, err := os.Stat(filePath); err == nil {
		return stat.ModTime().Equal(entry.ModTime) && stat.Size() == entry.Size
	}

	return false
}

func (fc *FileCache) Clear() {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	fc.entries = make(map[string]*CacheEntry)
	fc.usage = make(map[string]time.Time)
	fc.curSize = 0
}

func (fc *FileCache) Stats() (items int, sizeBytes int64, hitRate float64) {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	return len(fc.entries), fc.curSize, 0.0
}

func (fc *FileCache) loadAndCache(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	stat, err := os.Stat(filePath)
	if err != nil {
		return string(content), nil
	}

	contentStr := string(content)
	contentSize := int64(len(contentStr))

	fc.evictIfNeeded(contentSize)

	entry := &CacheEntry{
		Content:  contentStr,
		ModTime:  stat.ModTime(),
		Size:     stat.Size(),
		LoadTime: time.Now(),
	}
	fc.entries[filePath] = entry
	fc.usage[filePath] = time.Now()
	fc.curSize += contentSize

	return contentStr, nil
}

func (fc *FileCache) evictIfNeeded(newItemSize int64) {
	for (fc.curSize+newItemSize > fc.maxSize || len(fc.entries) >= fc.maxItems) && len(fc.entries) > 0 {
		fc.evictLRU()
	}
}

func (fc *FileCache) evictLRU() {
	var oldestPath string
	var oldestTime time.Time
	first := true

	for path, accessTime := range fc.usage {
		if first || accessTime.Before(oldestTime) {
			oldestPath = path
			oldestTime = accessTime
			first = false
		}
	}

	if oldestPath != "" {
		fc.removeEntry(oldestPath)
	}
}

func (fc *FileCache) loadContentForExistingEntry(filePath string, entry *CacheEntry) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	contentStr := string(content)
	contentSize := int64(len(contentStr))

	fc.evictIfNeeded(contentSize)

	entry.Content = contentStr
	fc.usage[filePath] = time.Now()
	fc.curSize += contentSize

	return contentStr, nil
}

func (fc *FileCache) removeEntry(filePath string) {
	if entry, exists := fc.entries[filePath]; exists {
		fc.curSize -= int64(len(entry.Content))
		delete(fc.entries, filePath)
		delete(fc.usage, filePath)
	}
}

