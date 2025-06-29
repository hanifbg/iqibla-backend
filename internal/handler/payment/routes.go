package payment

import (
	"github.com/hanifbg/landing_backend/internal/service"
	"github.com/labstack/echo/v4"
)

func InitRoute(e *echo.Echo, paymentService service.PaymentService) {
	handler := NewPaymentHandler(paymentService)
	registerRouter(e, handler)
}

func RegisterRoutes(e *echo.Echo, paymentService service.PaymentService) {
	handler := NewPaymentHandler(paymentService)
	registerRouter(e, handler)
}

func registerRouter(e *echo.Echo, handler *PaymentHandler) {
	// Order routes
	orderGroup := e.Group("/api/v1/orders")
	orderGroup.POST("", handler.CreateOrder)
	orderGroup.GET("/:order_id", handler.GetOrder)

	// Payment routes
	paymentGroup := e.Group("/api/v1/payments")
	paymentGroup.POST("/:order_id", handler.CreatePayment)
	paymentGroup.GET("/status/:payment_id", handler.GetPaymentStatus)
	paymentGroup.POST("/notification", handler.HandleNotification)
}