# Testing Strategy: Shipping Repository Refactor

## Testing Requirements

### Coverage Goals
- **Minimum Coverage**: 85% across all modified components
- **Target Coverage**: 90%+ for new repository layer
- **Critical Path Coverage**: 100% for external API interactions

**Note**: Since the service is not yet released, we can focus on comprehensive testing of the new architecture without backward compatibility test scenarios.

### Testing Pyramid Strategy

```
    /\
   /  \    E2E Tests (10%)
  /____\   Integration Tests (20%)
 /______\  Unit Tests (70%)
```

## Test Categories

### 1. Unit Tests (70% of test effort)

#### Repository Layer Unit Tests
**Location**: `internal/repository/rajaongkir/client_test.go`

**Coverage Areas**:
- HTTP client configuration and initialization
- Request construction and header setting
- Response parsing and error handling
- Timeout and retry mechanisms
- Input validation and sanitization

**Mock Strategy**:
```go
func TestRepository_GetProvinces(t *testing.T) {
    tests := []struct {
        name           string
        provinceID     string
        mockResponse   string
        mockStatusCode int
        expectedResult []response.RajaOngkirProvince
        expectedError  string
    }{
        {
            name:           "successful_request_with_valid_response",
            provinceID:     "1",
            mockResponse:   `{"meta":{"code":200},"data":[{"id":1,"province":"Bali"}]}`,
            mockStatusCode: 200,
            expectedResult: []response.RajaOngkirProvince{{ProvinceID: 1, Province: "Bali"}},
            expectedError:  "",
        },
        {
            name:           "api_error_response",
            provinceID:     "999",
            mockResponse:   `{"meta":{"code":400,"message":"Invalid province ID"}}`,
            mockStatusCode: 200,
            expectedResult: nil,
            expectedError:  "API error: Invalid province ID",
        },
        {
            name:           "network_timeout",
            provinceID:     "1",
            mockResponse:   "",
            mockStatusCode: 0, // Indicates timeout
            expectedResult: nil,
            expectedError:  "failed to make request",
        },
        {
            name:           "invalid_json_response",
            provinceID:     "1",
            mockResponse:   `{invalid json}`,
            mockStatusCode: 200,
            expectedResult: nil,
            expectedError:  "failed to unmarshal response",
        },
        {
            name:           "empty_province_id",
            provinceID:     "",
            mockResponse:   "",
            mockStatusCode: 0,
            expectedResult: []response.RajaOngkirProvince{},
            expectedError:  "",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Use httptest.Server for realistic HTTP testing
            server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                // Verify request headers and parameters
                assert.Equal(t, "test-api-key", r.Header.Get("key"))
                
                if tt.mockStatusCode == 0 {
                    // Simulate timeout
                    time.Sleep(100 * time.Millisecond)
                    return
                }
                
                w.WriteHeader(tt.mockStatusCode)
                w.Write([]byte(tt.mockResponse))
            }))
            defer server.Close()

            repo := NewRepository(Config{
                APIKey:  "test-api-key",
                BaseURL: server.URL,
                Timeout: 50 * time.Millisecond, // Short timeout for timeout tests
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
```

#### Service Layer Unit Tests
**Location**: `internal/service/shipping/impl_test.go`

**Coverage Areas**:
- Business logic and data transformation
- Error handling and propagation
- Request validation
- Response formatting

**Mock Strategy**:
```go
func TestShippingService_GetProvinces(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockShippingRepository(ctrl)
    service := &ShippingService{shippingRepo: mockRepo}

    tests := []struct {
        name           string
        request        request.GetProvincesRequest
        mockReturn     []response.RajaOngkirProvince
        mockError      error
        expectedResult []response.ProvinceResponse
        expectedError  string
    }{
        {
            name:    "successful_province_transformation",
            request: request.GetProvincesRequest{ID: "1"},
            mockReturn: []response.RajaOngkirProvince{
                {ProvinceID: 1, Province: "Bali"},
                {ProvinceID: 2, Province: "Jakarta"},
            },
            mockError: nil,
            expectedResult: []response.ProvinceResponse{
                {ProvinceID: "1", Province: "Bali"},
                {ProvinceID: "2", Province: "Jakarta"},
            },
            expectedError: "",
        },
        {
            name:           "repository_error_propagation",
            request:        request.GetProvincesRequest{ID: "999"},
            mockReturn:     nil,
            mockError:      errors.New("repository error"),
            expectedResult: nil,
            expectedError:  "repository error",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo.EXPECT().
                GetProvinces(tt.request.ID).
                Return(tt.mockReturn, tt.mockError)

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

### 2. Integration Tests (20% of test effort)

#### API Integration Tests
**Location**: `internal/integration/shipping_test.go`

**Coverage Areas**:
- End-to-end API workflows
- Real external API interactions (in test environment)
- Configuration validation
- Error scenarios with real timeouts

```go
func TestShippingIntegration_RealAPI(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration tests in short mode")
    }

    // Use test configuration with real API credentials
    cfg := &config.AppConfig{
        RajaOngkirAPIKey: os.Getenv("RAJAONGKIR_TEST_API_KEY"),
        RajaOngkirBaseURL: "https://api.rajaongkir.com/starter",
    }

    if cfg.RajaOngkirAPIKey == "" {
        t.Skip("Test API key not provided")
    }

    // Initialize full stack
    repoWrapper, err := repoUtil.New(cfg)
    require.NoError(t, err)

    serviceWrapper, err := serviceUtil.New(cfg, repoWrapper)
    require.NoError(t, err)

    tests := []struct {
        name          string
        testFunc      func(t *testing.T, service service.ShippingService)
        timeout       time.Duration
    }{
        {
            name:    "get_all_provinces",
            timeout: 10 * time.Second,
            testFunc: func(t *testing.T, service service.ShippingService) {
                req := request.GetProvincesRequest{}
                
                ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
                defer cancel()

                result, err := service.GetProvinces(req)
                
                assert.NoError(t, err)
                assert.NotEmpty(t, result)
                assert.Greater(t, len(result), 10) // Indonesia has 34 provinces
                
                // Verify expected provinces exist
                provinceNames := make([]string, len(result))
                for i, p := range result {
                    provinceNames[i] = p.Province
                }
                assert.Contains(t, provinceNames, "Bali")
                assert.Contains(t, provinceNames, "DKI Jakarta")
            },
        },
        {
            name:    "get_cities_for_jakarta",
            timeout: 10 * time.Second,
            testFunc: func(t *testing.T, service service.ShippingService) {
                // First get Jakarta province ID
                provincesReq := request.GetProvincesRequest{}
                provinces, err := service.GetProvinces(provincesReq)
                require.NoError(t, err)

                var jakartaID string
                for _, p := range provinces {
                    if strings.Contains(p.Province, "Jakarta") {
                        jakartaID = p.ProvinceID
                        break
                    }
                }
                require.NotEmpty(t, jakartaID, "Jakarta province not found")

                // Get cities for Jakarta
                citiesReq := request.GetCitiesRequest{ProvinceID: jakartaID}
                cities, err := service.GetCities(citiesReq)
                
                assert.NoError(t, err)
                assert.NotEmpty(t, cities)
                
                // Verify all cities belong to Jakarta
                for _, city := range cities {
                    assert.Equal(t, jakartaID, city.ProvinceID)
                    assert.Contains(t, city.Province, "Jakarta")
                }
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            timeout := tt.timeout
            if timeout == 0 {
                timeout = 30 * time.Second
            }

            done := make(chan bool)
            go func() {
                tt.testFunc(t, serviceWrapper.ShippingService)
                done <- true
            }()

            select {
            case <-done:
                // Test completed successfully
            case <-time.After(timeout):
                t.Fatalf("Test timed out after %v", timeout)
            }
        })
    }
}
```

#### Handler Integration Tests
**Location**: `internal/handler/shipping/integration_test.go`

```go
func TestShippingHandler_Integration(t *testing.T) {
    // Set up test server with real dependencies
    e := echo.New()
    cfg := loadTestConfig()
    
    repoWrapper, err := repoUtil.New(cfg)
    require.NoError(t, err)

    serviceWrapper, err := serviceUtil.New(cfg, repoWrapper)
    require.NoError(t, err)

    // Initialize handlers
    shipping.InitRoute(e, serviceWrapper)

    tests := []struct {
        name           string
        method         string
        path           string
        body           interface{}
        expectedStatus int
        validateBody   func(t *testing.T, body []byte)
    }{
        {
            name:           "get_provinces_success",
            method:         "GET",
            path:           "/api/v1/shipping/provinces",
            expectedStatus: 200,
            validateBody: func(t *testing.T, body []byte) {
                var response map[string]interface{}
                err := json.Unmarshal(body, &response)
                assert.NoError(t, err)
                
                data, exists := response["data"]
                assert.True(t, exists)
                
                provinces, ok := data.([]interface{})
                assert.True(t, ok)
                assert.Greater(t, len(provinces), 10)
            },
        },
        {
            name:           "calculate_shipping_cost",
            method:         "POST",
            path:           "/api/v1/shipping/cost",
            body: map[string]interface{}{
                "origin":      "501", // Jakarta Pusat
                "destination": "114", // Badung, Bali
                "weight":      1000,
                "courier":     "jne",
            },
            expectedStatus: 200,
            validateBody: func(t *testing.T, body []byte) {
                var response map[string]interface{}
                err := json.Unmarshal(body, &response)
                assert.NoError(t, err)
                
                data, exists := response["data"]
                assert.True(t, exists)
                
                costs, ok := data.([]interface{})
                assert.True(t, ok)
                assert.NotEmpty(t, costs)
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var req *http.Request
            var err error

            if tt.body != nil {
                bodyBytes, _ := json.Marshal(tt.body)
                req = httptest.NewRequest(tt.method, tt.path, bytes.NewReader(bodyBytes))
                req.Header.Set("Content-Type", "application/json")
            } else {
                req = httptest.NewRequest(tt.method, tt.path, nil)
            }

            rec := httptest.NewRecorder()
            e.ServeHTTP(rec, req)

            assert.Equal(t, tt.expectedStatus, rec.Code)
            
            if tt.validateBody != nil {
                tt.validateBody(t, rec.Body.Bytes())
            }
        })
    }
}
```

### 3. End-to-End Tests (10% of test effort)

#### Full Stack E2E Tests
**Location**: `test/e2e/shipping_test.go`

```go
func TestShippingE2E_CompleteUserJourney(t *testing.T) {
    if os.Getenv("E2E_TEST") != "true" {
        t.Skip("E2E tests skipped (set E2E_TEST=true to run)")
    }

    // Start the application in test mode
    app := startTestApplication()
    defer app.Shutdown()

    client := &http.Client{Timeout: 30 * time.Second}
    baseURL := "http://localhost:8080"

    t.Run("complete_shipping_calculation_flow", func(t *testing.T) {
        // Step 1: Get provinces
        resp, err := client.Get(baseURL + "/api/v1/shipping/provinces")
        require.NoError(t, err)
        require.Equal(t, 200, resp.StatusCode)

        var provincesResp map[string]interface{}
        err = json.NewDecoder(resp.Body).Decode(&provincesResp)
        require.NoError(t, err)
        resp.Body.Close()

        provinces := provincesResp["data"].([]interface{})
        require.NotEmpty(t, provinces)

        // Step 2: Get cities for first province
        firstProvince := provinces[0].(map[string]interface{})
        provinceID := firstProvince["province_id"].(string)

        resp, err = client.Get(fmt.Sprintf("%s/api/v1/shipping/cities/%s", baseURL, provinceID))
        require.NoError(t, err)
        require.Equal(t, 200, resp.StatusCode)

        var citiesResp map[string]interface{}
        err = json.NewDecoder(resp.Body).Decode(&citiesResp)
        require.NoError(t, err)
        resp.Body.Close()

        cities := citiesResp["data"].([]interface{})
        require.NotEmpty(t, cities)

        // Step 3: Get districts for first city
        firstCity := cities[0].(map[string]interface{})
        cityID := firstCity["city_id"].(string)

        resp, err = client.Get(fmt.Sprintf("%s/api/v1/shipping/districts/%s", baseURL, cityID))
        require.NoError(t, err)
        require.Equal(t, 200, resp.StatusCode)

        var districtsResp map[string]interface{}
        err = json.NewDecoder(resp.Body).Decode(&districtsResp)
        require.NoError(t, err)
        resp.Body.Close()

        districts := districtsResp["data"].([]interface{})
        require.NotEmpty(t, districts)

        // Step 4: Calculate shipping cost
        firstDistrict := districts[0].(map[string]interface{})
        districtID := firstDistrict["district_id"].(string)

        costRequest := map[string]interface{}{
            "origin":      "501", // Jakarta Pusat
            "destination": districtID,
            "weight":      1000,
            "courier":     "jne",
        }

        requestBody, _ := json.Marshal(costRequest)
        resp, err = client.Post(
            baseURL+"/api/v1/shipping/cost",
            "application/json",
            bytes.NewReader(requestBody),
        )
        require.NoError(t, err)
        require.Equal(t, 200, resp.StatusCode)

        var costResp map[string]interface{}
        err = json.NewDecoder(resp.Body).Decode(&costResp)
        require.NoError(t, err)
        resp.Body.Close()

        costs := costResp["data"].([]interface{})
        require.NotEmpty(t, costs)

        // Verify cost structure
        firstCost := costs[0].(map[string]interface{})
        require.Contains(t, firstCost, "service")
        require.Contains(t, firstCost, "cost")
        require.Contains(t, firstCost, "etd")
    })
}
```

## Performance Testing

### Load Testing
**Location**: `test/performance/shipping_load_test.go`

```go
func TestShippingPerformance_Load(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping performance tests in short mode")
    }

    const (
        concurrentUsers = 50
        requestsPerUser = 20
        maxResponseTime = 2 * time.Second
    )

    app := startTestApplication()
    defer app.Shutdown()

    client := &http.Client{
        Timeout: 5 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns:       100,
            IdleConnTimeout:    90 * time.Second,
            DisableCompression: true,
        },
    }

    var wg sync.WaitGroup
    results := make(chan time.Duration, concurrentUsers*requestsPerUser)
    errors := make(chan error, concurrentUsers*requestsPerUser)

    // Load test provinces endpoint
    for i := 0; i < concurrentUsers; i++ {
        wg.Add(1)
        go func(userID int) {
            defer wg.Done()
            
            for j := 0; j < requestsPerUser; j++ {
                start := time.Now()
                
                resp, err := client.Get("http://localhost:8080/api/v1/shipping/provinces")
                if err != nil {
                    errors <- err
                    continue
                }
                
                resp.Body.Close()
                duration := time.Since(start)
                results <- duration
                
                if resp.StatusCode != 200 {
                    errors <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
                }
            }
        }(i)
    }

    wg.Wait()
    close(results)
    close(errors)

    // Analyze results
    var responseTimes []time.Duration
    for duration := range results {
        responseTimes = append(responseTimes, duration)
    }

    var errorList []error
    for err := range errors {
        errorList = append(errorList, err)
    }

    // Performance assertions
    assert.Empty(t, errorList, "No errors should occur during load test")
    assert.NotEmpty(t, responseTimes, "Should have response times recorded")

    // Calculate statistics
    sort.Slice(responseTimes, func(i, j int) bool {
        return responseTimes[i] < responseTimes[j]
    })

    p50 := responseTimes[len(responseTimes)/2]
    p95 := responseTimes[int(float64(len(responseTimes))*0.95)]
    p99 := responseTimes[int(float64(len(responseTimes))*0.99)]

    t.Logf("Performance Results:")
    t.Logf("  Total requests: %d", len(responseTimes))
    t.Logf("  Errors: %d", len(errorList))
    t.Logf("  P50: %v", p50)
    t.Logf("  P95: %v", p95)
    t.Logf("  P99: %v", p99)

    // Performance requirements
    assert.Less(t, p95, maxResponseTime, "95th percentile should be under %v", maxResponseTime)
    assert.Less(t, float64(len(errorList))/float64(len(responseTimes)), 0.01, "Error rate should be under 1%")
}
```

## Test Coverage Monitoring

### Coverage Collection Script
**Location**: `scripts/test-coverage.sh`

```bash
#!/bin/bash
set -e

echo "Running test coverage analysis..."

# Clean previous coverage files
rm -f coverage.out coverage.html

# Run tests with coverage
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Extract coverage percentage
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')

echo "Total coverage: ${COVERAGE}%"

# Check coverage threshold
THRESHOLD=85
if (( $(echo "$COVERAGE < $THRESHOLD" | bc -l) )); then
    echo "Error: Coverage ${COVERAGE}% is below threshold ${THRESHOLD}%"
    exit 1
fi

echo "Coverage check passed!"

# Generate coverage report by package
echo "Coverage by package:"
go tool cover -func=coverage.out | grep -E "(internal/repository|internal/service)" | sort
```

### CI/CD Integration
**Location**: `.github/workflows/test-coverage.yml`

```yaml
name: Test Coverage

on:
  pull_request:
    branches: [ main ]
  push:
    branches: [ main ]

jobs:
  coverage:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Run tests with coverage
      run: |
        go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
    
    - name: Check coverage threshold
      run: |
        COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
        echo "Coverage: ${COVERAGE}%"
        if (( $(echo "$COVERAGE < 85" | bc -l) )); then
          echo "Error: Coverage ${COVERAGE}% is below 85%"
          exit 1
        fi
    
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
        flags: unittests
        name: codecov-umbrella
```

## Test Execution Strategy

### Local Development
```bash
# Unit tests only
make test-unit

# Integration tests (requires test API key)
export RAJAONGKIR_TEST_API_KEY="your-test-key"
make test-integration

# All tests including E2E
export E2E_TEST=true
make test-all

# Coverage report
make test-coverage
```

### CI/CD Pipeline
1. **Unit Tests**: Run on every commit
2. **Integration Tests**: Run on PR creation and merge
3. **E2E Tests**: Run nightly and before releases
4. **Performance Tests**: Run weekly and before major releases

### Test Data Management
- Use consistent test data across all test types
- Mock external API responses for unit tests
- Use test environment for integration tests
- Implement test data cleanup procedures

## Success Criteria

### Coverage Metrics
- [ ] Overall test coverage ≥ 85%
- [ ] Repository layer coverage ≥ 90%
- [ ] Service layer coverage ≥ 95%
- [ ] Critical path coverage = 100%

### Performance Metrics
- [ ] Unit tests complete in < 30 seconds
- [ ] Integration tests complete in < 2 minutes
- [ ] E2E tests complete in < 5 minutes
- [ ] API response times remain within 10% of baseline

### Quality Metrics
- [ ] Zero test flakiness
- [ ] All edge cases covered
- [ ] Error scenarios properly tested
- [ ] Mock usage follows best practices
