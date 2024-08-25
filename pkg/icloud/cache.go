package icloud

import (
	"sync"
	"time"
)

type CacheInterface interface {
	Set(key string, value interface{}, duration time.Duration)
	Get(key string) (interface{}, bool)
}

type Cache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

var _ CacheInterface = &Cache{}

type cacheItem struct {
	value      interface{}
	expiration int64
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]cacheItem),
	}
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(duration).UnixNano(),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[key]
	if !found || time.Now().UnixNano() > item.expiration {
		return nil, false
	}
	return item.value, true
}
