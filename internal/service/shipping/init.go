package shipping

import (
	"github.com/hanifbg/landing_backend/config"
	"github.com/hanifbg/landing_backend/internal/service"
)

type ShippingService struct {
	rajaOngkirClient RajaOngkirClientInterface
}

func New(cfg *config.AppConfig) service.ShippingService {
	rajaOngkirClient := NewRajaOngkirClient(cfg.RajaOngkirAPIKey, cfg.RajaOngkirBaseURL)
	return &ShippingService{
		rajaOngkirClient: rajaOngkirClient,
	}
}