package pokecache

import (
	"sync"
	"time"
)

// New cache
func NewCache(interval time.Duration) *Cache {
	newCache := Cache{
		interval:   interval,
		CachedData: make(map[string]cacheEntry),
	}

	// Peroidically deletes expired cache
	go newCache.reapLoop()

	return &newCache
}

// Cache
type Cache struct {
	CachedData map[string]cacheEntry
	Mu         sync.RWMutex
	interval   time.Duration
}

// Cache Entry
type cacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)

	for range ticker.C {
		c.Mu.Lock()
		// Remove cache
		for key, cache := range c.CachedData {
			// Calculate time between now and created cache
			timeSince := time.Since(cache.CreatedAt)

			// Remove the expired cached data
			if timeSince > c.interval {
				delete(c.CachedData, key)
			}
		}
		c.Mu.Unlock()
	}
}

// Add
func (c *Cache) Add(key string, val []byte) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	entry := cacheEntry{
		CreatedAt: time.Now(),
		Val:       val,
	}

	c.CachedData[key] = entry
}

// Get
func (c *Cache) Get(key string) ([]byte, bool) {
	c.Mu.RLock()
	defer c.Mu.RUnlock()

	entry, ok := c.CachedData[key]

	return entry.Val, ok
}
