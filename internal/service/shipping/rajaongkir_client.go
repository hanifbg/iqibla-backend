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

func (r *RajaOngkirClient) GetProvinces(provinceID string) ([]response.RajaOngkirProvince, error) {
	requestURL := fmt.Sprintf("%s/province", r.baseURL)
	if provinceID != "" {
		requestURL += "?id=" + provinceID
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

	var rajaOngkirResp response.RajaOngkirResponse
	if err := json.Unmarshal(body, &rajaOngkirResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if rajaOngkirResp.RajaOngkir.Status.Code != 200 {
		return nil, fmt.Errorf("RajaOngkir API error: %s", rajaOngkirResp.RajaOngkir.Status.Description)
	}

	// Convert results to provinces
	resultsBytes, err := json.Marshal(rajaOngkirResp.RajaOngkir.Results)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal results: %w", err)
	}

	var provinces []response.RajaOngkirProvince
	// Try to unmarshal as array first
	if err := json.Unmarshal(resultsBytes, &provinces); err != nil {
		// If that fails, try to unmarshal as single object
		var singleProvince response.RajaOngkirProvince
		if err := json.Unmarshal(resultsBytes, &singleProvince); err != nil {
			return nil, fmt.Errorf("failed to unmarshal provinces: %w", err)
		}
		provinces = []response.RajaOngkirProvince{singleProvince}
	}

	return provinces, nil
}

func (r *RajaOngkirClient) GetCities(provinceID, cityID string) ([]response.RajaOngkirCity, error) {
	requestURL := fmt.Sprintf("%s/city", r.baseURL)
	params := url.Values{}
	if provinceID != "" {
		params.Add("province", provinceID)
	}
	if cityID != "" {
		params.Add("id", cityID)
	}
	if len(params) > 0 {
		requestURL += "?" + params.Encode()
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

	var rajaOngkirResp response.RajaOngkirResponse
	if err := json.Unmarshal(body, &rajaOngkirResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if rajaOngkirResp.RajaOngkir.Status.Code != 200 {
		return nil, fmt.Errorf("RajaOngkir API error: %s", rajaOngkirResp.RajaOngkir.Status.Description)
	}

	// Convert results to cities
	resultsBytes, err := json.Marshal(rajaOngkirResp.RajaOngkir.Results)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal results: %w", err)
	}

	var cities []response.RajaOngkirCity
	// Try to unmarshal as array first
	if err := json.Unmarshal(resultsBytes, &cities); err != nil {
		// If that fails, try to unmarshal as single object
		var singleCity response.RajaOngkirCity
		if err := json.Unmarshal(resultsBytes, &singleCity); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cities: %w", err)
		}
		cities = []response.RajaOngkirCity{singleCity}
	}

	return cities, nil
}

func (r *RajaOngkirClient) CalculateShippingCost(origin, destination string, weight int, courier string) ([]response.RajaOngkirCost, error) {
	requestURL := fmt.Sprintf("%s/cost", r.baseURL)

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

	var rajaOngkirResp response.RajaOngkirResponse
	if err := json.Unmarshal(body, &rajaOngkirResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if rajaOngkirResp.RajaOngkir.Status.Code != 200 {
		return nil, fmt.Errorf("RajaOngkir API error: %s", rajaOngkirResp.RajaOngkir.Status.Description)
	}

	// Convert results to costs
	resultsBytes, err := json.Marshal(rajaOngkirResp.RajaOngkir.Results)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal results: %w", err)
	}

	var costs []response.RajaOngkirCost
	if err := json.Unmarshal(resultsBytes, &costs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal costs: %w", err)
	}

	return costs, nil
}