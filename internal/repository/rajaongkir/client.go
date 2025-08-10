package rajaongkir

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
	"github.com/hanifbg/landing_backend/internal/repository"
)

// Config holds configuration for the RajaOngkir repository
type Config struct {
	APIKey  string        // API key for RajaOngkir service
	BaseURL string        // Base URL for RajaOngkir API
	Timeout time.Duration // HTTP client timeout

	// Cache configuration
	CacheEnabled      bool // Whether caching is enabled
	CacheTTLHours     int  // Cache TTL in hours
	WarmupOnStartup   bool // Whether to warm up cache on startup
	WarmupTimeoutSecs int  // Warmup timeout in seconds
}

// Repository implements the ShippingRepository interface for RajaOngkir
type Repository struct {
	apiKey  string
	baseURL string
	client  *http.Client
	cache   *Cache
}

// Option is a functional option for configuring the Repository
type Option func(*Repository)

// WithTimeout sets a custom timeout for the HTTP client
func WithTimeout(timeout time.Duration) Option {
	return func(r *Repository) {
		r.client.Timeout = timeout
	}
}

// WithCustomClient sets a custom HTTP client
func WithCustomClient(client *http.Client) Option {
	return func(r *Repository) {
		r.client = client
	}
}

// NewRepository creates a new RajaOngkir repository
func NewRepository(cfg Config, opts ...Option) *Repository {
	// Set default timeout if not provided
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}

	// Initialize cache if enabled
	var cache *Cache
	if cfg.CacheEnabled {
		cacheTTL := time.Duration(cfg.CacheTTLHours) * time.Hour
		if cacheTTL == 0 {
			cacheTTL = 24 * time.Hour // Default to 24 hours
		}
		cache = NewCache(cacheTTL)
		fmt.Printf("🔧 CACHE INIT: Cache enabled with TTL=%v\n", cacheTTL)
	} else {
		fmt.Printf("🔧 CACHE INIT: Cache disabled\n")
	}

	// Create repository with defaults
	repo := &Repository{
		apiKey:  cfg.APIKey,
		baseURL: cfg.BaseURL,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
		cache: cache,
	}

	// Apply options
	for _, opt := range opts {
		opt(repo)
	}

	return repo
}

// GetProvinces retrieves a list of provinces from RajaOngkir API
func (r *Repository) GetProvinces(provinceID string) ([]response.RajaOngkirProvince, error) {
	// Check cache first if enabled and provinceID is empty (getting all provinces)
	if r.cache != nil && provinceID == "" {
		if cachedProvinces, found := r.cache.GetProvinces(); found {
			return cachedProvinces, nil
		}
	}

	requestURL := fmt.Sprintf("%s/destination/province", r.baseURL)
	if provinceID != "" {
		requestURL = fmt.Sprintf("%s/%s", requestURL, provinceID)
	}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, &repository.ShippingError{
			Operation: "GetProvinces.CreateRequest",
			Err:       err,
		}
	}

	req.Header.Set("key", r.apiKey)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, &repository.ShippingError{
			Operation: "GetProvinces.SendRequest",
			Err:       err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &repository.ShippingError{
			Operation: "GetProvinces.InvalidStatusCode",
			Err:       fmt.Errorf("received non-200 status code: %d", resp.StatusCode),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &repository.ShippingError{
			Operation: "GetProvinces.ReadResponseBody",
			Err:       err,
		}
	}

	// Parse response
	var rajaOngkirResp response.KomerceProvinceResponse
	if err := json.Unmarshal(body, &rajaOngkirResp); err != nil {
		return nil, &repository.ShippingError{
			Operation: "GetProvinces.ParseResponse",
			Err:       err,
		}
	}

	// Check response status
	if rajaOngkirResp.Meta.Code != http.StatusOK {
		return nil, &repository.ShippingError{
			Operation: "GetProvinces.APIError",
			Err:       fmt.Errorf("API error: %s (code: %d)", rajaOngkirResp.Meta.Message, rajaOngkirResp.Meta.Code),
		}
	}

	fmt.Printf("🌐 API CALL: Fetched %d provinces from RajaOngkir API\n", len(rajaOngkirResp.Data))

	// Cache the result if enabled and getting all provinces
	if r.cache != nil && provinceID == "" {
		r.cache.SetProvinces(rajaOngkirResp.Data)
	}

	return rajaOngkirResp.Data, nil
}

// GetCities retrieves a list of cities from RajaOngkir API
func (r *Repository) GetCities(provinceID, cityID string) ([]response.RajaOngkirCity, error) {
	// Check cache first if enabled and getting all cities for a province (no specific cityID)
	if r.cache != nil && provinceID != "" && cityID == "" {
		if cachedCities, found := r.cache.GetCities(provinceID); found {
			return cachedCities, nil
		}
	}

	// Construct the URL path - for province 1, it should be /destination/city/1
	requestURL := fmt.Sprintf("%s/destination/city", r.baseURL)

	// First prioritize province in the path
	if provinceID != "" {
		requestURL = fmt.Sprintf("%s/%s", requestURL, provinceID)
	}

	// Add cityID as a query parameter if provided
	params := url.Values{}
	if cityID != "" {
		params.Add("id", cityID)
	}

	if len(params) > 0 {
		requestURL = fmt.Sprintf("%s?%s", requestURL, params.Encode())
	}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, &repository.ShippingError{
			Operation: "GetCities.CreateRequest",
			Err:       err,
		}
	}

	req.Header.Set("key", r.apiKey)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, &repository.ShippingError{
			Operation: "GetCities.SendRequest",
			Err:       err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &repository.ShippingError{
			Operation: "GetCities.InvalidStatusCode",
			Err:       fmt.Errorf("received non-200 status code: %d", resp.StatusCode),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &repository.ShippingError{
			Operation: "GetCities.ReadResponseBody",
			Err:       err,
		}
	}

	// Parse response
	var rajaOngkirResp response.KomerceCityResponse
	if err := json.Unmarshal(body, &rajaOngkirResp); err != nil {
		return nil, &repository.ShippingError{
			Operation: "GetCities.ParseResponse",
			Err:       err,
		}
	}

	// Check response status
	if rajaOngkirResp.Meta.Code != http.StatusOK {
		return nil, &repository.ShippingError{
			Operation: "GetCities.APIError",
			Err:       fmt.Errorf("API error: %s (code: %d)", rajaOngkirResp.Meta.Message, rajaOngkirResp.Meta.Code),
		}
	}

	// Try to get province ID from path parameter and set it for each city
	parsedProvinceID := 0
	if provinceID != "" {
		parsed, err := strconv.Atoi(provinceID)
		if err == nil {
			parsedProvinceID = parsed
		}
	}

	// Set the province ID for each city based on the URL parameter
	if parsedProvinceID > 0 {
		for i := range rajaOngkirResp.Data {
			rajaOngkirResp.Data[i].ProvinceID = parsedProvinceID
		}
	}

	fmt.Printf("🌐 API CALL: Fetched %d cities from RajaOngkir API for province %s\n", len(rajaOngkirResp.Data), provinceID)

	// Cache the result if enabled and getting all cities for a province
	if r.cache != nil && provinceID != "" && cityID == "" {
		r.cache.SetCities(provinceID, rajaOngkirResp.Data)
	}

	return rajaOngkirResp.Data, nil
}

// GetDistricts retrieves a list of districts from RajaOngkir API
func (r *Repository) GetDistricts(cityID string) ([]response.RajaOngkirDistrict, error) {
	// Validate input
	if cityID == "" {
		return nil, &repository.ShippingError{
			Operation: "GetDistricts.ValidateInput",
			Err:       fmt.Errorf("city ID is required"),
		}
	}

	// Check cache first if enabled
	if r.cache != nil {
		if cachedDistricts, found := r.cache.GetDistricts(cityID); found {
			return cachedDistricts, nil
		}
	}

	requestURL := fmt.Sprintf("%s/destination/district/%s", r.baseURL, cityID)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, &repository.ShippingError{
			Operation: "GetDistricts.CreateRequest",
			Err:       err,
		}
	}

	req.Header.Set("key", r.apiKey)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, &repository.ShippingError{
			Operation: "GetDistricts.SendRequest",
			Err:       err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &repository.ShippingError{
			Operation: "GetDistricts.InvalidStatusCode",
			Err:       fmt.Errorf("received non-200 status code: %d", resp.StatusCode),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &repository.ShippingError{
			Operation: "GetDistricts.ReadResponseBody",
			Err:       err,
		}
	}

	var result struct {
		Meta struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Status  string `json:"status"`
		} `json:"meta"`
		Data []response.RajaOngkirDistrict `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, &repository.ShippingError{
			Operation: "GetDistricts.ParseResponse",
			Err:       err,
		}
	}

	if result.Meta.Code != http.StatusOK {
		return nil, &repository.ShippingError{
			Operation: "GetDistricts.APIError",
			Err:       fmt.Errorf("API error: %s (code: %d)", result.Meta.Message, result.Meta.Code),
		}
	}

	fmt.Printf("🌐 API CALL: Fetched %d districts from RajaOngkir API for city %s\n", len(result.Data), cityID)

	// Cache the result if enabled
	if r.cache != nil {
		r.cache.SetDistricts(cityID, result.Data)
	}

	return result.Data, nil
}

// CalculateShippingCost calculates shipping costs between origin and destination
func (r *Repository) CalculateShippingCost(origin, destination string, weight int, courier string) ([]response.RajaOngkirCost, error) {
	// Validate input
	if origin == "" || destination == "" || weight <= 0 || courier == "" {
		return nil, &repository.ShippingError{
			Operation: "CalculateShippingCost.ValidateInput",
			Err:       fmt.Errorf("invalid input: origin, destination, weight and courier are required"),
		}
	}

	requestURL := fmt.Sprintf("%s/calculate/district/domestic-cost", r.baseURL)

	formData := url.Values{}
	formData.Add("origin", origin)
	formData.Add("destination", destination)
	formData.Add("weight", strconv.Itoa(weight))
	formData.Add("courier", strings.ToLower(courier)) // Ensure courier is lowercase

	fmt.Printf("🌐 API REQUEST: POST %s with origin=%s, destination=%s, weight=%d, courier=%s\n",
		requestURL, origin, destination, weight, strings.ToLower(courier))

	req, err := http.NewRequest("POST", requestURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, &repository.ShippingError{
			Operation: "CalculateShippingCost.CreateRequest",
			Err:       err,
		}
	}

	req.Header.Set("key", r.apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	fmt.Printf("🌐 API CALL: Calculating shipping cost from %s to %s (weight: %dg, courier: %s)\n", origin, destination, weight, courier)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, &repository.ShippingError{
			Operation: "CalculateShippingCost.SendRequest",
			Err:       err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read the response body for more detailed error information
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ API ERROR: Status %d, Response: %s\n", resp.StatusCode, string(body))
		return nil, &repository.ShippingError{
			Operation: "CalculateShippingCost.InvalidStatusCode",
			Err:       fmt.Errorf("received non-200 status code: %d, response: %s", resp.StatusCode, string(body)),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &repository.ShippingError{
			Operation: "CalculateShippingCost.ReadResponseBody",
			Err:       err,
		}
	}

	fmt.Printf("📝 API RESPONSE: %s\n", string(body))

	var result struct {
		Meta struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Status  string `json:"status"`
		} `json:"meta"`
		Data []response.RajaOngkirCost `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, &repository.ShippingError{
			Operation: "CalculateShippingCost.ParseResponse",
			Err:       err,
		}
	}

	if result.Meta.Code != http.StatusOK {
		return nil, &repository.ShippingError{
			Operation: "CalculateShippingCost.APIError",
			Err:       fmt.Errorf("API error: %s (code: %d)", result.Meta.Message, result.Meta.Code),
		}
	}

	return result.Data, nil
}

// ValidateAWB validates AWB number with RajaOngkir tracking API
func (r *Repository) ValidateAWB(awbNumber, courier string, lastPhoneNumber *string) (*response.RajaOngkirTrackingResponse, error) {
	// Build the URL with query parameters
	endpoint := fmt.Sprintf("%s/track/waybill", r.baseURL)
	params := url.Values{}
	params.Add("awb", awbNumber)
	params.Add("courier", courier)

	if lastPhoneNumber != nil && *lastPhoneNumber != "" {
		params.Add("last_phone_number", *lastPhoneNumber)
	}

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	// Create POST request
	req, err := http.NewRequest("POST", fullURL, nil)
	if err != nil {
		return nil, &repository.ShippingError{
			Operation: "ValidateAWB.CreateRequest",
			Err:       err,
		}
	}

	// Set headers
	req.Header.Set("key", r.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, &repository.ShippingError{
			Operation: "ValidateAWB.HTTPRequest",
			Err:       err,
		}
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &repository.ShippingError{
			Operation: "ValidateAWB.ReadResponse",
			Err:       err,
		}
	}

	// Parse response
	var result response.RajaOngkirTrackingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, &repository.ShippingError{
			Operation: "ValidateAWB.ParseResponse",
			Err:       err,
		}
	}

	// Check if the response indicates an error
	if result.Meta.Code != http.StatusOK {
		return nil, &repository.ShippingError{
			Operation: "ValidateAWB.APIError",
			Err:       fmt.Errorf("invalid AWB: %s", result.Meta.Message),
		}
	}

	return &result, nil
}
