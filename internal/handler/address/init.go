package address

import (
	"github.com/labstack/echo/v4"
)

func InitRoute(e *echo.Echo) {
	e.GET("/api/address/search", SearchAddress)
}
