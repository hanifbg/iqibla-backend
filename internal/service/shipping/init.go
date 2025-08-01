package shipping

import (
	"log"
	"os"
	"time"

	"github.com/hanifbg/landing_backend/config"
	"github.com/hanifbg/landing_backend/internal/service"
)

type ShippingService struct {
	rajaOngkirClient RajaOngkirClientInterface
}

func New(cfg *config.AppConfig) service.ShippingService {
	// Create cache configuration from app config
	cacheConfig := CacheConfig{
		Enabled:         cfg.RajaOngkirCacheEnabled,
		DefaultTTL:      time.Duration(cfg.RajaOngkirCacheTTLHours) * time.Hour,
		WarmupOnStartup: cfg.RajaOngkirWarmupOnStartup,
		WarmupTimeout:   time.Duration(cfg.RajaOngkirWarmupTimeoutSecs) * time.Second,
	}

	// If config values are zero, use defaults
	if cacheConfig.DefaultTTL == 0 {
		cacheConfig.DefaultTTL = 24 * time.Hour
	}
	if cacheConfig.WarmupTimeout == 0 {
		cacheConfig.WarmupTimeout = 30 * time.Second
	}

	// Create a logger for the cached client
	logger := log.New(os.Stdout, "[RajaOngkir Cache] ", log.LstdFlags)

	// Create the RajaOngkir client (either cached or regular based on config)
	rajaOngkirClient := CreateRajaOngkirClient(cfg.RajaOngkirAPIKey, cfg.RajaOngkirBaseURL, cacheConfig)

	// Log cache status
	if cacheConfig.Enabled {
		logger.Printf("RajaOngkir cache enabled with TTL of %v", cacheConfig.DefaultTTL)
		if cacheConfig.WarmupOnStartup {
			logger.Printf("Cache warm-up on startup enabled with timeout of %v", cacheConfig.WarmupTimeout)
		}
	} else {
		logger.Printf("RajaOngkir cache disabled")
	}

	return &ShippingService{
		rajaOngkirClient: rajaOngkirClient,
	}
}
