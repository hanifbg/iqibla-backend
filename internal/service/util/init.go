package util

import (
	"github.com/hanifbg/landing_backend/config"
	"github.com/hanifbg/landing_backend/internal/repository/util"
	"github.com/hanifbg/landing_backend/internal/service"
	"github.com/hanifbg/landing_backend/internal/service/cart"
	"github.com/hanifbg/landing_backend/internal/service/payment"
	"github.com/hanifbg/landing_backend/internal/service/product"
	"github.com/hanifbg/landing_backend/internal/service/shipping"
)

type ServiceWrapper struct {
	ProductService  service.ProductService
	CartService     service.CartService
	PaymentService  service.PaymentService
	ShippingService service.ShippingService
}

func New(cfg *config.AppConfig, repoWrapper *util.RepoWrapper) (serviceWrapper *ServiceWrapper, err error) {
	serviceWrapper = &ServiceWrapper{
		ProductService:  product.New(cfg, repoWrapper),
		CartService:     cart.New(cfg, repoWrapper),
		PaymentService:  payment.New(cfg, repoWrapper),
		ShippingService: shipping.New(cfg),
	}

	return
}
