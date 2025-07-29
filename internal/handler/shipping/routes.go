package shipping

import (
	"github.com/hanifbg/landing_backend/internal/service"
	"github.com/hanifbg/landing_backend/internal/service/util"
	"github.com/labstack/echo/v4"
)

type ApiWrapper struct {
	shippingService service.ShippingService
}

func InitRoute(e *echo.Echo, servWrapper *util.ServiceWrapper) {
	api := ApiWrapper{
		shippingService: servWrapper.ShippingService,
	}
	api.registerRouter(e)
}

func (h *ApiWrapper) registerRouter(e *echo.Echo) {
	shippingGroup := e.Group("/api/v1/shipping")
	shippingGroup.GET("/provinces", h.GetProvinces)
	shippingGroup.GET("/cities/:province_id", h.GetCities)
	shippingGroup.POST("/cost", h.CalculateShippingCost)
}
