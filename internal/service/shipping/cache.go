package shipping

import (
	"sync"
	"time"
)

// CacheEntry represents a single cached item with expiration
type CacheEntry struct {
	Data      interface{}
	Timestamp time.Time
	TTL       time.Duration
}

// IsExpired checks if the cache entry has expired
func (c *CacheEntry) IsExpired() bool {
	return time.Since(c.Timestamp) > c.TTL
}

// CacheManager manages the in-memory cache with TTL
type CacheManager struct {
	provinces  *CacheEntry
	mutex      sync.RWMutex
	defaultTTL time.Duration
}

// NewCacheManager creates a new cache manager with the specified default TTL
func NewCacheManager(defaultTTL time.Duration) *CacheManager {
	return &CacheManager{
		defaultTTL: defaultTTL,
	}
}

// SetProvinces stores provinces data in the cache
func (c *CacheManager) SetProvinces(data interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.provinces = &CacheEntry{
		Data:      data,
		Timestamp: time.Now(),
		TTL:       c.defaultTTL,
	}
}

// GetProvinces retrieves provinces data from the cache if it exists and hasn't expired
func (c *CacheManager) GetProvinces() (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.provinces == nil || c.provinces.IsExpired() {
		return nil, false
	}

	return c.provinces.Data, true
}

// Clear removes all items from the cache
func (c *CacheManager) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.provinces = nil
}
