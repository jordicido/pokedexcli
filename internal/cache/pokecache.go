package Cache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]CacheEntry
	mutex sync.RWMutex
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	newCache := &Cache{cache: make(map[string]CacheEntry), mutex: sync.RWMutex{}}
	go newCache.reapLoop(interval)
	return newCache
}

func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache[key] = CacheEntry{val: val, createdAt: time.Now()}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		<-ticker.C
		c.mutex.Lock()

		c.cache = make(map[string]CacheEntry)

		c.mutex.Unlock()
	}
}
