package swagger

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func InitRoute(e *echo.Echo) {
	registerRouter(e)
}

// registerRouter initialize url route mapping
func registerRouter(e *echo.Echo) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
