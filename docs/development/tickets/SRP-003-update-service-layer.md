# Ticket SRP-003: Update Service Layer to Use Repository Interface

## Summary
Refactor the `ShippingService` to use the new `ShippingRepository` interface instead of directly depending on the `RajaOngkirClient`.

## Description
Update the shipping s### Backward Compatibility

### Service Interface Preservation
- [x] All existing method signatures remain unchanged
- [x] Response formats stay identical
- [x] Error handling behavior is consistent
- [x] Performance characteristics are maintained

### Handler Layer Impact
- [x] No changes required to handler layer
- [x] All existing API endpoints continue to work
- [x] Response formats remain identical to follow clean architecture principles by depending on the repository interface abstraction rather than the concrete RajaOngkir implementation. Since the service is not yet released, we can make breaking changes to achieve cleaner architecture without backward compatibility concerns.

## Acceptance Criteria

### Must Have
- [x] Update `ShippingService` struct to use `ShippingRepository` interface
- [x] Modify service initialization to inject repository dependency
- [x] Remove all references to `RajaOngkirClient` from service layer
- [x] Preserve all business logic and data transformation
- [x] Ensure proper error handling and propagation
- [x] Update service constructor and dependency injection
- [x] Clean up unused imports and dependencies

### Should Have
- [x] Improve error messages with business context
- [x] Add service-level input validation
- [x] Implement proper logging for business operations
- [x] Add performance monitoring hooks
- [x] Simplify service methods without compatibility layers

### Could Have
- [ ] Add request caching at service level
- [ ] Implement circuit breaker patterns for external dependencies
- [ ] Add business metrics collection

## Technical Requirements

### Service Structure Update
```go
// internal/service/shipping/init.go
type ShippingService struct {
    shippingRepo repository.ShippingRepository
}

func New(cfg *config.AppConfig, repoWrapper *util.RepoWrapper) service.ShippingService {
    return &ShippingService{
        shippingRepo: repoWrapper.ShippingRepo,
    }
}
```

### Method Implementation Example
```go
// internal/service/shipping/impl.go
func (s *ShippingService) GetProvinces(req request.GetProvincesRequest) ([]response.ProvinceResponse, error) {
    // Add service-level validation
    if err := s.validateGetProvincesRequest(req); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }

    // Call repository layer
    provinces, err := s.shippingRepo.GetProvinces(req.ID)
    if err != nil {
        // Add business context to error
        return nil, fmt.Errorf("failed to retrieve province data: %w", err)
    }

    // Business logic: transform repository response to service response
    var result []response.ProvinceResponse
    for _, province := range provinces {
        result = append(result, response.ProvinceResponse{
            ProvinceID: strconv.Itoa(province.ProvinceID),
            Province:   province.Province,
        })
    }

    return result, nil
}
```

### Dependency Injection Update
```go
// internal/repository/util/init.go - Update RepoWrapper
type RepoWrapper struct {
    ProductRepo  repository.ProductRepository
    CartRepo     repository.CartRepository
    PaymentRepo  repository.PaymentRepository
    CategoryRepo repository.CategoryRepository
    ShippingRepo repository.ShippingRepository // Add this
}

func New(cfg *config.AppConfig) (repoWrapper *RepoWrapper, err error) {
    // ... existing code ...

    // Initialize shipping repository
    shippingRepo := rajaongkir.NewRepository(rajaongkir.Config{
        APIKey:  cfg.RajaOngkirAPIKey,
        BaseURL: cfg.RajaOngkirBaseURL,
        Timeout: 30 * time.Second,
    })

    repoWrapper = &RepoWrapper{
        ProductRepo:  dbConnection,
        CartRepo:     dbConnection,
        PaymentRepo:  dbConnection,
        CategoryRepo: dbConnection,
        ShippingRepo: shippingRepo, // Add this
    }

    return
}
```

## Implementation Steps

1. **Update service struct** (`internal/service/shipping/init.go`)
   - Replace `rajaOngkirClient` field with `shippingRepo` field
   - Update field type to use repository interface
   - Modify constructor to accept repository from wrapper

2. **Update service methods** (`internal/service/shipping/impl.go`)
   - Change all method calls from client to repository
   - Maintain existing business logic
   - Improve error handling with business context
   - Add input validation where appropriate

3. **Update dependency injection** (`internal/repository/util/init.go`)
   - Add `ShippingRepo` field to `RepoWrapper`
   - Initialize shipping repository in constructor
   - Ensure proper configuration injection

4. **Update service initialization** (`internal/service/util/init.go`)
   - Pass repository wrapper to shipping service constructor
   - Ensure proper dependency chain

5. **Add validation and error handling**
   - Implement request validation functions
   - Add business-appropriate error messages
   - Implement proper logging

6. **Update imports and remove old dependencies**
   - Remove direct RajaOngkir client imports
   - Add repository interface imports
   - Clean up unused code

## Dependencies
- **Requires**: SRP-002 (Repository implementation must exist)
- **Files to modify**:
  - `internal/service/shipping/init.go`
  - `internal/service/shipping/impl.go`
  - `internal/repository/util/init.go`
  - `internal/service/util/init.go`

## Estimated Effort
- **Development**: 6 hours
- **Testing**: 4 hours
- **Integration testing**: 2 hours
- **Total**: 12 hours

## Testing Requirements

### Service Layer Tests Update
```go
// internal/service/shipping/impl_test.go
func TestShippingService_GetProvinces(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockShippingRepository(ctrl)
    service := &ShippingService{shippingRepo: mockRepo}

    tests := []struct {
        name           string
        request        request.GetProvincesRequest
        mockSetup      func()
        expectedResult []response.ProvinceResponse
        expectedError  string
    }{
        {
            name:    "successful_transformation",
            request: request.GetProvincesRequest{ID: "1"},
            mockSetup: func() {
                mockRepo.EXPECT().
                    GetProvinces("1").
                    Return([]response.RajaOngkirProvince{
                        {ProvinceID: 1, Province: "Bali"},
                    }, nil)
            },
            expectedResult: []response.ProvinceResponse{
                {ProvinceID: "1", Province: "Bali"},
            },
            expectedError: "",
        },
        {
            name:    "repository_error_handling",
            request: request.GetProvincesRequest{ID: "999"},
            mockSetup: func() {
                mockRepo.EXPECT().
                    GetProvinces("999").
                    Return(nil, errors.New("repository error"))
            },
            expectedResult: nil,
            expectedError:  "failed to retrieve province data",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mockSetup()

            result, err := service.GetProvinces(tt.request)

            if tt.expectedError != "" {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectedError)
                assert.Nil(t, result)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expectedResult, result)
            }
        })
    }
}
```

### Integration Tests
- [x] Test complete dependency injection chain
- [x] Verify handler → service → repository flow
- [x] Test error propagation through all layers
- [x] Validate business logic preservation

## Validation Requirements

### Input Validation Examples
```go
func (s *ShippingService) validateGetProvincesRequest(req request.GetProvincesRequest) error {
    // Add business-level validation if needed
    // Example: validate province ID format
    if req.ID != "" {
        if _, err := strconv.Atoi(req.ID); err != nil {
            return fmt.Errorf("province ID must be numeric: %s", req.ID)
        }
    }
    return nil
}

func (s *ShippingService) validateCalculateShippingRequest(req request.CalculateShippingRequest) error {
    if req.Origin == "" {
        return errors.New("origin is required")
    }
    if req.Destination == "" {
        return errors.New("destination is required")
    }
    if req.Weight <= 0 {
        return errors.New("weight must be positive")
    }
    return nil
}
```

## Backward Compatibility

### Service Interface Preservation
- [ ] All existing method signatures remain unchanged
- [ ] Response formats stay identical
- [ ] Error handling behavior is consistent
- [ ] Performance characteristics are maintained

### Handler Layer Impact
- [ ] No changes required to handler layer
- [ ] All existing API endpoints continue to work
- [ ] Response formats remain identical

## Definition of Done
- [x] Service layer uses repository interface instead of concrete client
- [x] All service methods are updated and tested
- [x] Dependency injection chain is complete and working
- [x] All existing functionality is preserved
- [x] Service-level tests are updated to use mocks
- [x] Integration tests pass
- [x] No breaking changes to external interfaces
- [x] Error handling is improved with business context
- [x] Security considerations are implemented
- [x] Code follows project style guidelines
- [x] Documentation is updated

## Risk Assessment

### Medium Risk
- **Dependency Injection Issues**: Repository might not be properly injected
  - **Mitigation**: Comprehensive integration testing, step-by-step validation

- **Business Logic Changes**: Accidentally modifying transformation logic
  - **Mitigation**: Careful code review, comprehensive test coverage

### Low Risk
- **Interface Mismatch**: Service expecting different interface methods
  - **Mitigation**: Compile-time verification, thorough testing

## Notes
- Focus on maintaining exact same business behavior
- Service layer should contain business logic, not data access logic
- Prepare service layer for future enhancements (caching, multiple providers)
- Ensure clean separation between service and repository concerns

## Security Considerations

### Business-Layer Validation
- [x] Implement comprehensive input validation at service layer
- [x] Validate business rules before passing data to repository
- [x] Prevent business logic bypass through invalid inputs

### Error Handling and Information Leakage
- [x] Ensure service errors don't expose implementation details
- [x] Add business context to errors without leaking sensitive data
- [x] Implement proper error wrapping for debugging without compromising security

### Security Boundaries
- [x] Maintain clear security boundaries between layers
- [x] Ensure service layer doesn't expose repository implementation details
- [x] Sanitize data crossing layer boundaries

## Related Tickets
- **Depends on**: SRP-002 (Move RajaOngkir Implementation)
- **Blocks**: SRP-004 (Update Dependency Injection)
- **Related**: SRP-005 (Update Tests and Mocks)

## Labels
- `refactor`
- `service-layer`
- `dependency-injection`
- `business-logic`
