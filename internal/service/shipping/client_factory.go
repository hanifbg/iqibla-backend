package shipping

import (
	"log"
	"os"
	"time"
)

// CreateRajaOngkirClient creates a RajaOngkirClient instance
// If caching is enabled, it returns a CachedRajaOngkirClient wrapping the regular client
func CreateRajaOngkirClient(apiKey, baseURL string, config CacheConfig) RajaOngkirClientInterface {
	// Create the base client
	baseClient := NewRajaOngkirClient(apiKey, baseURL)

	// If caching is disabled, return the base client
	if !config.Enabled {
		return baseClient
	}

	// Create a logger for the cached client
	logger := log.New(os.Stdout, "[RajaOngkir Cache] ", log.LstdFlags)

	// Create the cached client
	cachedClient := NewCachedRajaOngkirClient(baseClient, config.DefaultTTL, logger)

	// If warm-up on startup is enabled, initialize the cache
	if config.WarmupOnStartup {
		// Create a channel to handle timeout
		done := make(chan bool)
		var err error

		// Initialize cache in a goroutine
		go func() {
			err = cachedClient.InitCache()
			done <- true
		}()

		// Wait for initialization or timeout
		select {
		case <-done:
			if err != nil {
				logger.Printf("Cache warm-up failed: %v", err)
			}
		case <-time.After(config.WarmupTimeout):
			logger.Printf("Cache warm-up timed out after %v", config.WarmupTimeout)
		}
	}

	return cachedClient
}
