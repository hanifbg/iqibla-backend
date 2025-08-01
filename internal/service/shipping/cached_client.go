package shipping

import (
	"fmt"
	"log"
	"time"

	"github.com/hanifbg/landing_backend/internal/model/response"
)

// CachedRajaOngkirClient is a decorator that wraps the RajaOngkirClientInterface
// and implements caching for API responses
type CachedRajaOngkirClient struct {
	client      RajaOngkirClientInterface
	cache       *CacheManager
	logger      *log.Logger
	initialized bool
}

// NewCachedRajaOngkirClient creates a new cached client that wraps the original client
func NewCachedRajaOngkirClient(client RajaOngkirClientInterface, ttl time.Duration, logger *log.Logger) *CachedRajaOngkirClient {
	return &CachedRajaOngkirClient{
		client:      client,
		cache:       NewCacheManager(ttl),
		logger:      logger,
		initialized: false,
	}
}

// Ensure CachedRajaOngkirClient implements RajaOngkirClientInterface
var _ RajaOngkirClientInterface = (*CachedRajaOngkirClient)(nil)

// InitCache warms up the cache with initial data
func (c *CachedRajaOngkirClient) InitCache() error {
	if c.initialized {
		return nil
	}

	// Warm up the provinces cache
	c.logger.Printf("Warming up provinces cache...")
	provinces, err := c.client.GetProvinces("")
	if err != nil {
		c.logger.Printf("Failed to warm up provinces cache: %v", err)
		return fmt.Errorf("failed to warm up provinces cache: %w", err)
	}

	c.cache.SetProvinces(provinces)
	c.logger.Printf("Provinces cache warmed up successfully with %d provinces", len(provinces))
	c.initialized = true

	return nil
}

// GetProvinces retrieves a list of provinces, using cache when available
func (c *CachedRajaOngkirClient) GetProvinces(provinceID string) ([]response.RajaOngkirProvince, error) {
	// If specific province is requested, bypass cache
	if provinceID != "" {
		c.logger.Printf("Cache bypassed for specific province ID: %s", provinceID)
		return c.client.GetProvinces(provinceID)
	}

	// Try to get from cache first
	if data, found := c.cache.GetProvinces(); found {
		c.logger.Printf("Cache hit: returning provinces from cache")
		if provinces, ok := data.([]response.RajaOngkirProvince); ok {
			return provinces, nil
		}
		// If type assertion fails, log and continue to API call
		c.logger.Printf("Cache error: type assertion failed for provinces")
	} else {
		c.logger.Printf("Cache miss: fetching provinces from API")
	}

	// Not in cache or cache expired, call the API
	provinces, err := c.client.GetProvinces(provinceID)
	if err != nil {
		return nil, err
	}

	// Update cache with new data
	c.cache.SetProvinces(provinces)
	c.logger.Printf("Cache updated with %d provinces", len(provinces))

	return provinces, nil
}

// GetCities simply delegates to the underlying client (caching will be added in Phase 2)
func (c *CachedRajaOngkirClient) GetCities(provinceID, cityID string) ([]response.RajaOngkirCity, error) {
	return c.client.GetCities(provinceID, cityID)
}

// GetDistricts simply delegates to the underlying client (caching will be added in Phase 2)
func (c *CachedRajaOngkirClient) GetDistricts(cityID string) ([]response.RajaOngkirDistrict, error) {
	return c.client.GetDistricts(cityID)
}

// CalculateShippingCost delegates to the underlying client (should not be cached as shipping costs are dynamic)
func (c *CachedRajaOngkirClient) CalculateShippingCost(origin, destination string, weight int, courier string) ([]response.RajaOngkirCost, error) {
	return c.client.CalculateShippingCost(origin, destination, weight, courier)
}
