package product

import (
	"github.com/hanifbg/landing_backend/internal/service"
	"github.com/hanifbg/landing_backend/internal/service/util"
	"github.com/labstack/echo/v4"
)

type ApiWrapper struct {
	ProductService service.ProductService
}

func InitRoute(e *echo.Echo, servWrapper *util.ServiceWrapper) {
	api := ApiWrapper{
		ProductService: servWrapper.ProductService,
	}
	api.registerRouter(e)
}

// registerRouter initialize url route mapping
func (h *ApiWrapper) registerRouter(e *echo.Echo) {
	const UPLOAD_DIR = "./uploads"
	e.Static("/uploads", UPLOAD_DIR)
	productV1 := e.Group("api/v1/products")
	productV1.GET("/:id", h.GetProductByID)
	productV1.GET("", h.GetAllProducts)
}
