package service

import (
	"github.com/hanifbg/landing_backend/internal/model/request"
	"github.com/hanifbg/landing_backend/internal/model/response"
)

type ShippingService interface {
	GetProvinces(req request.GetProvincesRequest) ([]response.ProvinceResponse, error)
	GetCities(req request.GetCitiesRequest) ([]response.CityResponse, error)
	GetDistricts(req request.GetDistrictsRequest) ([]response.DistrictResponse, error)
	CalculateShippingCost(req request.CalculateShippingRequest) ([]response.ShippingCostResponse, error)
	ValidateAndSaveAWB(req request.ValidateAWBRequest) (*response.ValidateAWBResponse, error)
}
