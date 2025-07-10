package category

import (
	"github.com/hanifbg/landing_backend/internal/service"
	"github.com/hanifbg/landing_backend/internal/service/util"
	"github.com/labstack/echo/v4"
)

type ApiWrapper struct {
	CategoryService service.CategoryService
}

func InitRoute(e *echo.Echo, servWrapper *util.ServiceWrapper) {
	api := ApiWrapper{
		CategoryService: servWrapper.CategoryService,
	}
	api.registerRouter(e)
}

func (h *ApiWrapper) registerRouter(e *echo.Echo) {
	categoryV1 := e.Group("api/v1/categories")
	categoryV1.GET("/:slug", h.GetCategoryBySlug)
}
