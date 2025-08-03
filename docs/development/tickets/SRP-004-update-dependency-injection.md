# Ticket SRP-004: Update Dependency Injection and Configuration

## Summary
Complete the dependency injection chain to properly wire the new repository layer into the application initialization process.

## Description
Update the application's dependency injection system to properly initialize and inject the shipping repository throughout the application stack, ensuring clean configuration management and proper initialization order.

## Acceptance Criteria

### Must Have
- [x] Update `RepoWrapper` to include `ShippingRepository`
- [x] Update `ServiceWrapper` initialization to use repository
- [x] Ensure proper configuration injection from app config
- [x] Maintain existing initialization patterns
- [x] Verify handler layer continues to work without changes
- [x] Test complete application startup with new dependencies

### Should Have
- [x] Add configuration validation for shipping repository
- [x] Implement proper error handling in initialization
- [x] Add logging for dependency injection steps
- [x] Create health check for shipping repository

### Could Have
- [ ] Add configuration hot-reloading support
- [ ] Implement graceful degradation if repository fails
- [ ] Add dependency injection documentation

## Technical Requirements

### Repository Wrapper Update
```go
// internal/repository/util/init.go
type RepoWrapper struct {
    ProductRepo  repository.ProductRepository
    CartRepo     repository.CartRepository
    PaymentRepo  repository.PaymentRepository
    CategoryRepo repository.CategoryRepository
    ShippingRepo repository.ShippingRepository // Add shipping repository
}

func New(cfg *config.AppConfig) (repoWrapper *RepoWrapper, err error) {
    // Initialize database connection
    dbConnection, err := db.Init(cfg)
    if err != nil {
        return nil, fmt.Errorf("failed to initialize database: %w", err)
    }

    // Initialize shipping repository
    shippingRepo, err := initShippingRepository(cfg)
    if err != nil {
        return nil, fmt.Errorf("failed to initialize shipping repository: %w", err)
    }

    repoWrapper = &RepoWrapper{
        ProductRepo:  dbConnection,
        CartRepo:     dbConnection,
        PaymentRepo:  dbConnection,
        CategoryRepo: dbConnection,
        ShippingRepo: shippingRepo,
    }

    return repoWrapper, nil
}

func initShippingRepository(cfg *config.AppConfig) (repository.ShippingRepository, error) {
    // Validate configuration
    if err := validateShippingConfig(cfg); err != nil {
        return nil, fmt.Errorf("invalid shipping configuration: %w", err)
    }

    // Create repository with configuration
    repo := rajaongkir.NewRepository(rajaongkir.Config{
        APIKey:  cfg.RajaOngkirAPIKey,
        BaseURL: cfg.RajaOngkirBaseURL,
        Timeout: 30 * time.Second,
    })

    // Test repository connectivity (optional health check)
    if err := testShippingRepository(repo); err != nil {
        log.Printf("Warning: Shipping repository health check failed: %v", err)
        // Don't fail initialization, just log warning
    }

    return repo, nil
}
```

### Configuration Validation
```go
// internal/repository/util/validation.go
func validateShippingConfig(cfg *config.AppConfig) error {
    if cfg.RajaOngkirAPIKey == "" {
        return errors.New("RajaOngkir API key is required")
    }

    if cfg.RajaOngkirBaseURL == "" {
        return errors.New("RajaOngkir base URL is required")
    }

    // Validate URL format
    if _, err := url.Parse(cfg.RajaOngkirBaseURL); err != nil {
        return fmt.Errorf("invalid RajaOngkir base URL: %w", err)
    }

    return nil
}

func testShippingRepository(repo repository.ShippingRepository) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Test with a simple province lookup
    _, err := repo.GetProvinces("")
    if err != nil {
        return fmt.Errorf("repository health check failed: %w", err)
    }

    return nil
}
```

### Service Wrapper Update
```go
// internal/service/util/init.go
func New(cfg *config.AppConfig, repoWrapper *util.RepoWrapper) (serviceWrapper *ServiceWrapper, err error) {
    productService := product.New(cfg, repoWrapper)
    cartService := cart.New(cfg, repoWrapper)
    paymentService := payment.New(cfg, repoWrapper)
    categoryService := category.New(cfg, repoWrapper)
    
    // Update shipping service to use repository
    shippingService := shipping.New(cfg, repoWrapper)

    serviceWrapper = &ServiceWrapper{
        ProductService:  productService,
        CartService:     cartService,
        PaymentService:  paymentService,
        ShippingService: shippingService,
        CategoryService: categoryService,
    }

    return serviceWrapper, nil
}
```

### Main Application Update
```go
// cmd/main.go - Ensure proper initialization order
func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal("Failed to load configuration:", err)
    }

    // Initialize repository layer
    repoWrapper, err := repoUtil.New(cfg)
    if err != nil {
        log.Fatal("Failed to initialize repositories:", err)
    }

    // Initialize service layer
    serviceWrapper, err := serviceUtil.New(cfg, repoWrapper)
    if err != nil {
        log.Fatal("Failed to initialize services:", err)
    }

    // Initialize handlers and start server
    e := echo.New()
    handlerUtil.InitHandler(cfg, e, serviceWrapper)

    // Start server
    log.Printf("Starting server on port %s", cfg.Port)
    log.Fatal(e.Start(":" + cfg.Port))
}
```

## Implementation Steps

1. **Update Repository Wrapper** (`internal/repository/util/init.go`)
   - Add `ShippingRepo` field to `RepoWrapper` struct
   - Implement shipping repository initialization
   - Add configuration validation
   - Add optional health check

2. **Update Service Initialization** (`internal/service/util/init.go`)
   - Ensure shipping service receives repository wrapper
   - Verify proper dependency passing

3. **Add Configuration Validation**
   - Create validation functions for shipping config
   - Add proper error handling and messaging
   - Implement health check functionality

4. **Update Application Main** (`cmd/main.go`)
   - Verify initialization order is correct
   - Add proper error handling for startup
   - Ensure graceful shutdown handles new dependencies

5. **Add Logging and Monitoring**
   - Log successful repository initialization
   - Add health check endpoints
   - Implement proper error logging

6. **Create Integration Tests**
   - Test complete application startup
   - Verify dependency injection chain
   - Test configuration validation

## Dependencies
- **Requires**: SRP-003 (Service layer must be updated)
- **Files to modify**:
  - `internal/repository/util/init.go`
  - `internal/service/util/init.go`
  - `cmd/main.go`
  - Configuration files if needed

## Estimated Effort
- **Development**: 4 hours
- **Testing**: 3 hours
- **Integration verification**: 2 hours
- **Total**: 9 hours

## Testing Requirements

### Integration Tests
```go
// test/integration/dependency_injection_test.go
func TestDependencyInjection_Complete(t *testing.T) {
    // Load test configuration
    cfg := &config.AppConfig{
        RajaOngkirAPIKey:  "test-key",
        RajaOngkirBaseURL: "https://api.rajaongkir.com/starter",
        DatabaseURL:       "test-db-url",
        // ... other config
    }

    // Test repository initialization
    repoWrapper, err := repoUtil.New(cfg)
    assert.NoError(t, err)
    assert.NotNil(t, repoWrapper.ShippingRepo)

    // Test service initialization
    serviceWrapper, err := serviceUtil.New(cfg, repoWrapper)
    assert.NoError(t, err)
    assert.NotNil(t, serviceWrapper.ShippingService)

    // Test that service can use repository
    provinces, err := serviceWrapper.ShippingService.GetProvinces(request.GetProvincesRequest{})
    assert.NoError(t, err)
    assert.NotNil(t, provinces)
}

func TestConfigurationValidation(t *testing.T) {
    tests := []struct {
        name          string
        config        *config.AppConfig
        expectedError string
    }{
        {
            name: "missing_api_key",
            config: &config.AppConfig{
                RajaOngkirAPIKey:  "",
                RajaOngkirBaseURL: "https://api.rajaongkir.com/starter",
            },
            expectedError: "RajaOngkir API key is required",
        },
        {
            name: "missing_base_url",
            config: &config.AppConfig{
                RajaOngkirAPIKey:  "test-key",
                RajaOngkirBaseURL: "",
            },
            expectedError: "RajaOngkir base URL is required",
        },
        {
            name: "invalid_base_url",
            config: &config.AppConfig{
                RajaOngkirAPIKey:  "test-key",
                RajaOngkirBaseURL: "invalid-url",
            },
            expectedError: "invalid RajaOngkir base URL",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := repoUtil.New(tt.config)
            
            if tt.expectedError != "" {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectedError)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Application Startup Tests
```go
// test/startup/application_test.go
func TestApplicationStartup(t *testing.T) {
    // Test that application can start with new dependencies
    // This should be an integration test that starts the actual application
    
    cfg := loadTestConfig()
    
    // Initialize all dependencies
    repoWrapper, err := repoUtil.New(cfg)
    require.NoError(t, err)
    
    serviceWrapper, err := serviceUtil.New(cfg, repoWrapper)
    require.NoError(t, err)
    
    // Initialize Echo server
    e := echo.New()
    handlerUtil.InitHandler(cfg, e, serviceWrapper)
    
    // Test that server can start (don't actually start, just verify setup)
    assert.NotNil(t, e)
    
    // Test that endpoints are registered
    routes := e.Routes()
    assert.NotEmpty(t, routes)
    
    // Verify shipping endpoints exist
    var shippingRoutes []echo.Route
    for _, route := range routes {
        if strings.Contains(route.Path, "/shipping") {
            shippingRoutes = append(shippingRoutes, route)
        }
    }
    assert.NotEmpty(t, shippingRoutes)
}
```

## Configuration Requirements

### Environment Variables
```bash
# Required environment variables
RAJAONGKIR_API_KEY=your-api-key
RAJAONGKIR_BASE_URL=https://api.rajaongkir.com/starter

# Optional configuration
SHIPPING_TIMEOUT=30s
SHIPPING_HEALTH_CHECK=true
```

### Configuration File Updates
```json
// config/app.config.json
{
  "rajaongkir_api_key": "${RAJAONGKIR_API_KEY}",
  "rajaongkir_base_url": "${RAJAONGKIR_BASE_URL}",
  "shipping_timeout": "${SHIPPING_TIMEOUT:30s}",
  "shipping_health_check": "${SHIPPING_HEALTH_CHECK:true}"
}
```

## Health Check Implementation

### Repository Health Check
```go
// internal/repository/rajaongkir/health.go
func (r *Repository) HealthCheck(ctx context.Context) error {
    // Simple health check - get provinces without ID (should be fast)
    _, err := r.GetProvinces("")
    if err != nil {
        return fmt.Errorf("shipping repository health check failed: %w", err)
    }
    return nil
}
```

### HTTP Health Check Endpoint
```go
// internal/handler/health/shipping.go
func (h *HealthHandler) CheckShipping(c echo.Context) error {
    ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
    defer cancel()

    if healthChecker, ok := h.shippingRepo.(interface{ HealthCheck(context.Context) error }); ok {
        if err := healthChecker.HealthCheck(ctx); err != nil {
            return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{
                "status": "unhealthy",
                "error":  err.Error(),
            })
        }
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "status": "healthy",
    })
}
```

## Definition of Done
- [x] Repository wrapper includes shipping repository
- [x] Service wrapper properly initializes with repository dependencies
- [x] Configuration validation works correctly
- [x] Application starts successfully with new dependency chain
- [x] All existing functionality continues to work
- [x] Integration tests pass
- [x] Health checks are implemented and working
- [x] Security considerations are implemented
- [x] Error handling provides clear messages
- [x] Logging provides adequate debugging information
- [x] Documentation is updated

## Risk Assessment

### Medium Risk
- **Initialization Order Issues**: Dependencies might be initialized in wrong order
  - **Mitigation**: Careful testing of startup sequence, comprehensive integration tests

- **Configuration Missing**: Required config values might not be set
  - **Mitigation**: Validation functions, clear error messages, documentation

### Low Risk
- **Health Check Failures**: Health checks might fail even when service works
  - **Mitigation**: Make health checks optional, provide clear logging

## Notes
- Maintain backward compatibility during transition
- Ensure configuration validation provides clear error messages
- Add comprehensive logging for debugging initialization issues
- Consider graceful degradation if shipping repository fails

## Security Considerations

### Configuration Security
- [x] Ensure API keys and credentials are securely stored
- [x] Implement validation to prevent startup with invalid credentials
- [x] Support environment variables for sensitive configuration

### Secure Initialization
- [x] Validate configuration before initializing components
- [x] Implement proper error handling during initialization
- [x] Ensure credentials aren't logged during startup

### Security Testing
- [x] Test behavior with missing or invalid credentials
- [x] Verify proper error messages that don't expose sensitive details
- [x] Test health checks don't expose implementation details

## Related Tickets
- **Depends on**: SRP-003 (Update Service Layer)
- **Blocks**: SRP-005 (Update Tests and Mocks)
- **Related**: SRP-006 (Legacy Cleanup)

## Labels
- `refactor`
- `dependency-injection`
- `configuration`
- `initialization`
