package shipping

import (
	"github.com/hanifbg/landing_backend/internal/model/request"
	"github.com/hanifbg/landing_backend/internal/model/response"
)

func (s *ShippingService) GetProvinces(req request.GetProvincesRequest) ([]response.ProvinceResponse, error) {
	provinces, err := s.rajaOngkirClient.GetProvinces(req.ID)
	if err != nil {
		return nil, err
	}

	var result []response.ProvinceResponse
	for _, province := range provinces {
		result = append(result, response.ProvinceResponse{
			ProvinceID: province.ProvinceID,
			Province:   province.Province,
		})
	}

	return result, nil
}

func (s *ShippingService) GetCities(req request.GetCitiesRequest) ([]response.CityResponse, error) {
	cities, err := s.rajaOngkirClient.GetCities(req.ProvinceID, req.ID)
	if err != nil {
		return nil, err
	}

	var result []response.CityResponse
	for _, city := range cities {
		result = append(result, response.CityResponse{
			CityID:     city.CityID,
			ProvinceID: city.ProvinceID,
			Province:   city.Province,
			Type:       city.Type,
			CityName:   city.CityName,
			PostalCode: city.PostalCode,
		})
	}

	return result, nil
}

func (s *ShippingService) CalculateShippingCost(req request.CalculateShippingRequest) ([]response.ShippingCostResponse, error) {
	costs, err := s.rajaOngkirClient.CalculateShippingCost(req.Origin, req.Destination, req.Weight, req.Courier)
	if err != nil {
		return nil, err
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