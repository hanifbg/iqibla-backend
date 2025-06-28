package util

import (
	"github.com/hanifbg/landing_backend/config"
	"github.com/hanifbg/landing_backend/internal/handler/cart"
	"github.com/hanifbg/landing_backend/internal/handler/payment"
	"github.com/hanifbg/landing_backend/internal/handler/product"
	"github.com/hanifbg/landing_backend/internal/handler/swagger"
	serv "github.com/hanifbg/landing_backend/internal/service/util"
	"github.com/labstack/echo/v4"
)

func InitHandler(cfg *config.AppConfig, e *echo.Echo, servWrapper *serv.ServiceWrapper) {
	// Initialize product routes
	product.InitRoute(e, servWrapper)

	// Initialize cart routes
	cart.InitRoute(e, servWrapper)

	// Initialize payment routes
	payment.RegisterRoutes(e, servWrapper.PaymentService)

	// Init swagger
	swagger.InitRoute(e)
}
