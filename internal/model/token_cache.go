package model

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/epilande/codegrab/internal/utils"
)

// TokenResult represents a cached token count result
type TokenResult struct {
	Tokens    int
	Error     error
	Timestamp time.Time
}

// TokenCache manages async token estimation with caching
type TokenCache struct {
	cache     map[string]TokenResult
	mutex     sync.RWMutex
	workQueue chan string
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

// NewTokenCache creates a new token cache with background workers
func NewTokenCache() *TokenCache {
	ctx, cancel := context.WithCancel(context.Background())
	tc := &TokenCache{
		cache:     make(map[string]TokenResult),
		workQueue: make(chan string, 100),
		ctx:       ctx,
		cancel:    cancel,
	}
	
	numWorkers := 2
	for i := 0; i < numWorkers; i++ {
		tc.wg.Add(1)
		go tc.worker()
	}
	
	return tc
}

// GetTokens returns cached token count or queues for background calculation
func (tc *TokenCache) GetTokens(filePath string) (int, bool) {
	tc.mutex.RLock()
	result, exists := tc.cache[filePath]
	tc.mutex.RUnlock()
	
	if exists {
		if result.Error != nil {
			return 0, true
		}
		return result.Tokens, true
	}
	
	select {
	case tc.workQueue <- filePath:
	default:
	}
	
	return 0, false
}

// GetTokensFormatted returns formatted token string or empty if not cached
func (tc *TokenCache) GetTokensFormatted(filePath string) string {
	if tokens, cached := tc.GetTokens(filePath); cached && tokens > 0 {
		return fmt.Sprintf(" [%d tokens]", tokens)
	}
	return ""
}

// InvalidateFile removes a file from the cache
func (tc *TokenCache) InvalidateFile(filePath string) {
	tc.mutex.Lock()
	delete(tc.cache, filePath)
	tc.mutex.Unlock()
}

// ClearCache removes all cached entries
func (tc *TokenCache) ClearCache() {
	tc.mutex.Lock()
	tc.cache = make(map[string]TokenResult)
	tc.mutex.Unlock()
}

// Close shuts down the token cache and background workers
func (tc *TokenCache) Close() {
	tc.cancel()
	close(tc.workQueue)
	tc.wg.Wait()
}

// worker processes token calculation requests in background
func (tc *TokenCache) worker() {
	defer tc.wg.Done()
	
	for {
		select {
		case filePath, ok := <-tc.workQueue:
			if !ok {
				return
			}
			tc.calculateTokens(filePath)
			
		case <-tc.ctx.Done():
			return
		}
	}
}

// calculateTokens performs the actual token calculation and caching
func (tc *TokenCache) calculateTokens(filePath string) {
	tc.mutex.RLock()
	_, exists := tc.cache[filePath]
	tc.mutex.RUnlock()
	
	if exists {
		return
	}
	
	result := TokenResult{
		Timestamp: time.Now(),
	}
	
	if ok, err := utils.IsTextFile(filePath); !ok || err != nil {
		result.Error = fmt.Errorf("not a text file or error checking: %v", err)
	} else {
		if contentBytes, err := os.ReadFile(filePath); err != nil {
			result.Error = err
		} else {
			content := string(contentBytes)
			result.Tokens = utils.EstimateTokens(content)
		}
	}
	
	tc.mutex.Lock()
	tc.cache[filePath] = result
	tc.mutex.Unlock()
}

// Stats returns cache statistics for debugging
func (tc *TokenCache) Stats() (cached int, queued int) {
	tc.mutex.RLock()
	cached = len(tc.cache)
	tc.mutex.RUnlock()
	
	queued = len(tc.workQueue)
	return
}