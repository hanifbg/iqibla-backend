package shipping

import (
	"github.com/hanifbg/landing_backend/config"
	"github.com/hanifbg/landing_backend/internal/repository"
	"github.com/hanifbg/landing_backend/internal/repository/util"
	"github.com/hanifbg/landing_backend/internal/service"
)

type ShippingService struct {
	ShippingRepo    repository.ShippingRepository
	AWBTrackingRepo repository.AWBTrackingRepository
}

func New(cfg *config.AppConfig, repoWrapper *util.RepoWrapper) service.ShippingService {
	return &ShippingService{
		ShippingRepo:    repoWrapper.ShippingRepo,
		AWBTrackingRepo: repoWrapper.AWBTrackingRepo,
	}
}
