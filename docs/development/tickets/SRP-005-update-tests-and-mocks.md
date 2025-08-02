# Ticket SRP-005: Update Tests and Mock Generation

## Summary
Update all existing tests to use the new repository interface pattern and regenerate mocks to maintain >85% test coverage throughout the refactoring process.

## Description
Comprehensive update of the test suite to use the new repository pattern, including updating existing tests, creating new repository-specific tests, and ensuring proper mock generation for all interfaces.

## Acceptance Criteria

### Must Have
- [x] Update all service layer tests to use repository mocks
- [x] Create comprehensive repository layer unit tests
- [x] Regenerate all mock interfaces with correct paths
- [x] Maintain test coverage above 85%
- [x] Ensure all existing test scenarios are preserved
- [x] Update integration tests to use new architecture

### Should Have
- [x] Improve test coverage for edge cases and error scenarios
- [x] Add performance benchmarks for repository layer
- [x] Create test utilities for common scenarios
- [x] Add comprehensive validation testing

### Could Have
- [ ] Add mutation testing to verify test quality
- [ ] Create load testing for repository layer
- [ ] Add chaos engineering tests for external API failures

## Technical Requirements

### Mock Generation Updates
```go
// internal/repository/shipping.go
//go:generate mockgen -source=shipping.go -destination=../service/shipping/mocks/shipping_repository_mock.go -package=mocks

// internal/service/shipping.go (if needed)
//go:generate mockgen -source=shipping.go -destination=mocks/shipping_service_mock.go -package=mocks
```

### Repository Layer Tests
```go
// internal/repository/rajaongkir/client_test.go
package rajaongkir

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/hanifbg/landing_backend/internal/model/response"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestRepository_GetProvinces(t *testing.T) {
    tests := []struct {
        name           string
        provinceID     string
        mockResponse   interface{}
        mockStatusCode int
        mockDelay      time.Duration
        expectedResult []response.RajaOngkirProvince
        expectedError  string
    }{
        {
            name:           "successful_all_provinces",
            provinceID:     "",
            mockStatusCode: 200,
            mockResponse: response.KomerceProvinceResponse{
                Meta: response.Meta{Code: 200, Status: "success"},
                Data: []response.RajaOngkirProvince{
                    {ProvinceID: 1, Province: "Bali"},
                    {ProvinceID: 2, Province: "Bangka Belitung"},
                },
            },
            expectedResult: []response.RajaOngkirProvince{
                {ProvinceID: 1, Province: "Bali"},
                {ProvinceID: 2, Province: "Bangka Belitung"},
            },
            expectedError: "",
        },
        {
            name:           "successful_specific_province",
            provinceID:     "1",
            mockStatusCode: 200,
            mockResponse: response.KomerceProvinceResponse{
                Meta: response.Meta{Code: 200, Status: "success"},
                Data: []response.RajaOngkirProvince{
                    {ProvinceID: 1, Province: "Bali"},
                },
            },
            expectedResult: []response.RajaOngkirProvince{
                {ProvinceID: 1, Province: "Bali"},
            },
            expectedError: "",
        },
        {
            name:           "api_error_response",
            provinceID:     "999",
            mockStatusCode: 200,
            mockResponse: response.KomerceProvinceResponse{
                Meta: response.Meta{Code: 400, Status: "error", Message: "Invalid province ID"},
                Data: []response.RajaOngkirProvince{},
            },
            expectedResult: nil,
            expectedError:  "API error: Invalid province ID",
        },
        {
            name:           "network_error",
            provinceID:     "1",
            mockStatusCode: 500,
            mockResponse:   "Internal Server Error",
            expectedResult: nil,
            expectedError:  "failed to make request",
        },
        {
            name:           "timeout_error",
            provinceID:     "1",
            mockStatusCode: 200,
            mockDelay:      100 * time.Millisecond,
            mockResponse: response.KomerceProvinceResponse{
                Meta: response.Meta{Code: 200, Status: "success"},
                Data: []response.RajaOngkirProvince{{ProvinceID: 1, Province: "Bali"}},
            },
            expectedResult: nil,
            expectedError:  "failed to make request",
        },
        {
            name:           "invalid_json_response",
            provinceID:     "1",
            mockStatusCode: 200,
            mockResponse:   "invalid json response",
            expectedResult: nil,
            expectedError:  "failed to unmarshal response",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                // Verify request headers
                assert.Equal(t, "test-api-key", r.Header.Get("key"))
                
                // Verify request path
                expectedPath := "/destination/province"
                if tt.provinceID != "" {
                    expectedPath += "/" + tt.provinceID
                }
                assert.Equal(t, expectedPath, r.URL.Path)

                // Simulate delay for timeout tests
                if tt.mockDelay > 0 {
                    time.Sleep(tt.mockDelay)
                }

                w.WriteHeader(tt.mockStatusCode)

                if resp, ok := tt.mockResponse.(response.KomerceProvinceResponse); ok {
                    json.NewEncoder(w).Encode(resp)
                } else {
                    w.Write([]byte(tt.mockResponse.(string)))
                }
            }))
            defer server.Close()

            // Create repository with short timeout for timeout tests
            timeout := 30 * time.Second
            if tt.mockDelay > 0 {
                timeout = 50 * time.Millisecond
            }

            repo := NewRepository(Config{
                APIKey:  "test-api-key",
                BaseURL: server.URL,
                Timeout: timeout,
            })

            result, err := repo.GetProvinces(tt.provinceID)

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

// Similar comprehensive tests for GetCities, GetDistricts, CalculateShippingCost
```

### Service Layer Tests Update
```go
// internal/service/shipping/impl_test.go
package shipping

import (
    "errors"
    "testing"

    "github.com/golang/mock/gomock"
    "github.com/hanifbg/landing_backend/internal/model/request"
    "github.com/hanifbg/landing_backend/internal/model/response"
    "github.com/hanifbg/landing_backend/internal/service/shipping/mocks"
    "github.com/stretchr/testify/assert"
)

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
            name:    "successful_province_transformation",
            request: request.GetProvincesRequest{ID: "1"},
            mockSetup: func() {
                mockRepo.EXPECT().
                    GetProvinces("1").
                    Return([]response.RajaOngkirProvince{
                        {ProvinceID: 1, Province: "Bali"},
                        {ProvinceID: 2, Province: "Jakarta"},
                    }, nil)
            },
            expectedResult: []response.ProvinceResponse{
                {ProvinceID: "1", Province: "Bali"},
                {ProvinceID: "2", Province: "Jakarta"},
            },
            expectedError: "",
        },
        {
            name:    "empty_province_list",
            request: request.GetProvincesRequest{ID: ""},
            mockSetup: func() {
                mockRepo.EXPECT().
                    GetProvinces("").
                    Return([]response.RajaOngkirProvince{}, nil)
            },
            expectedResult: []response.ProvinceResponse{},
            expectedError:  "",
        },
        {
            name:    "repository_error",
            request: request.GetProvincesRequest{ID: "999"},
            mockSetup: func() {
                mockRepo.EXPECT().
                    GetProvinces("999").
                    Return(nil, errors.New("repository connection failed"))
            },
            expectedResult: nil,
            expectedError:  "failed to retrieve province data",
        },
        {
            name:    "repository_api_error",
            request: request.GetProvincesRequest{ID: "999"},
            mockSetup: func() {
                mockRepo.EXPECT().
                    GetProvinces("999").
                    Return(nil, &RepositoryError{
                        Operation: "GetProvinces",
                        Cause:     errors.New("API error: Invalid province ID"),
                    })
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

### Integration Tests Update
```go
// internal/integration/shipping_integration_test.go
package integration

import (
    "context"
    "testing"
    "time"

    "github.com/hanifbg/landing_backend/config"
    "github.com/hanifbg/landing_backend/internal/model/request"
    repoUtil "github.com/hanifbg/landing_backend/internal/repository/util"
    serviceUtil "github.com/hanifbg/landing_backend/internal/service/util"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestShippingIntegration_FullStack(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration tests in short mode")
    }

    cfg := loadIntegrationTestConfig()
    
    // Initialize full dependency chain
    repoWrapper, err := repoUtil.New(cfg)
    require.NoError(t, err)
    
    serviceWrapper, err := serviceUtil.New(cfg, repoWrapper)
    require.NoError(t, err)

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    t.Run("complete_shipping_flow", func(t *testing.T) {
        // Test provinces
        provinces, err := serviceWrapper.ShippingService.GetProvinces(request.GetProvincesRequest{})
        assert.NoError(t, err)
        assert.NotEmpty(t, provinces)

        // Get first province for cities test
        firstProvince := provinces[0]
        
        // Test cities
        cities, err := serviceWrapper.ShippingService.GetCities(request.GetCitiesRequest{
            ProvinceID: firstProvince.ProvinceID,
        })
        assert.NoError(t, err)
        assert.NotEmpty(t, cities)

        // Get first city for districts test
        firstCity := cities[0]
        
        // Test districts
        districts, err := serviceWrapper.ShippingService.GetDistricts(request.GetDistrictsRequest{
            CityID: firstCity.CityID,
        })
        assert.NoError(t, err)
        assert.NotEmpty(t, districts)

        // Test shipping cost calculation
        if len(districts) > 0 {
            costs, err := serviceWrapper.ShippingService.CalculateShippingCost(request.CalculateShippingRequest{
                Origin:      "501", // Jakarta Pusat
                Destination: districts[0].DistrictID,
                Weight:      1000,
                Courier:     "jne",
            })
            assert.NoError(t, err)
            assert.NotEmpty(t, costs)
        }
    })
}
```

## Implementation Steps

1. **Generate Repository Mocks**
   - Run `go generate` for repository interfaces
   - Verify mock files are created in correct locations
   - Update import paths in existing test files

2. **Update Service Layer Tests**
   - Replace `RajaOngkirClientInterface` mocks with `ShippingRepository` mocks
   - Update all test setups to use new mock interface
   - Verify all existing test scenarios are preserved

3. **Create Repository Unit Tests**
   - Create comprehensive test suite for repository layer
   - Test all HTTP client scenarios
   - Test error handling and edge cases
   - Test timeout and retry behavior

4. **Update Integration Tests**
   - Update integration tests to use new dependency chain
   - Test complete application flow
   - Verify external API integration still works

5. **Add Performance Tests**
   - Create benchmark tests for repository layer
   - Test concurrent access patterns
   - Measure performance impact of refactoring

6. **Update Test Utilities**
   - Create helper functions for common test scenarios
   - Update test configuration management
   - Add test data factories

## Dependencies
- **Requires**: SRP-004 (Dependency injection must be complete)
- **Files to modify**:
  - All test files in `internal/service/shipping/`
  - All integration test files
  - Mock generation directives
  - Test utilities and helpers

## Estimated Effort
- **Repository test creation**: 8 hours
- **Service test updates**: 6 hours
- **Integration test updates**: 4 hours
- **Mock generation and fixes**: 2 hours
- **Performance testing**: 3 hours
- **Total**: 23 hours

## Test Coverage Requirements

### Coverage Targets
- **Repository Layer**: 95%+ coverage
- **Service Layer**: 90%+ coverage  
- **Integration Tests**: Cover all critical paths
- **Overall**: Maintain >85% total coverage

### Coverage Verification
```bash
# Run coverage analysis
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# Generate coverage report
go tool cover -html=coverage.out -o coverage.html

# Check coverage by package
go tool cover -func=coverage.out | grep -E "(repository|service)"

# Verify minimum coverage threshold
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
if (( $(echo "$COVERAGE < 85" | bc -l) )); then
    echo "Error: Coverage $COVERAGE% below 85% threshold"
    exit 1
fi
```

### Test Quality Metrics
```go
// Add test quality validation
func TestCoverage_Quality(t *testing.T) {
    // Verify all public methods are tested
    // Verify error paths are tested
    // Verify edge cases are covered
}
```

## Performance Testing

### Repository Benchmarks
```go
// internal/repository/rajaongkir/benchmark_test.go
func BenchmarkRepository_GetProvinces(b *testing.B) {
    repo := setupBenchmarkRepository()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := repo.GetProvinces("")
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkRepository_GetProvincesParallel(b *testing.B) {
    repo := setupBenchmarkRepository()
    
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            _, err := repo.GetProvinces("")
            if err != nil {
                b.Fatal(err)
            }
        }
    })
}
```

### Load Testing
```go
// test/load/shipping_load_test.go
func TestShippingLoad_ConcurrentRequests(t *testing.T) {
    const (
        numGoroutines = 50
        requestsPerGoroutine = 20
    )

    // Test concurrent repository access
    // Measure response times and error rates
    // Verify no race conditions
}
```

## Mock Generation Management

### Mock Update Script
```bash
#!/bin/bash
# scripts/update-mocks.sh

echo "Updating all mocks..."

# Generate repository mocks
go generate ./internal/repository/...

# Generate service mocks (if needed)
go generate ./internal/service/...

# Verify mocks compile
go build ./internal/service/shipping/mocks/...

echo "Mock generation complete"
```

### Mock Verification
```go
// test/mock_verification_test.go
func TestMocks_CompileAndImplementInterfaces(t *testing.T) {
    // Verify all mocks implement their interfaces correctly
    var _ repository.ShippingRepository = &mocks.MockShippingRepository{}
    var _ service.ShippingService = &mocks.MockShippingService{}
}
```

## Definition of Done
- [x] All repository methods have comprehensive unit tests
- [x] All service layer tests use repository mocks instead of client mocks
- [x] Integration tests pass with new architecture
- [x] Test coverage is above 85% overall
- [x] Repository layer coverage is above 95%
- [x] All mocks are generated and up-to-date
- [x] Performance benchmarks show no significant regression
- [x] All existing test scenarios are preserved
- [x] Security considerations are properly tested
- [x] Test utilities are updated and documented
- [x] CI/CD pipeline passes all tests

## Risk Assessment

### High Risk
- **Test Coverage Drop**: Refactoring might reduce test coverage below 85%
  - **Mitigation**: Incremental testing updates, continuous coverage monitoring

### Medium Risk
- **Mock Generation Issues**: Generated mocks might have incorrect interfaces
  - **Mitigation**: Automated mock verification, compilation tests

- **Integration Test Failures**: External API dependencies might cause flaky tests
  - **Mitigation**: Proper test environment setup, retry mechanisms

### Low Risk
- **Performance Test Failures**: Benchmarks might show temporary regression
  - **Mitigation**: Baseline establishment, performance monitoring

## Notes
- Maintain test quality while updating architecture
- Focus on comprehensive error scenario testing
- Ensure external API integration tests are stable
- Document any changes to test patterns or utilities

## Security Considerations

### Secure Testing
- [x] Ensure tests don't expose API keys or credentials
- [x] Test security-related edge cases (empty keys, invalid credentials)
- [x] Verify error handling doesn't expose sensitive information

### Test Data Security
- [x] Use fake/mock API keys in tests
- [x] Sanitize test output to prevent credential leakage
- [x] Ensure test environment configurations are secure

### Mock Security
- [x] Ensure mocks properly simulate security boundaries
- [x] Test authorization and authentication failure scenarios
- [x] Verify proper input validation at boundary layers

## Related Tickets
- **Depends on**: SRP-004 (Update Dependency Injection)
- **Blocks**: SRP-006 (Legacy Cleanup)
- **Related**: All previous SRP tickets

## Labels
- `refactor`
- `testing`
- `mocks`
- `test-coverage`
