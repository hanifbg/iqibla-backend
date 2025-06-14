package util

import (
	"github.com/hanifbg/landing_backend/config"
	"github.com/hanifbg/landing_backend/internal/repository/util"
	"github.com/hanifbg/landing_backend/internal/service"
	"github.com/hanifbg/landing_backend/internal/service/product"
)

type ServiceWrapper struct {
	ProductService service.ProductService
}

func New(cfg *config.AppConfig, repoWrapper *util.RepoWrapper) (serviceWrapper *ServiceWrapper, err error) {

	serviceWrapper = &ServiceWrapper{
		ProductService: product.New(cfg, repoWrapper),
	}

	return
}
