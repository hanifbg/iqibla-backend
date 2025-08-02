# RajaOngkir Caching Documentation

## Overview

This document provides information about the RajaOngkir API caching implementation, which helps reduce API usage and stay within the 100 requests/day limit.

## Features

- In-memory caching of province data with configurable TTL
- Thread-safe cache access for concurrent requests
- Cache warm-up on application startup
- Automatic fallback to API when cache is empty or expired
- Comprehensive logging of cache operations
- Seamless integration with existing code via Decorator pattern

## Configuration

The following configuration options are available in the `app.config.json` file under the `shipping` section:

```json
"shipping": {
    "rajaongkir_api_key": "your_api_key",
    "rajaongkir_base_url": "https://rajaongkir.komerce.id/api/v1",
    "rajaongkir_cache_enabled": true,
    "rajaongkir_cache_ttl_hours": 24,
    "rajaongkir_warmup_on_startup": true,
    "rajaongkir_warmup_timeout_secs": 30
}
```

### Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `rajaongkir_cache_enabled` | Enable/disable caching | `true` |
| `rajaongkir_cache_ttl_hours` | Cache time-to-live in hours | `24` |
| `rajaongkir_warmup_on_startup` | Pre-populate cache on application startup | `true` |
| `rajaongkir_warmup_timeout_secs` | Timeout for cache warm-up in seconds | `30` |

## How It Works

1. When the application starts, it creates a `CachedRajaOngkirClient` that wraps the regular `RajaOngkirClient`
2. If warm-up is enabled, it pre-fetches province data to populate the cache
3. When a client makes a request for province data:
   - If the cache contains valid (non-expired) data, it's returned immediately
   - If the cache is empty or expired, the client falls back to the API
   - After a successful API call, the cache is updated with fresh data
4. The cache automatically expires after the configured TTL (default: 24 hours)

## Logging

Cache operations are logged with the prefix `[RajaOngkir Cache]` and include:
- Cache hit/miss events
- Warm-up success/failure
- Cache initialization status

## Troubleshooting

### Cache not working

1. Verify that `rajaongkir_cache_enabled` is set to `true` in your configuration
2. Check logs for any errors during cache initialization
3. Ensure the API key is valid and the service is reachable

### Slow startup

If application startup is slow due to cache warm-up:
1. Adjust `rajaongkir_warmup_timeout_secs` to a lower value
2. Set `rajaongkir_warmup_on_startup` to `false` to disable warm-up

### Stale data

If you're seeing outdated province data:
1. Decrease `rajaongkir_cache_ttl_hours` to refresh data more frequently
2. Restart the application to force a cache refresh

## Future Enhancements

Currently, only province data is cached. Future enhancements will include:
- Caching of city data
- Caching of district data
- Background cache refreshing
- More comprehensive metrics and monitoring
