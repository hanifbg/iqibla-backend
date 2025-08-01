package shipping

import (
	"bytes"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/hanifbg/landing_backend/internal/model/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRajaOngkirClient is a mock implementation of RajaOngkirClientInterface for testing
type MockRajaOngkirClient struct {
	mock.Mock
}

func (m *MockRajaOngkirClient) GetProvinces(provinceID string) ([]response.RajaOngkirProvince, error) {
	args := m.Called(provinceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]response.RajaOngkirProvince), args.Error(1)
}

func (m *MockRajaOngkirClient) GetCities(provinceID, cityID string) ([]response.RajaOngkirCity, error) {
	args := m.Called(provinceID, cityID)
	return args.Get(0).([]response.RajaOngkirCity), args.Error(1)
}

func (m *MockRajaOngkirClient) GetDistricts(cityID string) ([]response.RajaOngkirDistrict, error) {
	args := m.Called(cityID)
	return args.Get(0).([]response.RajaOngkirDistrict), args.Error(1)
}

func (m *MockRajaOngkirClient) CalculateShippingCost(origin, destination string, weight int, courier string) ([]response.RajaOngkirCost, error) {
	args := m.Called(origin, destination, weight, courier)
	return args.Get(0).([]response.RajaOngkirCost), args.Error(1)
}

func TestCachedRajaOngkirClient_GetProvinces(t *testing.T) {
	// Setup
	mockClient := new(MockRajaOngkirClient)
	var logBuffer bytes.Buffer
	logger := log.New(&logBuffer, "", 0)
	cachedClient := NewCachedRajaOngkirClient(mockClient, 1*time.Hour, logger)

	// Test data
	testProvinces := []response.RajaOngkirProvince{
		{ProvinceID: 1, Province: "Test Province 1"},
		{ProvinceID: 2, Province: "Test Province 2"},
	}

	// Test scenario 1: Cache miss, successful API call
	mockClient.On("GetProvinces", "").Return(testProvinces, nil).Once()

	provinces, err := cachedClient.GetProvinces("")
	assert.NoError(t, err, "GetProvinces should not return an error")
	assert.Equal(t, testProvinces, provinces, "Should return provinces from API")
	assert.Contains(t, logBuffer.String(), "Cache miss", "Should log cache miss")

	// Clear log buffer
	logBuffer.Reset()

	// Test scenario 2: Cache hit
	provinces, err = cachedClient.GetProvinces("")
	assert.NoError(t, err, "GetProvinces should not return an error on cache hit")
	assert.Equal(t, testProvinces, provinces, "Should return provinces from cache")
	assert.Contains(t, logBuffer.String(), "Cache hit", "Should log cache hit")

	// Clear log buffer
	logBuffer.Reset()

	// Test scenario 3: Specific province ID bypasses cache
	specificProvinces := []response.RajaOngkirProvince{{ProvinceID: 1, Province: "Specific Province"}}
	mockClient.On("GetProvinces", "1").Return(specificProvinces, nil).Once()

	provinces, err = cachedClient.GetProvinces("1")
	assert.NoError(t, err, "GetProvinces with specific ID should not return an error")
	assert.Equal(t, specificProvinces, provinces, "Should return specific province from API")
	assert.Contains(t, logBuffer.String(), "Cache bypassed", "Should log cache bypass")

	// Clear log buffer
	logBuffer.Reset()

	// Test scenario 4: API error
	mockClient.On("GetProvinces", "").Return(nil, errors.New("API error")).Once()

	// Clear cache
	cachedClient.cache.Clear()

	provinces, err = cachedClient.GetProvinces("")
	assert.Error(t, err, "GetProvinces should return an error when API fails")
	assert.Nil(t, provinces, "No provinces should be returned on API error")
	assert.Contains(t, logBuffer.String(), "Cache miss", "Should log cache miss")

	// Verify all expectations were met
	mockClient.AssertExpectations(t)
}

func TestCachedRajaOngkirClient_InitCache(t *testing.T) {
	// Setup
	mockClient := new(MockRajaOngkirClient)
	var logBuffer bytes.Buffer
	logger := log.New(&logBuffer, "", 0)
	cachedClient := NewCachedRajaOngkirClient(mockClient, 1*time.Hour, logger)

	// Test data
	testProvinces := []response.RajaOngkirProvince{
		{ProvinceID: 1, Province: "Test Province 1"},
		{ProvinceID: 2, Province: "Test Province 2"},
	}

	// Test scenario 1: Successful initialization
	mockClient.On("GetProvinces", "").Return(testProvinces, nil).Once()

	err := cachedClient.InitCache()
	assert.NoError(t, err, "InitCache should not return an error")
	assert.True(t, cachedClient.initialized, "Client should be marked as initialized")
	assert.Contains(t, logBuffer.String(), "Provinces cache warmed up successfully", "Should log successful warm-up")

	// Clear log buffer
	logBuffer.Reset()

	// Test scenario 2: Already initialized
	err = cachedClient.InitCache()
	assert.NoError(t, err, "InitCache should not return an error when already initialized")
	assert.Empty(t, logBuffer.String(), "Should not log anything when already initialized")

	// Clear log buffer
	logBuffer.Reset()

	// Test scenario 3: Initialization failure
	// Reset client
	cachedClient = NewCachedRajaOngkirClient(mockClient, 1*time.Hour, logger)
	mockClient.On("GetProvinces", "").Return(nil, errors.New("API error")).Once()

	err = cachedClient.InitCache()
	assert.Error(t, err, "InitCache should return an error when API fails")
	assert.False(t, cachedClient.initialized, "Client should not be marked as initialized on error")
	assert.Contains(t, logBuffer.String(), "Failed to warm up", "Should log initialization failure")

	// Verify all expectations were met
	mockClient.AssertExpectations(t)
}

func TestCachedRajaOngkirClient_DelegatedMethods(t *testing.T) {
	// Setup
	mockClient := new(MockRajaOngkirClient)
	var logBuffer bytes.Buffer
	logger := log.New(&logBuffer, "", 0)
	cachedClient := NewCachedRajaOngkirClient(mockClient, 1*time.Hour, logger)

	// Test data
	testCities := []response.RajaOngkirCity{
		{CityID: 1, CityName: "Test City 1"},
	}
	testDistricts := []response.RajaOngkirDistrict{
		{DistrictID: 1, DistrictName: "Test District 1"},
	}
	testCosts := []response.RajaOngkirCost{
		{Code: "jne", Name: "JNE"},
	}

	// Test GetCities
	mockClient.On("GetCities", "1", "").Return(testCities, nil).Once()
	cities, err := cachedClient.GetCities("1", "")
	assert.NoError(t, err, "GetCities should not return an error")
	assert.Equal(t, testCities, cities, "Should delegate GetCities to underlying client")

	// Test GetDistricts
	mockClient.On("GetDistricts", "1").Return(testDistricts, nil).Once()
	districts, err := cachedClient.GetDistricts("1")
	assert.NoError(t, err, "GetDistricts should not return an error")
	assert.Equal(t, testDistricts, districts, "Should delegate GetDistricts to underlying client")

	// Test CalculateShippingCost
	mockClient.On("CalculateShippingCost", "1", "2", 1000, "jne").Return(testCosts, nil).Once()
	costs, err := cachedClient.CalculateShippingCost("1", "2", 1000, "jne")
	assert.NoError(t, err, "CalculateShippingCost should not return an error")
	assert.Equal(t, testCosts, costs, "Should delegate CalculateShippingCost to underlying client")

	// Verify all expectations were met
	mockClient.AssertExpectations(t)
}
