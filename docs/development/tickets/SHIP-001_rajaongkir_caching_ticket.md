# Tracking Ticket: RajaOngkir API Caching Implementation

**Ticket ID**: SHIP-001  
**Type**: Enhancement  
**Priority**: High  
**Epic**: Shipping Optimization  
**Created**: August 1, 2025  
**Assignee**: Backend Team  

## User Story

**As a** backend service  
**I want** to cache RajaOngkir API responses in memory  
**So that** I can reduce API usage and stay within the 100 requests/day limit while maintaining fast response times  

## Background

RajaOngkir API has a strict limitation of 100 requests per day. Our current implementation makes direct API calls for every request to get provinces, cities, and districts data. Since geographic data is relatively static, implementing caching will significantly reduce API usage while improving performance.

## Acceptance Criteria

### Must Have (Phase 1)
- [x] **Cache provinces data** in memory with 24-hour TTL
- [x] **Thread-safe cache access** using proper synchronization
- [x] **Cache warm-up on startup** to pre-populate province data
- [x] **Fallback to API** when cache is empty or expired
- [x] **Maintain existing API interface** (no breaking changes)
- [x] **Configurable cache settings** (TTL, enable/disable)
- [x] **Error handling** for cache warm-up failures
- [x] **Logging** for cache hit/miss events

### Should Have (Phase 2) 
- [ ] **Extend caching to cities** data (by province)
- [ ] **Extend caching to districts** data (by city)
- [ ] **Hierarchical cache warm-up** (provinces → cities → districts)
- [ ] **Background cache refresh** to prevent expiration
- [ ] **Cache metrics** for monitoring hit/miss ratios

### Could Have (Phase 3)
- [ ] **Memory usage monitoring** and alerts
- [ ] **Circuit breaker pattern** for API resilience
- [ ] **Cache size limits** to prevent memory overflow
- [ ] **Performance benchmarks** comparing cached vs non-cached

### Won't Have (This Release)
- Persistent caching (Redis/Database)
- Distributed cache synchronization
- Cache warming based on usage patterns
- Dynamic TTL adjustment

## Technical Requirements

### Architecture
- Implement **Decorator Pattern** wrapping existing `RajaOngkirClient`
- Maintain `RajaOngkirClientInterface` compatibility
- Use **in-memory map** with proper synchronization (sync.RWMutex)
- Implement **TTL-based expiration** mechanism

### Configuration
```json
{
  "rajaongkir_cache": {
    "enabled": true,
    "default_ttl": "24h",
    "warmup_on_startup": true,
    "warmup_timeout": "30s",
    "background_refresh": true
  }
}
```

### Cache Structure
- **Provinces**: Single cache entry for all provinces
- **Cities**: One cache entry per province (key: `province_id`)
- **Districts**: One cache entry per city (key: `city_id`)

### Error Handling
- **Startup**: If cache warm-up fails, continue without cache (graceful degradation)
- **Runtime**: If cache read fails, fallback to direct API call
- **Refresh**: If background refresh fails, log error and continue with existing cache

## Implementation Tasks

### Phase 1: Basic Province Caching (Sprint 1)

#### Backend Implementation
- [x] **Create cache infrastructure** 
  - [x] `CacheEntry` struct with data and timestamp
  - [x] `CacheManager` with thread-safe operations
  - [x] TTL validation methods

- [x] **Implement CachedRajaOngkirClient**
  - [x] Decorator struct implementing `RajaOngkirClientInterface`
  - [x] Province caching logic in `GetProvinces()`
  - [x] Cache hit/miss logging
  - [x] API fallback mechanism

- [x] **Cache initialization service**
  - [x] Startup warm-up process
  - [x] Configuration integration
  - [x] Error handling and graceful degradation

- [x] **Testing**
  - [x] Unit tests for cache operations
  - [x] Integration tests for API fallback
  - [x] Concurrent access tests
  - [x] Error scenario tests

#### Documentation
- [x] Update API documentation
- [x] Add configuration guide
- [x] Create troubleshooting guide

### Phase 2: Extended Caching (Sprint 2)

- [ ] **Cities caching implementation**
  - [ ] Cache key strategy
  - [ ] Extend warm-up process
  - [ ] Update `GetCities()` method

- [ ] **Districts caching implementation**
  - [ ] Cache key strategy  
  - [ ] Complete hierarchical warm-up
  - [ ] Update `GetDistricts()` method

- [ ] **Background refresh mechanism**
  - [ ] Periodic cache refresh
  - [ ] Error handling for refresh failures
  - [ ] Configurable refresh intervals

### Phase 3: Monitoring & Optimization (Sprint 3)

- [ ] **Metrics implementation**
  - [ ] Cache hit/miss ratios
  - [ ] Memory usage tracking
  - [ ] API call reduction metrics

- [ ] **Performance optimization**
  - [ ] Memory usage optimization
  - [ ] Cache key optimization
  - [ ] Response time improvements

## Testing Strategy

### Unit Tests
- [x] Cache hit scenarios
- [x] Cache miss scenarios  
- [x] TTL expiration handling
- [x] Concurrent access safety
- [x] Error handling paths

### Integration Tests
- [x] End-to-end shipping flow
- [x] Cache warm-up process
- [x] API fallback scenarios
- [x] Configuration changes

### Performance Tests
- [ ] Memory usage under load
- [ ] Response time comparison
- [ ] Concurrent user simulation
- [ ] Cache performance benchmarks

## Definition of Done

### Code Quality
- [x] All unit tests passing (>90% coverage)
- [x] Integration tests passing
- [ ] Code review completed and approved
- [x] No performance regression
- [x] Memory usage within acceptable limits

### Documentation
- [x] Technical documentation updated
- [x] Configuration guide created
- [x] Deployment instructions updated
- [x] Troubleshooting guide available

### Deployment
- [ ] Successfully deployed to staging
- [ ] Performance monitoring configured
- [ ] Cache metrics available
- [ ] Error monitoring in place
- [ ] Rollback plan tested

## Risks and Mitigation

### High Priority
1. **Memory Leaks**
   - *Risk*: Improper cache cleanup leading to memory issues
   - *Mitigation*: Implement proper TTL cleanup and memory monitoring
   - *Owner*: Backend Team

2. **Race Conditions** 
   - *Risk*: Concurrent cache access causing data corruption
   - *Mitigation*: Proper mutex usage and thread-safe operations
   - *Owner*: Backend Team

### Medium Priority
1. **Startup Delays**
   - *Risk*: Cache warm-up delaying application startup
   - *Mitigation*: Asynchronous warm-up with timeout
   - *Owner*: DevOps Team

2. **Stale Data**
   - *Risk*: Cached data becoming outdated
   - *Mitigation*: Reasonable TTL and background refresh
   - *Owner*: Backend Team

## Success Metrics

### Primary KPIs
- **API Usage Reduction**: Target 80-90% reduction in geographic API calls
- **Response Time**: Maintain current response times (<200ms for cached data)
- **Cache Hit Ratio**: Achieve >90% hit ratio after warm-up period

### Secondary KPIs  
- **Memory Usage**: Additional memory consumption <100MB
- **Error Rate**: Maintain current error rates (<0.1%)
- **Startup Time**: Increase in startup time <30 seconds

## Dependencies

### Internal
- [x] Configuration system updates
- [x] Logging framework integration
- [ ] Monitoring system updates

### External
- None for Phase 1-3

## Timeline

| Phase | Duration | Deliverables |
|-------|----------|--------------|
| Phase 1 | Week 1 | Province caching implementation |
| Phase 2 | Week 2 | Cities and districts caching |
| Phase 3 | Week 3 | Monitoring and optimization |
| Testing | Week 4 | End-to-end testing and deployment |

**Total Estimated Effort**: 4 weeks (1 developer)

## Rollback Plan

1. **Configuration Rollback**: Disable cache via configuration flag
2. **Code Rollback**: Revert to previous version if needed
3. **Data Cleanup**: Clear cache if memory issues occur

## Related Tickets

- SHIP-002: Performance monitoring for shipping APIs
- SHIP-003: Redis-based persistent caching (Future enhancement)
- CONFIG-001: Configuration management improvements

## Notes

- Keep `CalculateShippingCost()` uncached as shipping costs are dynamic
- Consider expanding caching to other static data in future iterations
- Monitor API usage carefully during initial rollout
- Ensure proper error logging for debugging cache issues
