# Ticket SRP-002: Move RajaOngkir Implementation to Repository Layer

## Summary
Move the `RajaOngkirClient` implementation from the service layer to the repository layer and refactor it to implement the `ShippingRepository` interface.

## Description
Migrate the existing `RajaOngkirClient` from `internal/service/shipping/` to `internal/repository/rajaongkir/` and refactor it to properly implement the repository pattern while maintaining all existing functionality. Since the service is not yet released, we can perform a direct migration without maintaining compatibility layers.

## Acceptance Criteria

### Must Have
- [x] Move `RajaOngkirClient` code to `internal/repository/rajaongkir/client.go`
- [x] Refactor struct to implement `ShippingRepository` interface
- [x] Maintain all existing HTTP client functionality
- [x] Preserve API key security and configuration patterns
- [x] Ensure proper error handling and wrapping
- [x] Maintain existing timeout and retry behavior
- [x] Remove old client files from service layer immediately

### Should Have
- [x] Improve error messages with context information
- [x] Add structured logging for debugging
- [x] Implement proper request/response validation
- [x] Add metrics collection preparation
- [x] Clean up unused imports and dependencies

### Could Have
- [ ] Optimize HTTP client configuration
- [ ] Prepare hooks for caching layer
- [ ] Add request tracing capabilities

## Technical Requirements

### Implementation Structure
```go
// internal/repository/rajaongkir/client.go
package rajaongkir

type Repository struct {
    apiKey  string
    baseURL string
    client  *http.Client
    logger  *log.Logger
}

func NewRepository(cfg Config, opts ...Option) *Repository {
    // Implementation with proper dependency injection
}

func (r *Repository) GetProvinces(provinceID string) ([]response.RajaOngkirProvince, error) {
    // Migrated implementation with improved error handling
}
```

### Error Handling Enhancement
```go
func (r *Repository) handleAPIError(operation string, resp *http.Response, err error) error {
    if err != nil {
        return &RepositoryError{
            Operation: operation,
            Cause:     err,
            Details: map[string]interface{}{
                "api_endpoint": r.baseURL,
                "status_code":  resp.StatusCode,
            },
        }
    }
    return nil
}
```

### Configuration Integration
```go
// Ensure proper configuration injection from app config
func NewRepositoryFromAppConfig(cfg *config.AppConfig) *Repository {
    return NewRepository(Config{
        APIKey:  cfg.RajaOngkirAPIKey,
        BaseURL: cfg.RajaOngkirBaseURL,
        Timeout: 30 * time.Second,
    })
}
```

## Implementation Steps

1. **Create repository implementation** (`internal/repository/rajaongkir/client.go`)
   - Copy existing `RajaOngkirClient` code
   - Rename struct to `Repository`
   - Implement `ShippingRepository` interface
   - Add proper constructor function

2. **Enhance error handling**
   - Wrap errors with operation context
   - Add structured error types
   - Improve error messages

3. **Update HTTP client configuration**
   - Review and optimize timeout settings
   - Ensure proper TLS configuration
   - Add connection pooling if needed

4. **Add logging and observability**
   - Add structured logging for debugging
   - Prepare metrics collection points
   - Add request/response tracing

5. **Implement configuration patterns**
   - Create flexible configuration options
   - Add validation for required settings
   - Support environment-specific settings

6. **Create comprehensive tests**
   - Unit tests for all methods
   - HTTP client mock tests
   - Error scenario testing
   - Timeout and retry testing

## Dependencies
- **Requires**: SRP-001 (Repository interface must exist)
- **Files to modify**: 
  - `internal/service/shipping/rajaongkir_client.go` (source)
  - `internal/repository/rajaongkir/client.go` (destination)

## Estimated Effort
- **Development**: 8 hours
- **Testing**: 6 hours
- **Code review and refinement**: 2 hours
- **Total**: 16 hours

## Testing Requirements

### Unit Tests Coverage
- [x] All repository methods (GetProvinces, GetCities, GetDistricts, CalculateShippingCost)
- [x] HTTP client configuration and initialization
- [x] Error handling scenarios
- [x] Timeout and network failure scenarios
- [x] Invalid response parsing scenarios

### Mock Testing Strategy
```go
func TestRepository_GetProvinces(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Mock different API responses
    }))
    defer server.Close()

    repo := NewRepository(Config{
        APIKey:  "test-key",
        BaseURL: server.URL,
        Timeout: 5 * time.Second,
    })

    // Test various scenarios
}
```

## Security Considerations

### API Key Protection
- [x] Ensure API keys are not logged or exposed in error messages
- [x] Validate API key format and presence
- [x] Use secure HTTP client configuration

### Input Validation
- [x] Validate all input parameters before API calls
- [x] Sanitize user input to prevent injection attacks
- [x] Implement proper URL encoding

### Error Information Leakage
- [x] Ensure internal details are not exposed in error messages
- [x] Log sensitive information securely
- [x] Return user-appropriate error messages

## Performance Considerations

### HTTP Client Optimization
- [x] Use connection pooling for better performance
- [x] Implement proper timeout configurations
- [x] Add retry logic with exponential backoff

### Memory Management
- [x] Ensure proper resource cleanup
- [x] Optimize JSON parsing and memory allocation
- [x] Implement streaming for large responses if needed

## Definition of Done
- [x] Repository implementation compiles and runs correctly
- [x] All existing functionality is preserved
- [x] Repository implements the ShippingRepository interface
- [x] Comprehensive unit tests are written and passing
- [x] Error handling is improved and consistent
- [x] Security requirements are met
- [x] Performance meets or exceeds current benchmarks
- [x] Code follows project style guidelines
- [x] Documentation is updated

## Risk Assessment

### High Risk
- **API Integration Breaking**: Moving HTTP client code might break external API calls
  - **Mitigation**: Comprehensive testing with real API calls, feature flags for rollback

### Medium Risk
- **Configuration Issues**: API keys might not be properly injected
  - **Mitigation**: Extensive configuration testing, validation functions

### Low Risk
- **Performance Regression**: Additional abstraction might impact performance
  - **Mitigation**: Performance benchmarking, optimization as needed

## Notes
- Preserve all existing HTTP client behavior exactly
- Focus on improving structure and maintainability
- Prepare foundation for future caching implementation
- Maintain backward compatibility during transition

## Related Tickets
- **Depends on**: SRP-001 (Create Repository Interface)
- **Blocks**: SRP-003 (Update Service Layer)
- **Related**: SRP-005 (Update Tests and Mocks)

## Labels
- `refactor`
- `repository`
- `migration`
- `http-client`
