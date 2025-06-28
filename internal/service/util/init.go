package util

import (
	"github.com/hanifbg/landing_backend/config"
	"github.com/hanifbg/landing_backend/internal/repository/util"
	"github.com/hanifbg/landing_backend/internal/service"
	"github.com/hanifbg/landing_backend/internal/service/cart"
	"github.com/hanifbg/landing_backend/internal/service/payment"
	"github.com/hanifbg/landing_backend/internal/service/product"
)

type ServiceWrapper struct {
	ProductService service.ProductService
	CartService    service.CartService
	PaymentService service.PaymentService
}

func New(cfg *config.AppConfig, repoWrapper *util.RepoWrapper) (serviceWrapper *ServiceWrapper, err error) {
	serviceWrapper = &ServiceWrapper{
		ProductService: product.New(cfg, repoWrapper),
		CartService:    cart.New(cfg, repoWrapper),
		PaymentService: payment.NewPaymentServiceWithMidtrans(repoWrapper.PaymentRepo, repoWrapper.CartRepo, cfg.MidtransServerKey, cfg.IsProduction),
	}

	return
}
