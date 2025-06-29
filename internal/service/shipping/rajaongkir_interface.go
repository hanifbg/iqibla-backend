package shipping

import "github.com/hanifbg/landing_backend/internal/model/response"

//go:generate mockgen -source=rajaongkir_interface.go -destination=mocks/rajaongkir_client_mock.go -package=mocks

type RajaOngkirClientInterface interface {
	GetProvinces(provinceID string) ([]response.RajaOngkirProvince, error)
	GetCities(provinceID, cityID string) ([]response.RajaOngkirCity, error)
	CalculateShippingCost(origin, destination string, weight int, courier string) ([]response.RajaOngkirCost, error)
}