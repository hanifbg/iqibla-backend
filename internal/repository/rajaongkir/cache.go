package rajaongkir

import (
	"fmt"
	"sync"
	"time"

	"github.com/hanifbg/landing_backend/internal/model/response"
)

// CacheItem represents a cached item with expiry
type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

// Cache provides thread-safe caching with TTL
type Cache struct {
	mu    sync.RWMutex
	items map[string]*CacheItem
	ttl   time.Duration
}

// NewCache creates a new cache with the specified TTL
func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		items: make(map[string]*CacheItem),
		ttl:   ttl,
	}
}

// Get retrieves an item from cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.ExpiresAt) {
		// Item has expired
		delete(c.items, key)
		return nil, false
	}

	return item.Data, true
}

// Set stores an item in cache
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &CacheItem{
		Data:      value,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// GetProvinces retrieves provinces from cache
func (c *Cache) GetProvinces() ([]response.RajaOngkirProvince, bool) {
	if data, found := c.Get("provinces"); found {
		if provinces, ok := data.([]response.RajaOngkirProvince); ok {
			fmt.Printf("✅ CACHE HIT: Provinces loaded from cache\n")
			return provinces, true
		}
	}
	fmt.Printf("⚠️ CACHE MISS: Provinces not found in cache\n")
	return nil, false
}

// SetProvinces stores provinces in cache
func (c *Cache) SetProvinces(provinces []response.RajaOngkirProvince) {
	c.Set("provinces", provinces)
	fmt.Printf("💾 CACHE STORE: Stored %d provinces in cache\n", len(provinces))
}

// GetCities retrieves cities from cache
func (c *Cache) GetCities(provinceID string) ([]response.RajaOngkirCity, bool) {
	key := fmt.Sprintf("cities_%s", provinceID)
	if data, found := c.Get(key); found {
		if cities, ok := data.([]response.RajaOngkirCity); ok {
			fmt.Printf("✅ CACHE HIT: Cities for province %s loaded from cache\n", provinceID)
			return cities, true
		}
	}
	fmt.Printf("⚠️ CACHE MISS: Cities for province %s not found in cache\n", provinceID)
	return nil, false
}

// SetCities stores cities in cache
func (c *Cache) SetCities(provinceID string, cities []response.RajaOngkirCity) {
	key := fmt.Sprintf("cities_%s", provinceID)
	c.Set(key, cities)
	fmt.Printf("💾 CACHE STORE: Stored %d cities for province %s in cache\n", len(cities), provinceID)
}

// GetDistricts retrieves districts from cache
func (c *Cache) GetDistricts(cityID string) ([]response.RajaOngkirDistrict, bool) {
	key := fmt.Sprintf("districts_%s", cityID)
	if data, found := c.Get(key); found {
		if districts, ok := data.([]response.RajaOngkirDistrict); ok {
			fmt.Printf("✅ CACHE HIT: Districts for city %s loaded from cache\n", cityID)
			return districts, true
		}
	}
	fmt.Printf("⚠️ CACHE MISS: Districts for city %s not found in cache\n", cityID)
	return nil, false
}

// SetDistricts stores districts in cache
func (c *Cache) SetDistricts(cityID string, districts []response.RajaOngkirDistrict) {
	key := fmt.Sprintf("districts_%s", cityID)
	c.Set(key, districts)
	fmt.Printf("💾 CACHE STORE: Stored %d districts for city %s in cache\n", len(districts), cityID)
}
