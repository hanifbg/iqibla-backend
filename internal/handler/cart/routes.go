package cart

import (
	"github.com/hanifbg/landing_backend/internal/service"
	"github.com/hanifbg/landing_backend/internal/service/util"
	"github.com/labstack/echo/v4"
)

type ApiWrapper struct {
	cartService service.CartService
}

func InitRoute(e *echo.Echo, servWrapper *util.ServiceWrapper) {
	api := ApiWrapper{
		cartService: servWrapper.CartService,
	}
	api.registerRouter(e)
}

func (h *ApiWrapper) registerRouter(e *echo.Echo) {
	cartGroup := e.Group("/api/v1/cart")

	// Cart management endpoints
	cartGroup.POST("/add", h.AddItem)
	cartGroup.POST("/update-quantity", h.UpdateItemQuantity)
	cartGroup.POST("/remove", h.RemoveItem)
	cartGroup.GET("/:cart_id", h.GetCart)
	cartGroup.POST("/apply-discount", h.ApplyDiscount)
}
