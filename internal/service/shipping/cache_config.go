package shipping

import "time"

// CacheConfig holds configuration for the RajaOngkir API cache
type CacheConfig struct {
	Enabled         bool          `json:"enabled"`
	DefaultTTL      time.Duration `json:"default_ttl"`
	WarmupOnStartup bool          `json:"warmup_on_startup"`
	WarmupTimeout   time.Duration `json:"warmup_timeout"`
}

// DefaultCacheConfig returns a default cache configuration
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		Enabled:         true,
		DefaultTTL:      24 * time.Hour, // 24 hours default TTL
		WarmupOnStartup: true,
		WarmupTimeout:   30 * time.Second,
	}
}
