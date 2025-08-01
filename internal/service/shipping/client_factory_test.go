package shipping

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultCacheConfig(t *testing.T) {
	config := DefaultCacheConfig()

	assert.True(t, config.Enabled, "Default config should have caching enabled")
	assert.Equal(t, 24*time.Hour, config.DefaultTTL, "Default TTL should be 24 hours")
	assert.True(t, config.WarmupOnStartup, "Default config should warm up cache on startup")
	assert.Equal(t, 30*time.Second, config.WarmupTimeout, "Default warmup timeout should be 30 seconds")
}

func TestCreateRajaOngkirClient(t *testing.T) {
	// Test with caching disabled
	config := CacheConfig{
		Enabled: false,
	}

	client := CreateRajaOngkirClient("test-key", "http://test-url", config)
	_, isRegularClient := client.(*RajaOngkirClient)
	assert.True(t, isRegularClient, "Should return regular client when caching is disabled")

	// Test with caching enabled but no warm-up
	config = CacheConfig{
		Enabled:         true,
		DefaultTTL:      1 * time.Hour,
		WarmupOnStartup: false,
	}

	client = CreateRajaOngkirClient("test-key", "http://test-url", config)
	cachedClient, isCachedClient := client.(*CachedRajaOngkirClient)
	assert.True(t, isCachedClient, "Should return cached client when caching is enabled")
	assert.False(t, cachedClient.initialized, "Should not initialize cache when warm-up is disabled")

	// Test with caching and warm-up enabled, but with a very short timeout to avoid actual API calls
	// This is more of an integration test and might be flaky, so we're just checking type
	config = CacheConfig{
		Enabled:         true,
		DefaultTTL:      1 * time.Hour,
		WarmupOnStartup: true,
		WarmupTimeout:   1 * time.Millisecond, // Very short timeout to avoid actual API calls
	}

	client = CreateRajaOngkirClient("test-key", "http://test-url", config)
	_, isCachedClient = client.(*CachedRajaOngkirClient)
	assert.True(t, isCachedClient, "Should return cached client with warm-up enabled")
}
