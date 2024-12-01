package pokeapi

import (
	"sync"
	"time"
)

type Cache struct {
	Entries  map[string]cacheEntry
	Mutex    *sync.Mutex
	Interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	m := make(map[string]cacheEntry)
	mutex := sync.Mutex{}

	cache := Cache{
		Entries:  m,
		Mutex:    &mutex,
		Interval: interval,
	}

	go cache.reapLoop()

	return &cache
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.Interval)
	for range ticker.C {
		go c.reap()
	}
}

func (c *Cache) reap() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	for key, entry := range c.Entries {
		expiry := time.Now().UTC().Add(-c.Interval)
		if entry.createdAt.Before(expiry) || entry.createdAt.Equal(expiry) {
			delete(c.Entries, key)
		}
	}
}

func (c *Cache) Add(key string, val []byte) bool {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	return true
}

func (c *Cache) Get(key string) ([]byte, bool) {
	var val []byte
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	entry, ok := c.Entries[key]
	if !ok {
		return nil, false
	}
	val = entry.val
	return val, true
}

func (c *Cache) GetAll() [][]byte {
	var res [][]byte
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	for key := range c.Entries {
		if val, ok := c.Get(key); ok {
			res = append(res, val)
		}
	}
	return res
}
