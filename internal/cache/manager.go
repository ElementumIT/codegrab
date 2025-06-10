package cache

import (
	"sync"
)

type Manager struct {
	fileCache *FileCache
	once      sync.Once
}

var globalManager = &Manager{}

func (m *Manager) GetFileCache() *FileCache {
	m.once.Do(func() {
		m.fileCache = DefaultFileCache()
	})
	return m.fileCache
}

func Global() *Manager {
	return globalManager
}

func GetGlobalFileCache() *FileCache {
	return Global().GetFileCache()
}

func ResetGlobalCache() {
	globalManager = &Manager{}
}

