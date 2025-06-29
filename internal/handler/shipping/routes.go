package shipping

import (
	"github.com/hanifbg/landing_backend/internal/handler"
	"github.com/hanifbg/landing_backend/internal/service"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, shippingService service.ShippingService) {
	shippingHandler := handler.NewShippingHandler(shippingService)

	// Shipping routes
	shippingGroup := e.Group("/api/v1/shipping")
	shippingGroup.GET("/provinces", shippingHandler.GetProvinces)
	shippingGroup.GET("/cities", shippingHandler.GetCities)
	shippingGroup.POST("/cost", shippingHandler.CalculateShippingCost)
}