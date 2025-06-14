package util

import (
	"github.com/hanifbg/landing_backend/config"
	"github.com/hanifbg/landing_backend/internal/handler/product"
	serv "github.com/hanifbg/landing_backend/internal/service/util"
	"github.com/labstack/echo/v4"
)

func InitHandler(cfg *config.AppConfig, e *echo.Echo, servWrapper *serv.ServiceWrapper) {
	// var err error
	product.InitRoute(e, servWrapper)
}
