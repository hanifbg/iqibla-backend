package repository

import "github.com/hanifbg/landing_backend/internal/model/response"

//go:generate mockgen -source=shipping.go -destination=../service/shipping/mocks/shipping_repository_mock.go -package=mocks

// ShippingRepository defines the interface for shipping data operations
// It abstracts the actual implementation details for retrieving shipping data
type ShippingRepository interface {
	// GetProvinces retrieves a list of provinces from the shipping provider
	// If provinceID is provided, it retrieves a specific province
	GetProvinces(provinceID string) ([]response.RajaOngkirProvince, error)

	// GetCities retrieves a list of cities from the shipping provider
	// If provinceID is provided, it retrieves cities in that province
	// If cityID is also provided, it retrieves a specific city
	GetCities(provinceID, cityID string) ([]response.RajaOngkirCity, error)

	// GetDistricts retrieves a list of districts from the shipping provider
	// cityID is required to specify which city's districts to retrieve
	GetDistricts(cityID string) ([]response.RajaOngkirDistrict, error)

	// CalculateShippingCost calculates shipping costs between origin and destination
	// origin and destination are location IDs
	// weight is in grams
	// courier is the shipping provider code (e.g., "jne", "pos", "tiki")
	CalculateShippingCost(origin, destination string, weight int, courier string) ([]response.RajaOngkirCost, error)
}

// ShippingError represents errors from the shipping repository
type ShippingError struct {
	Operation string // Operation that failed
	Err       error  // Original error
}

// Error returns the string representation of the error
func (e *ShippingError) Error() string {
	if e.Err != nil {
		return e.Operation + ": " + e.Err.Error()
	}
	return e.Operation
}

// Unwrap returns the underlying error
func (e *ShippingError) Unwrap() error {
	return e.Err
}
