package pokecache

import (
	"sync"
	"time"
)

// Package pokecache provides caching logic for the Pokedex application.

// entry represents a cache value and the time it was added.
type entry struct {
	val     []byte
	addedAt time.Time
}

// Cache is a thread-safe cache with expiration.
type Cache struct {
	mu       sync.Mutex
	entries  map[string]entry
	interval time.Duration
}

// NewCache creates a new Cache and starts the reap loop.
func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries:  make(map[string]entry),
		interval: interval,
	}
	go c.reapLoop()
	return c
}

// Add adds a new entry to the cache.
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = entry{val: val, addedAt: time.Now()}
}

// Get retrieves an entry from the cache.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	return e.val, true
}

// reapLoop periodically removes entries older than the interval.
func (c *Cache) reapLoop() {
	for {
		time.Sleep(c.interval)
		now := time.Now()
		c.mu.Lock()
		for k, e := range c.entries {
			if now.Sub(e.addedAt) > c.interval {
				delete(c.entries, k)
			}
		}
		c.mu.Unlock()
	}
}
