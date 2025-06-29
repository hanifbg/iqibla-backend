package shipping

import (
	"github.com/labstack/echo/v4"
)

func InitRoute(e *echo.Echo) {
	e.POST("/api/shipping/rates", GetShippingRate)
}
