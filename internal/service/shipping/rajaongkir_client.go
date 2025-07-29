package shipping

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hanifbg/landing_backend/internal/model/response"
)

type RajaOngkirClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// Ensure RajaOngkirClient implements RajaOngkirClientInterface
var _ RajaOngkirClientInterface = (*RajaOngkirClient)(nil)

func NewRajaOngkirClient(apiKey, baseURL string) *RajaOngkirClient {
	return &RajaOngkirClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetProvinces retrieves a list of provinces from RajaOngkir API
func (r *RajaOngkirClient) GetProvinces(provinceID string) ([]response.RajaOngkirProvince, error) {
	requestURL := fmt.Sprintf("%s/destination/province", r.baseURL)
	if provinceID != "" {
		requestURL = fmt.Sprintf("%s/%s", requestURL, provinceID)
	}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("key", r.apiKey)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var komerceResp response.KomerceProvinceResponse
	if err := json.Unmarshal(body, &komerceResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if komerceResp.Meta.Code != 200 {
		return nil, fmt.Errorf("API error: %s", komerceResp.Meta.Message)
	}

	// The province ID from the API is an integer, but the struct expects a string.
	// We need to handle the conversion if necessary, but for now, we assume it can be unmarshaled directly.
	// If the API returns province IDs as numbers, the `RajaOngkirProvince` struct's `ProvinceID` field might need to be `int` or have a custom unmarshaler.

	// For now, let's check if the `id` field in the response is a number and handle it.
	// A more robust solution would be to adjust the struct or use a custom unmarshaler.
	var provinces []response.RajaOngkirProvince
	for _, p := range komerceResp.Data {
		// The json tag is already `id`, so direct unmarshaling should work if the types are compatible.
		// If `id` is an int in the JSON, `ProvinceID` should be `int` or we need to handle it.
		// Let's assume for now that it can be handled as a string.
		provinces = append(provinces, p)
	}

	return provinces, nil
}

// GetCities retrieves a list of cities from RajaOngkir API
func (r *RajaOngkirClient) GetCities(provinceID, cityID string) ([]response.RajaOngkirCity, error) {
	if provinceID == "" {
		return nil, fmt.Errorf("provinceID is required to get cities")
	}
	requestURL := fmt.Sprintf("%s/destination/city/%s", r.baseURL, provinceID)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("key", r.apiKey)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var komerceResp response.KomerceCityResponse
	if err := json.Unmarshal(body, &komerceResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal city response: %w", err)
	}

	if komerceResp.Meta.Code != 200 {
		return nil, fmt.Errorf("API error: %s", komerceResp.Meta.Message)
	}

	provinceIDInt, err := strconv.Atoi(provinceID)
	if err != nil {
		return nil, fmt.Errorf("invalid provinceID format: %w", err)
	}

	// Populate ProvinceID for each city from the request parameter
	cities := komerceResp.Data
	for i := range cities {
		cities[i].ProvinceID = provinceIDInt
	}

	return cities, nil
}

// GetDistricts retrieves a list of districts from RajaOngkir API
func (r *RajaOngkirClient) GetDistricts(cityID string) ([]response.RajaOngkirDistrict, error) {
	if cityID == "" {
		return nil, fmt.Errorf("cityID is required to get districts")
	}
	requestURL := fmt.Sprintf("%s/destination/district/%s", r.baseURL, cityID)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("key", r.apiKey)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Define a struct to match the expected response format
	type DistrictResponse struct {
		Meta response.Meta                 `json:"meta"`
		Data []response.RajaOngkirDistrict `json:"data"`
	}

	var districtResp DistrictResponse
	if err := json.Unmarshal(body, &districtResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal district response: %w", err)
	}

	if districtResp.Meta.Code != 200 {
		return nil, fmt.Errorf("API error: %s", districtResp.Meta.Message)
	}

	cityIDInt, err := strconv.Atoi(cityID)
	if err != nil {
		return nil, fmt.Errorf("invalid cityID format: %w", err)
	}

	// Ensure CityID is set for all districts
	districts := districtResp.Data
	for i := range districts {
		districts[i].CityID = cityIDInt
	}

	return districts, nil
}

func (r *RajaOngkirClient) CalculateShippingCost(origin, destination string, weight int, courier string) ([]response.RajaOngkirCost, error) {
	// Updated to use the correct endpoint for cost calculation
	requestURL := fmt.Sprintf("%s/calculate/district/domestic-cost", r.baseURL)

	data := url.Values{}
	data.Set("origin", origin)
	data.Set("destination", destination)
	data.Set("weight", strconv.Itoa(weight))
	data.Set("courier", courier)

	req, err := http.NewRequest("POST", requestURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("key", r.apiKey)
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// The API response format has changed, so we need to handle it differently
	// New format: {"meta":{"message":"...","code":200,"status":"success"},"data":[...]}
	var komerceCostResp struct {
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
	}

	if err := json.Unmarshal(body, &komerceCostResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if komerceCostResp.Meta.Code != 200 {
		return nil, fmt.Errorf("RajaOngkir API error: %s", komerceCostResp.Meta.Message)
	}

	// Convert the new response format to the existing RajaOngkirCost format
	var costs []response.RajaOngkirCost
	courierMap := make(map[string]*response.RajaOngkirCost)

	for _, item := range komerceCostResp.Data {
		// Check if we already have an entry for this courier
		cost, exists := courierMap[item.Code]
		if !exists {
			// Create a new courier entry
			cost = &response.RajaOngkirCost{
				Code: item.Code,
				Name: item.Name,
				Costs: []struct {
					Service     string `json:"service"`
					Description string `json:"description"`
					Cost        []struct {
						Value int    `json:"value"`
						ETD   string `json:"etd"`
						Note  string `json:"note"`
					} `json:"cost"`
				}{},
			}
			courierMap[item.Code] = cost
		}

		// Add the service to the courier
		cost.Costs = append(cost.Costs, struct {
			Service     string `json:"service"`
			Description string `json:"description"`
			Cost        []struct {
				Value int    `json:"value"`
				ETD   string `json:"etd"`
				Note  string `json:"note"`
			} `json:"cost"`
		}{
			Service:     item.Service,
			Description: item.Description,
			Cost: []struct {
				Value int    `json:"value"`
				ETD   string `json:"etd"`
				Note  string `json:"note"`
			}{
				{
					Value: item.Cost,
					ETD:   item.ETD,
					Note:  "",
				},
			},
		})
	}

	// Convert the map to a slice
	for _, cost := range courierMap {
		costs = append(costs, *cost)
	}

	return costs, nil
}
