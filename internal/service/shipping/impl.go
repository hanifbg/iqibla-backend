package shipping

import (
	"fmt"
	"strconv"

	"github.com/hanifbg/landing_backend/internal/model/request"
	"github.com/hanifbg/landing_backend/internal/model/response"
)

func (s *ShippingService) GetProvinces(req request.GetProvincesRequest) ([]response.ProvinceResponse, error) {
	// Validate input if necessary

	// Call repository
	provinces, err := s.shippingRepo.GetProvinces(req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get provinces: %w", err)
	}

	// Transform to service response format
	var result []response.ProvinceResponse
	for _, province := range provinces {
		result = append(result, response.ProvinceResponse{
			ProvinceID: strconv.Itoa(province.ProvinceID),
			Province:   province.Province,
		})
	}

	return result, nil
}

func (s *ShippingService) GetCities(req request.GetCitiesRequest) ([]response.CityResponse, error) {
	// Validate input if necessary

	// Call repository
	cities, err := s.shippingRepo.GetCities(req.ProvinceID, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cities: %w", err)
	}

	// Transform to service response format
	var result []response.CityResponse
	for _, city := range cities {
		result = append(result, response.CityResponse{
			CityID:     strconv.Itoa(city.CityID),
			ProvinceID: strconv.Itoa(city.ProvinceID),
			CityName:   city.CityName,
		})
	}

	return result, nil
}

func (s *ShippingService) GetDistricts(req request.GetDistrictsRequest) ([]response.DistrictResponse, error) {
	// Validate input
	if req.CityID == "" {
		return nil, fmt.Errorf("city ID is required")
	}

	// Call repository
	districts, err := s.shippingRepo.GetDistricts(req.CityID)
	if err != nil {
		return nil, err
	}

	var result []response.DistrictResponse
	for _, district := range districts {
		result = append(result, response.DistrictResponse{
			DistrictID:   strconv.Itoa(district.DistrictID),
			CityID:       req.CityID, // Use the city ID from the request parameter
			DistrictName: district.DistrictName,
		})
	}

	return result, nil
}

func (s *ShippingService) CalculateShippingCost(req request.CalculateShippingRequest) ([]response.ShippingCostResponse, error) {
	// Validate input
	if req.Origin == "" || req.Destination == "" || req.Weight <= 0 || req.Courier == "" {
		return nil, fmt.Errorf("origin, destination, weight, and courier are required")
	}

	// Call repository
	costs, err := s.shippingRepo.CalculateShippingCost(req.Origin, req.Destination, req.Weight, req.Courier)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate shipping cost: %w", err)
	}

	var result []response.ShippingCostResponse
	for _, cost := range costs {
		for _, costDetail := range cost.Costs {
			for _, costValue := range costDetail.Cost {
				result = append(result, response.ShippingCostResponse{
					Service:     costDetail.Service,
					Description: costDetail.Description,
					Cost:        float64(costValue.Value),
					ETD:         costValue.ETD,
				})
			}
		}
	}

	return result, nil
}
