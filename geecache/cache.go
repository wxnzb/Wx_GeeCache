package geecache

import (
	"fmt"
	"geecache/lru"
	"sync"
)

type cache struct {
	mu      sync.Mutex
	ca      *lru.Cache
	maxSize int
}

// 将lru中的lru.go实现加锁功能
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ca == nil {
		c.ca = lru.NewCache(c.maxSize, nil)
	}
	c.ca.Add(key, value)
}
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c == nil {
		return
	}
	fmt.Printf("%p", c)
	if v, ok := c.ca.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
