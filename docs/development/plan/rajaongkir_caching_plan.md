# RajaOngkir API Caching Implementation

## Overview

Implement in-memory caching for RajaOngkir API responses to reduce API usage and work within the 100 requests/day limitation.

## Current Situation Analysis

### Problems Identified
1. **API Rate Limit**: RajaOngkir API has a strict 100 requests/day limitation
2. **Frequent Calls**: Geographic data (provinces, cities, districts) are requested frequently
3. **Static Data**: Geographic data rarely changes, making it ideal for caching
4. **No Caching**: Current implementation makes direct API calls for every request

### Current API Usage Patterns
- `GetProvinces()` - Called when users need to select their province
- `GetCities(provinceID)` - Called when users select a province and need cities  
- `GetDistricts(cityID)` - Called when users select a city and need districts
- `CalculateShippingCost()` - Called for shipping cost calculation (should not be cached due to dynamic pricing)

## Solution Design

### Architecture Overview

Implement a **Decorator Pattern** with a caching layer that wraps the existing `RajaOngkirClient`:

```
[Handler] → [CachedRajaOngkirClient] → [RajaOngkirClient] → [RajaOngkir API]
                    ↓
              [In-Memory Cache]
```

### Key Components

1. **CachedRajaOngkirClient**: Decorator that implements `RajaOngkirClientInterface`
2. **CacheManager**: Manages in-memory cache with TTL
3. **CacheEntry**: Wrapper for cached data with timestamp
4. **Cache Warm-up Service**: Pre-populates cache on application startup

## Implementation Plan

### Phase 1: Basic Province Caching (Week 1)

**Scope**: Implement caching for `GetProvinces()` method only

#### 1.1 Create Cache Infrastructure
- [ ] Create `CacheEntry` struct with data and timestamp
- [ ] Create `CacheManager` with thread-safe operations
- [ ] Implement TTL checking mechanism

#### 1.2 Implement CachedRajaOngkirClient
- [ ] Create decorator struct implementing `RajaOngkirClientInterface`
- [ ] Implement province caching logic
- [ ] Add cache hit/miss logging
- [ ] Implement fallback to API on cache miss

#### 1.3 Cache Initialization
- [ ] Create cache warm-up service
- [ ] Integrate with application startup
- [ ] Add graceful degradation if warm-up fails

### Phase 2: Extended Caching (Week 2)

**Scope**: Extend caching to cities and districts

#### 2.1 Cities Caching
- [ ] Implement cache key strategy for cities (by province)
- [ ] Add cities caching to `GetCities()` method
- [ ] Implement hierarchical cache warm-up (provinces → cities)

#### 2.2 Districts Caching  
- [ ] Implement cache key strategy for districts (by city)
- [ ] Add districts caching to `GetDistricts()` method
- [ ] Complete hierarchical cache warm-up (provinces → cities → districts)

#### 2.3 Configuration
- [ ] Add configuration for cache TTL
- [ ] Add configuration to enable/disable caching
- [ ] Add configuration for warm-up strategy

### Phase 3: Monitoring & Optimization (Week 3)

#### 3.1 Metrics & Monitoring
- [ ] Add cache hit/miss ratio metrics
- [ ] Add cache size metrics
- [ ] Add API call reduction metrics
- [ ] Add cache refresh timing metrics

#### 3.2 Error Handling & Resilience
- [ ] Implement circuit breaker pattern for API calls
- [ ] Add retry mechanism for failed cache refreshes
- [ ] Implement graceful degradation strategies

#### 3.3 Performance Optimization
- [ ] Implement background cache refresh
- [ ] Add memory usage monitoring
- [ ] Optimize cache key strategies

## Technical Specifications

### Cache Structure

```go
type CacheEntry struct {
    Data      interface{}
    Timestamp time.Time
    TTL       time.Duration
}

type CacheManager struct {
    provinces map[string]*CacheEntry
    cities    map[string]*CacheEntry // key: "province_id"
    districts map[string]*CacheEntry // key: "city_id"
    mutex     sync.RWMutex
    defaultTTL time.Duration
}

type CachedRajaOngkirClient struct {
    client RajaOngkirClientInterface
    cache  *CacheManager
    logger *log.Logger
}
```

### Cache Key Strategy
- **Provinces**: `"provinces"` (single entry for all provinces)
- **Cities**: `"cities:province_id:{provinceID}"` 
- **Districts**: `"districts:city_id:{cityID}"`

### Configuration Options
```go
type CacheConfig struct {
    Enabled           bool          `json:"enabled"`
    DefaultTTL        time.Duration `json:"default_ttl"`
    WarmupOnStartup   bool          `json:"warmup_on_startup"`
    WarmupTimeout     time.Duration `json:"warmup_timeout"`
    BackgroundRefresh bool          `json:"background_refresh"`
}
```

## Testing Strategy

### Unit Tests
- [ ] Cache hit/miss scenarios
- [ ] TTL expiration handling
- [ ] Concurrent access safety
- [ ] Error handling and fallbacks

### Integration Tests
- [ ] Cache warm-up process
- [ ] API fallback scenarios
- [ ] End-to-end shipping flow with caching

### Performance Tests
- [ ] Memory usage under load
- [ ] Cache performance vs API calls
- [ ] Concurrent access performance

## Risk Assessment

### High Priority Risks
1. **Memory Leaks**: Improper cache cleanup could lead to memory issues
   - *Mitigation*: Implement proper TTL cleanup and memory monitoring

2. **Stale Data**: Cached data might become outdated
   - *Mitigation*: Reasonable TTL (24 hours) and background refresh

3. **Race Conditions**: Concurrent cache access issues
   - *Mitigation*: Proper mutex usage and thread-safe operations

### Medium Priority Risks
1. **Startup Delays**: Cache warm-up might delay application startup
   - *Mitigation*: Asynchronous warm-up with fallback to API

2. **API Dependency**: Cache warm-up failure could impact startup
   - *Mitigation*: Graceful degradation and retry mechanisms

## Migration Strategy

### Backward Compatibility
- Implement decorator pattern to maintain existing interface
- No changes required in handlers or service layer
- Configuration-based enabling of cache

### Deployment Plan
1. **Stage 1**: Deploy with cache disabled by default
2. **Stage 2**: Enable cache in staging environment
3. **Stage 3**: Monitor performance and API usage reduction
4. **Stage 4**: Enable cache in production with monitoring

## Success Metrics

### Primary Metrics
- **API Usage Reduction**: Target 80-90% reduction in province/city/district API calls
- **Response Time**: Maintain or improve current response times
- **Cache Hit Ratio**: Target >90% hit ratio after warm-up

### Secondary Metrics
- **Memory Usage**: Monitor additional memory consumption
- **Error Rate**: Maintain current error rates
- **Startup Time**: Acceptable increase in startup time (<30 seconds)

## Future Enhancements

### Phase 4: Persistent Caching (Future)
- Implement Redis-based caching for multi-instance deployments
- Add database persistence for cache data
- Implement distributed cache invalidation

### Phase 5: Smart Caching (Future)
- Implement predictive cache loading based on usage patterns
- Add cache analytics and optimization recommendations
- Implement dynamic TTL based on data update frequency

## Dependencies

### External Dependencies
- No new external dependencies required for Phase 1-3
- Future Redis integration will require Redis client library

### Internal Dependencies
- Configuration system for cache settings
- Logging framework for cache metrics
- Existing RajaOngkir client implementation

## Rollback Plan

1. **Immediate Rollback**: Disable cache via configuration
2. **Code Rollback**: Revert to direct API client usage
3. **Cleanup**: Remove cache infrastructure if needed

The decorator pattern ensures easy rollback without breaking existing functionality.
