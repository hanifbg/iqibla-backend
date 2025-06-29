package service

import (
	"github.com/hanifbg/landing_backend/internal/model/request"
	"github.com/hanifbg/landing_backend/internal/model/response"
)

type ShippingService interface {
	GetProvinces(req request.GetProvincesRequest) ([]response.ProvinceResponse, error)
	GetCities(req request.GetCitiesRequest) ([]response.CityResponse, error)
	CalculateShippingCost(req request.CalculateShippingRequest) ([]response.ShippingCostResponse, error)
}