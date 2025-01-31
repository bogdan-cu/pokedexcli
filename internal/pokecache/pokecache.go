package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry
	mu      *sync.RWMutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	entries := make(map[string]cacheEntry)
	cache := Cache{entries: entries, mu: &sync.RWMutex{}}
	go cache.ReapLoop(interval)
	return &cache
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := cacheEntry{createdAt: time.Now(), val: value}
	c.entries[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	return entry.val, ok
}

func (c *Cache) ReapLoop(interval time.Duration) {
	timer := time.NewTicker(interval)
	defer timer.Stop()
	for {
		t := <-timer.C
		c.mu.RLock()
		for key, entry := range c.entries {
			if t.Add(-interval).After(entry.createdAt) {
				c.mu.RUnlock()
				c.mu.Lock()
				delete(c.entries, key)
				c.mu.Unlock()
				c.mu.RLock()
			}
		}
		c.mu.RUnlock()
	}
}
