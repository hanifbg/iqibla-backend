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
		baseURL := "https://api.rajaongkir.com/starter"

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
			assert.Equal(t, "/province", r.URL.Path)
			assert.Equal(t, "test-api-key", r.Header.Get("key"))

			response := response.RajaOngkirResponse{
				RajaOngkir: struct {
					Query   interface{} `json:"query"`
					Status  response.Status      `json:"status"`
					Results interface{} `json:"results"`
				}{
					Status: response.Status{
						Code:        200,
						Description: "OK",
					},
					Results: []response.RajaOngkirProvince{
						{
							ProvinceID: "1",
							Province:   "Bali",
						},
						{
							ProvinceID: "2",
							Province:   "Bangka Belitung",
						},
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewRajaOngkirClient("test-api-key", server.URL)
		result, err := client.GetProvinces("")

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "1", result[0].ProvinceID)
		assert.Equal(t, "Bali", result[0].Province)
	})

	t.Run("Success - Get specific province by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/province", r.URL.Path)
			assert.Equal(t, "id=1", r.URL.RawQuery)

			// RajaOngkir API returns a single object when querying by ID
			response := response.RajaOngkirResponse{
				RajaOngkir: struct {
					Query   interface{} `json:"query"`
					Status  response.Status      `json:"status"`
					Results interface{} `json:"results"`
				}{
					Status: response.Status{
						Code:        200,
						Description: "OK",
					},
					// Single object response for specific province
					Results: response.RajaOngkirProvince{
						ProvinceID: "1",
						Province:   "Bali",
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewRajaOngkirClient("test-api-key", server.URL)
		result, err := client.GetProvinces("1")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "1", result[0].ProvinceID)
		assert.Equal(t, "Bali", result[0].Province)
	})

	t.Run("Error - API returns error status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := response.RajaOngkirResponse{
				RajaOngkir: struct {
					Query   interface{} `json:"query"`
					Status  response.Status      `json:"status"`
					Results interface{} `json:"results"`
				}{
					Status: response.Status{
						Code:        400,
						Description: "Invalid API key",
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewRajaOngkirClient("invalid-key", server.URL)
		result, err := client.GetProvinces("")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Invalid API key")
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
	t.Run("Success - Get all cities", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/city", r.URL.Path)
			assert.Equal(t, "test-api-key", r.Header.Get("key"))

			response := response.RajaOngkirResponse{
				RajaOngkir: struct {
					Query   interface{} `json:"query"`
					Status  response.Status      `json:"status"`
					Results interface{} `json:"results"`
				}{
					Status: response.Status{
						Code:        200,
						Description: "OK",
					},
					Results: []response.RajaOngkirCity{
						{
							CityID:     "1",
							ProvinceID: "1",
							Province:   "Bali",
							Type:       "Kabupaten",
							CityName:   "Badung",
							PostalCode: "80351",
						},
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewRajaOngkirClient("test-api-key", server.URL)
		result, err := client.GetCities("", "")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "1", result[0].CityID)
		assert.Equal(t, "Badung", result[0].CityName)
	})

	t.Run("Success - Get cities by province", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "province=1", r.URL.RawQuery)

			response := response.RajaOngkirResponse{
				RajaOngkir: struct {
					Query   interface{} `json:"query"`
					Status  response.Status      `json:"status"`
					Results interface{} `json:"results"`
				}{
					Status: response.Status{
						Code:        200,
						Description: "OK",
					},
					Results: []response.RajaOngkirCity{
						{
							CityID:     "1",
							ProvinceID: "1",
							Province:   "Bali",
							Type:       "Kabupaten",
							CityName:   "Badung",
							PostalCode: "80351",
						},
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewRajaOngkirClient("test-api-key", server.URL)
		result, err := client.GetCities("1", "")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "1", result[0].ProvinceID)
	})

	t.Run("Error - API returns error status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := response.RajaOngkirResponse{
				RajaOngkir: struct {
					Query   interface{} `json:"query"`
					Status  response.Status      `json:"status"`
					Results interface{} `json:"results"`
				}{
					Status: response.Status{
						Code:        400,
						Description: "Invalid province ID",
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewRajaOngkirClient("test-api-key", server.URL)
		result, err := client.GetCities("999", "")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Invalid province ID")
	})
}

// Test CalculateShippingCost method
func TestRajaOngkirClient_CalculateShippingCost(t *testing.T) {
	t.Run("Success - Calculate shipping cost", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/cost", r.URL.Path)
			assert.Equal(t, "test-api-key", r.Header.Get("key"))
			assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("content-type"))

			// Parse form data
			r.ParseForm()
			assert.Equal(t, "501", r.FormValue("origin"))
			assert.Equal(t, "114", r.FormValue("destination"))
			assert.Equal(t, "1000", r.FormValue("weight"))
			assert.Equal(t, "jne", r.FormValue("courier"))

			response := response.RajaOngkirResponse{
				RajaOngkir: struct {
					Query   interface{} `json:"query"`
					Status  response.Status      `json:"status"`
					Results interface{} `json:"results"`
				}{
					Status: response.Status{
						Code:        200,
						Description: "OK",
					},
					Results: []response.RajaOngkirCost{
						{
							Code: "jne",
							Name: "Jalur Nugraha Ekakurir (JNE)",
							Costs: []struct {
								Service     string `json:"service"`
								Description string `json:"description"`
								Cost        []struct {
									Value int    `json:"value"`
									ETD   string `json:"etd"`
									Note  string `json:"note"`
								} `json:"cost"`
							}{
								{
									Service:     "REG",
									Description: "Layanan Reguler",
									Cost: []struct {
										Value int    `json:"value"`
										ETD   string `json:"etd"`
										Note  string `json:"note"`
									}{
										{
											Value: 15000,
											ETD:   "1-2",
											Note:  "",
										},
									},
								},
							},
						},
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
			response := response.RajaOngkirResponse{
				RajaOngkir: struct {
					Query   interface{} `json:"query"`
					Status  response.Status      `json:"status"`
					Results interface{} `json:"results"`
				}{
					Status: response.Status{
						Code:        400,
						Description: "Invalid courier",
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewRajaOngkirClient("test-api-key", server.URL)
		result, err := client.CalculateShippingCost("501", "114", 1000, "invalid")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "Invalid courier")
	})

	t.Run("Error - Server unavailable", func(t *testing.T) {
		client := NewRajaOngkirClient("test-api-key", "http://invalid-url")
		result, err := client.CalculateShippingCost("501", "114", 1000, "jne")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to make request")
	})
}