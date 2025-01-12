package pokecache

import (
	"sync"
	"time"
)

// cacheEntry represents a single cache item
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// Cache manages a map of cache entries with a mutex for thread-safety
type Cache struct {
	mu       sync.Mutex
	entries  map[string]cacheEntry
	interval time.Duration
}

// NewCache creates a new cache with a specified cleanup interval
func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
	}

	// Start the reap loop in a goroutine
	go cache.reapLoop()

	return cache
}

// Add adds a new entry to the cache
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

// Get retrieves an entry from the cache
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, found := c.entries[key]
	if !found {
		return nil, false
	}
	return entry.val, true
}

// reapLoop periodically removes expired entries from the cache
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.reap()
	}
}

// reap removes entries older than the cache interval
func (c *Cache) reap() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.entries {
		if now.Sub(entry.createdAt) > c.interval {
			delete(c.entries, key)
		}
	}
}
