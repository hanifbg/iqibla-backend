package rajaongkir

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/hanifbg/landing_backend/internal/model/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_GetProvinces(t *testing.T) {
	tests := []struct {
		name           string
		provinceID     string
		mockResponse   string
		mockStatusCode int
		expectedResult []response.RajaOngkirProvince
		expectError    bool
	}{
		{
			name:       "success_all_provinces",
			provinceID: "",
			mockResponse: `{
				"meta": {"code": 200, "message": "OK", "status": "success"},
				"data": [
					{"id": 1, "name": "Bali"},
					{"id": 2, "name": "Jawa Barat"}
				]
			}`,
			mockStatusCode: http.StatusOK,
			expectedResult: []response.RajaOngkirProvince{
				{ProvinceID: 1, Province: "Bali"},
				{ProvinceID: 2, Province: "Jawa Barat"},
			},
			expectError: false,
		},
		{
			name:       "success_specific_province",
			provinceID: "1",
			mockResponse: `{
				"meta": {"code": 200, "message": "OK", "status": "success"},
				"data": [{"id": 1, "name": "Bali"}]
			}`,
			mockStatusCode: http.StatusOK,
			expectedResult: []response.RajaOngkirProvince{
				{ProvinceID: 1, Province: "Bali"},
			},
			expectError: false,
		},
		{
			name:           "error_api_failure",
			provinceID:     "",
			mockResponse:   `{"meta": {"code": 400, "message": "Bad Request", "status": "error"}}`,
			mockStatusCode: http.StatusOK,
			expectedResult: nil,
			expectError:    true,
		},
		{
			name:           "error_invalid_json",
			provinceID:     "",
			mockResponse:   `{invalid json}`,
			mockStatusCode: http.StatusOK,
			expectedResult: nil,
			expectError:    true,
		},
		{
			name:           "error_http_error",
			provinceID:     "",
			mockResponse:   ``,
			mockStatusCode: http.StatusInternalServerError,
			expectedResult: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check request method
				assert.Equal(t, "GET", r.Method)

				// Check API key in header
				assert.Equal(t, "test-api-key", r.Header.Get("key"))

				// Check URL path - should include province ID if provided
				expectedPath := "/destination/province"
				if tt.provinceID != "" {
					expectedPath += "/" + tt.provinceID
				}
				assert.Equal(t, expectedPath, r.URL.Path)

				// Return mock response
				w.WriteHeader(tt.mockStatusCode)
				_, err := w.Write([]byte(tt.mockResponse))
				assert.NoError(t, err)
			}))
			defer server.Close()

			// Create repository with test server URL and a non-nil HTTP client
			repo := NewRepository(Config{
				APIKey:  "test-api-key",
				BaseURL: server.URL,
				Client:  server.Client(),
			})

			// Execute test
			provinces, err := repo.GetProvinces(tt.provinceID)

			// Check results
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, provinces)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, provinces)
			}
		})
	}
}

func TestRepository_GetCities(t *testing.T) {
	tests := []struct {
		name           string
		provinceID     string
		cityID         string
		mockResponse   string
		mockStatusCode int
		expectedResult []response.RajaOngkirCity
		expectError    bool
	}{
		{
			name:       "success_all_cities",
			provinceID: "",
			cityID:     "",
			mockResponse: `{
				"meta": {"code": 200, "message": "OK", "status": "success"},
				"data": [
					{"id": 1, "province_id": 1, "province_name": "Bali", "type": "Kabupaten", "name": "Badung", "postal_code": "80351"},
					{"id": 2, "province_id": 1, "province_name": "Bali", "type": "Kabupaten", "name": "Bangli", "postal_code": "80619"}
				]
			}`,
			mockStatusCode: http.StatusOK,
			expectedResult: []response.RajaOngkirCity{
				{CityID: 1, ProvinceID: 1, Province: "Bali", Type: "Kabupaten", CityName: "Badung", PostalCode: "80351"},
				{CityID: 2, ProvinceID: 1, Province: "Bali", Type: "Kabupaten", CityName: "Bangli", PostalCode: "80619"},
			},
			expectError: false,
		},
		{
			name:       "success_province_filter",
			provinceID: "1",
			cityID:     "",
			mockResponse: `{
				"meta": {"code": 200, "message": "OK", "status": "success"},
				"data": [
					{"id": 1, "province_id": 1, "province_name": "Bali", "type": "Kabupaten", "name": "Badung", "postal_code": "80351"},
					{"id": 2, "province_id": 1, "province_name": "Bali", "type": "Kabupaten", "name": "Bangli", "postal_code": "80619"}
				]
			}`,
			mockStatusCode: http.StatusOK,
			expectedResult: []response.RajaOngkirCity{
				{CityID: 1, ProvinceID: 1, Province: "Bali", Type: "Kabupaten", CityName: "Badung", PostalCode: "80351"},
				{CityID: 2, ProvinceID: 1, Province: "Bali", Type: "Kabupaten", CityName: "Bangli", PostalCode: "80619"},
			},
			expectError: false,
		},
		{
			name:       "success_specific_city",
			provinceID: "1",
			cityID:     "1",
			mockResponse: `{
				"meta": {"code": 200, "message": "OK", "status": "success"},
				"data": [
					{"id": 1, "province_id": 1, "province_name": "Bali", "type": "Kabupaten", "name": "Badung", "postal_code": "80351"}
				]
			}`,
			mockStatusCode: http.StatusOK,
			expectedResult: []response.RajaOngkirCity{
				{CityID: 1, ProvinceID: 1, Province: "Bali", Type: "Kabupaten", CityName: "Badung", PostalCode: "80351"},
			},
			expectError: false,
		},
		{
			name:           "error_api_failure",
			provinceID:     "1",
			cityID:         "",
			mockResponse:   `{"meta": {"code": 400, "message": "Bad Request", "status": "error"}}`,
			mockStatusCode: http.StatusOK,
			expectedResult: nil,
			expectError:    true,
		},
		{
			name:           "error_http_error",
			provinceID:     "",
			cityID:         "",
			mockResponse:   ``,
			mockStatusCode: http.StatusInternalServerError,
			expectedResult: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check request method and headers
				assert.Equal(t, "GET", r.Method)
				assert.Equal(t, "test-api-key", r.Header.Get("key"))

				// Check URL path - for province ID
				expectedPath := "/destination/city"
				if tt.provinceID != "" {
					expectedPath += "/" + tt.provinceID
				}
				assert.Equal(t, expectedPath, r.URL.Path)

				// Check query parameters for cityID
				if tt.cityID != "" {
					assert.Equal(t, tt.cityID, r.URL.Query().Get("id"))
				}

				// Modify the mock response to use the simplified format if it contains city data
				// This is needed because the API has changed its response format
				if !tt.expectError && strings.Contains(tt.mockResponse, `"province_id"`) {
					// Convert the old format to new simplified format
					var oldResp response.KomerceCityResponse
					err := json.Unmarshal([]byte(tt.mockResponse), &oldResp)
					require.NoError(t, err)

					// Create simplified response with only id and name
					for i := range oldResp.Data {
						// Clear fields that are no longer in the API response
						oldResp.Data[i].Province = ""
						oldResp.Data[i].Type = ""
						oldResp.Data[i].PostalCode = ""
						// ProvinceID is set based on URL in the actual implementation
						if tt.provinceID != "" {
							provinceID, err := strconv.Atoi(tt.provinceID)
							if err == nil {
								oldResp.Data[i].ProvinceID = provinceID
							} else {
								oldResp.Data[i].ProvinceID = 0
							}
						} else {
							oldResp.Data[i].ProvinceID = 0
						}
					}

					// Marshal to JSON
					newResponse, err := json.Marshal(oldResp)
					require.NoError(t, err)

					// Return the simplified response
					w.WriteHeader(tt.mockStatusCode)
					_, err = w.Write(newResponse)
					require.NoError(t, err)
					return
				}

				// Return original mock response
				w.WriteHeader(tt.mockStatusCode)
				_, err := w.Write([]byte(tt.mockResponse))
				require.NoError(t, err)
			}))
			defer server.Close()

			// Create repository with test server URL and a non-nil HTTP client
			repo := NewRepository(Config{
				APIKey:  "test-api-key",
				BaseURL: server.URL,
				Client:  server.Client(),
			})

			// Execute test
			cities, err := repo.GetCities(tt.provinceID, tt.cityID)

			// Check results
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, cities)
			} else {
				// For successful cases, we only check that we get cities with the correct IDs and names
				// We don't check other fields because they're no longer provided by the API
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedResult), len(cities), "Should have the same number of cities")

				for i, expectedCity := range tt.expectedResult {
					assert.Equal(t, expectedCity.CityID, cities[i].CityID, "CityID should match")
					assert.Equal(t, expectedCity.CityName, cities[i].CityName, "CityName should match")

					// Check provinceID if we specifically passed one
					if tt.provinceID != "" {
						provinceID, err := strconv.Atoi(tt.provinceID)
						if err == nil {
							assert.Equal(t, provinceID, cities[i].ProvinceID, "ProvinceID should match the path parameter")
						}
					}
				}
			}
		})
	}
}

// Additional tests for GetDistricts and CalculateShippingCost methods would follow the same pattern
