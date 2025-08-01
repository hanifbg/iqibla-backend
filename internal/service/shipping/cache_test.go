package shipping

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheEntry(t *testing.T) {
	// Test non-expired entry
	entry := CacheEntry{
		Data:      "test data",
		Timestamp: time.Now(),
		TTL:       1 * time.Hour,
	}
	assert.False(t, entry.IsExpired(), "New entry should not be expired")

	// Test expired entry
	expiredEntry := CacheEntry{
		Data:      "test data",
		Timestamp: time.Now().Add(-2 * time.Hour),
		TTL:       1 * time.Hour,
	}
	assert.True(t, expiredEntry.IsExpired(), "Entry older than TTL should be expired")
}

func TestCacheManager(t *testing.T) {
	// Setup
	cache := NewCacheManager(1 * time.Hour)
	testData := []string{"test1", "test2"}

	// Test set and get
	cache.SetProvinces(testData)
	data, found := cache.GetProvinces()
	assert.True(t, found, "Data should be found after being set")
	assert.Equal(t, testData, data, "Retrieved data should match set data")

	// Test clear
	cache.Clear()
	_, found = cache.GetProvinces()
	assert.False(t, found, "Data should not be found after cache is cleared")

	// Test expiration
	shortTTLCache := NewCacheManager(1 * time.Millisecond)
	shortTTLCache.SetProvinces(testData)
	time.Sleep(10 * time.Millisecond) // Wait for expiration
	_, found = shortTTLCache.GetProvinces()
	assert.False(t, found, "Data should not be found after TTL expires")
}
