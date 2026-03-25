package cache

import (
	"sync"
	"time"
)

// entry 单条缓存条目
type entry struct {
	value     any
	expiresAt time.Time
}

func (e *entry) expired() bool {
	return time.Now().After(e.expiresAt)
}

// Cache 基于内存的 TTL 缓存，并发安全
type Cache struct {
	mu  sync.RWMutex
	ttl time.Duration
	m   map[string]*entry
}

// New 创建缓存实例，ttl 为缓存有效期
func New(ttl time.Duration) *Cache {
	return &Cache{
		ttl: ttl,
		m:   make(map[string]*entry),
	}
}

// Get 获取缓存值，不存在或已过期返回 nil, false
func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	e, ok := c.m[key]
	c.mu.RUnlock()

	if !ok || e.expired() {
		return nil, false
	}
	return e.value, true
}

// Set 写入缓存
func (c *Cache) Set(key string, value any) {
	c.mu.Lock()
	c.m[key] = &entry{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
	c.mu.Unlock()
}

// Delete 删除指定缓存
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.m, key)
	c.mu.Unlock()
}

// Clear 清空所有缓存
func (c *Cache) Clear() {
	c.mu.Lock()
	c.m = make(map[string]*entry)
	c.mu.Unlock()
}
