package shipping

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hanifbg/landing_backend/internal/model/response"
	"github.com/stretchr/testify/assert"
)

// Test NewRajaOngkirClient constructor
func TestNewRajaOngkirClient(t *testing.T) {
	t.Run("Success - Create RajaOngkir client", func(t *testing.T) {
		apiKey := "test-api-key"
		baseURL := "https://api.rajaongkir.com/starter/v2"

		client := NewRajaOngkirClient(apiKey, baseURL)

		assert.NotNil(t, client)
		assert.Equal(t, apiKey, client.apiKey)
		assert.Equal(t, baseURL, client.baseURL)
		assert.NotNil(t, client.client)
		assert.Equal(t, 30*time.Second, client.client.Timeout)
	})
}

// Test GetProvinces method
func TestRajaOngkirClient_GetProvinces(t *testing.T) {
	t.Run("Success - Get all provinces", func(t *testing.T) {
		// Create mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/destination/province", r.URL.Path)
			assert.Equal(t, "test-api-key", r.Header.Get("key"))

			komerceResp := response.KomerceProvinceResponse{
				Meta: response.Meta{
					Code:    200,
					Message: "OK",
				},
				Data: []response.RajaOngkirProvince{
					{
						ProvinceID: 1,
						Province:   "Bali",
					},
					{
						ProvinceID: 2,
						Province:   "Bangka Belitung",
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(komerceResp)
		}))
		defer server.Close()

		client := NewRajaOngkirClient("test-api-key", server.URL)
		result, err := client.GetProvinces("")

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, 1, result[0].ProvinceID)
		assert.Equal(t, "Bali", result[0].Province)
	})

	t.Run("Success - Get specific province by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/destination/province/1", r.URL.Path)

			// Komerce API returns a single object in the data array
			komerceResp := response.KomerceProvinceResponse{
				Meta: response.Meta{
					Code:    200,
					Message: "OK",
				},
				Data: []response.RajaOngkirProvince{
					{
						ProvinceID: 1,
						Province:   "Bali",
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(komerceResp)
		}))
		defer server.Close()

		client := NewRajaOngkirClient("test-api-key", server.URL)
		result, err := client.GetProvinces("1")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, 1, result[0].ProvinceID)
		assert.Equal(t, "Bali", result[0].Province)
	})

	t.Run("Error - API returns error status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/destination/province", r.URL.Path)
			komerceResp := response.KomerceProvinceResponse{
				Meta: response.Meta{
					Code:    400,
					Message: "Invalid API key",
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(komerceResp)
		}))
		defer server.Close()

		client := NewRajaOngkirClient("invalid-key", server.URL)
		result, err := client.GetProvinces("")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "API error: Invalid API key")
	})

	t.Run("Error - Invalid JSON response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("invalid json"))
		}))
		defer server.Close()

		client := NewRajaOngkirClient("test-api-key", server.URL)
		result, err := client.GetProvinces("")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to unmarshal response")
	})

	t.Run("Error - Server unavailable", func(t *testing.T) {
		client := NewRajaOngkirClient("test-api-key", "http://invalid-url")
		result, err := client.GetProvinces("")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to make request")
	})
}

// Test GetCities method
func TestRajaOngkirClient_GetCities(t *testing.T) {
		t.Run("Success - Get cities by province ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/destination/city/1", r.URL.Path)
			assert.Equal(t, "test-api-key", r.Header.Get("key"))

			komerceResp := response.KomerceCityResponse{
				Meta: response.Meta{
					Code:    200,
					Message: "OK",
				},
				Data: []response.RajaOngkirCity{
					{
						CityID:     1,
						// ProvinceID is populated from the request
						Province:   "Bali",
						Type:       "Kabupaten",
						CityName:   "Badung",
						PostalCode: "80351",
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(komerceResp)
		}))
		defer server.Close()

		client := NewRajaOngkirClient("test-api-key", server.URL)
		result, err := client.GetCities("1", "")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, 1, result[0].CityID)
		assert.Equal(t, 1, result[0].ProvinceID)
	})



		t.Run("Error - API returns error status for cities", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/destination/city/999", r.URL.Path)
			komerceResp := response.KomerceCityResponse{
				Meta: response.Meta{
					Code:    400,
					Message: "Invalid province ID",
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(komerceResp)
		}))
		defer server.Close()

		client := NewRajaOngkirClient("test-api-key", server.URL)
		result, err := client.GetCities("999", "")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "API error: Invalid province ID")
	})
}

// Test CalculateShippingCost method
func TestRajaOngkirClient_CalculateShippingCost(t *testing.T) {
	t.Run("Success - Calculate shipping cost", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/calculate/district/domestic-cost", r.URL.Path)
			assert.Equal(t, "test-api-key", r.Header.Get("key"))
			assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("content-type"))

			// Parse form data
			r.ParseForm()
			assert.Equal(t, "501", r.FormValue("origin"))
			assert.Equal(t, "114", r.FormValue("destination"))
			assert.Equal(t, "1000", r.FormValue("weight"))
			assert.Equal(t, "jne", r.FormValue("courier"))

			// Mock the new API response format
			response := struct {
				Meta struct {
					Message string `json:"message"`
					Code    int    `json:"code"`
					Status  string `json:"status"`
				} `json:"meta"`
				Data []struct {
					Name        string `json:"name"`
					Code        string `json:"code"`
					Service     string `json:"service"`
					Description string `json:"description"`
					Cost        int    `json:"cost"`
					ETD         string `json:"etd"`
				} `json:"data"`
			}{
				Meta: struct {
					Message string `json:"message"`
					Code    int    `json:"code"`
					Status  string `json:"status"`
				}{
					Message: "Success Calculate Domestic Shipping cost",
					Code:    200,
					Status:  "success",
				},
				Data: []struct {
					Name        string `json:"name"`
					Code        string `json:"code"`
					Service     string `json:"service"`
					Description string `json:"description"`
					Cost        int    `json:"cost"`
					ETD         string `json:"etd"`
				}{
					{
						Name:        "Jalur Nugraha Ekakurir (JNE)",
						Code:        "jne",
						Service:     "REG",
						Description: "Layanan Reguler",
						Cost:        15000,
						ETD:         "1-2 day",
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewRajaOngkirClient("test-api-key", server.URL)
		result, err := client.CalculateShippingCost("501", "114", 1000, "jne")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "jne", result[0].Code)
		assert.Equal(t, "Jalur Nugraha Ekakurir (JNE)", result[0].Name)
		assert.Len(t, result[0].Costs, 1)
		assert.Equal(t, "REG", result[0].Costs[0].Service)
	})

	t.Run("Error - API returns error status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Mock the new API error response format
			response := struct {
				Meta struct {
					Message string `json:"message"`
					Code    int    `json:"code"`
					Status  string `json:"status"`
				} `json:"meta"`
				Data interface{} `json:"data"`
			}{
				Meta: struct {
					Message string `json:"message"`
					Code    int    `json:"code"`
					Status  string `json:"status"`
				}{
					Message: "Invalid courier",
					Code:    400,
					Status:  "error",
				},
				Data: nil,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewRajaOngkirClient("test-api-key", server.URL)
		result, err := client.CalculateShippingCost("501", "114", 1000, "invalid")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "RajaOngkir API error: Invalid courier")
	})

	t.Run("Error - Server unavailable", func(t *testing.T) {
		client := NewRajaOngkirClient("test-api-key", "http://invalid-url")
		result, err := client.CalculateShippingCost("501", "114", 1000, "jne")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to make request")
	})
}