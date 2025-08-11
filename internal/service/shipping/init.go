package shipping

import (
	"github.com/hanifbg/landing_backend/config"
	"github.com/hanifbg/landing_backend/internal/repository"
	"github.com/hanifbg/landing_backend/internal/repository/util"
	"github.com/hanifbg/landing_backend/internal/service"
)

type ShippingService struct {
	shippingRepo repository.ShippingRepository
}

func New(cfg *config.AppConfig, repoWrapper *util.RepoWrapper) service.ShippingService {
	return &ShippingService{
		shippingRepo: repoWrapper.ShippingRepo,
	}
}
