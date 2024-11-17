package cache

import (
	"sync"
	"time"
)

type HTMLCache struct {
	content    string
	lastUpdate time.Time
	ttl        time.Duration
	mu         sync.RWMutex
}

func NewHTMLCache(ttl time.Duration) *HTMLCache {
	return &HTMLCache{
		ttl: ttl,
	}
}

func (c *HTMLCache) Get() (string, bool) {
	if c.content == "" {
		return "", false
	}

	if time.Since(c.lastUpdate) > c.ttl {
		return "", false
	}

	return c.content, true
}

func (c *HTMLCache) Set(content string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.content = content
	c.lastUpdate = time.Now()
}
