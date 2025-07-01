package payment

import (
	"encoding/json"
	"net/http"

	"github.com/hanifbg/landing_backend/internal/model/request"
	"github.com/hanifbg/landing_backend/internal/service"
	"github.com/labstack/echo/v4"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// CreateOrder godoc
// @Summary Create a new order from cart
// @Description Create a new order from the items in the cart
// @Tags orders
// @Accept json
// @Produce json
// @Param request body request.CreateOrderRequest true "Order details"
// @Success 200 {object} response.OrderResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/orders [post]
func (h *PaymentHandler) CreateOrder(c echo.Context) error {
	var req request.CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request: " + err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Validation error: " + err.Error(),
		})
	}

	order, err := h.paymentService.CreateOrder(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to create order: " + err.Error(),
		})
	}

	order.Message = "Order created successfully"
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	return c.JSON(http.StatusOK, order)
}

// GetOrder godoc
// @Summary Get order details
// @Description Get details of an order by ID
// @Tags orders
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Success 200 {object} response.OrderResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/orders/{order_id} [get]
func (h *PaymentHandler) GetOrder(c echo.Context) error {
	orderID := c.Param("order_id")
	if orderID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Order ID is required",
		})
	}

	order, err := h.paymentService.GetOrder(orderID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "Order not found: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, order)
}

// CreatePayment godoc
// @Summary Create payment for an order
// @Description Create a payment transaction for an order
// @Tags payments
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Success 200 {object} response.PaymentResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/payments/{order_id} [post]
func (h *PaymentHandler) CreatePayment(c echo.Context) error {
	orderID := c.Param("order_id")
	if orderID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Order ID is required",
		})
	}

	payment, err := h.paymentService.CreatePayment(orderID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to create payment: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, payment)
}

// GetPaymentStatus godoc
// @Summary Get payment status
// @Description Get the status of a payment by ID
// @Tags payments
// @Accept json
// @Produce json
// @Param payment_id path string true "Payment ID"
// @Success 200 {object} response.PaymentStatusResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/payments/status/{payment_id} [get]
func (h *PaymentHandler) GetPaymentStatus(c echo.Context) error {
	paymentID := c.Param("payment_id")
	if paymentID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Payment ID is required",
		})
	}

	status, err := h.paymentService.GetPaymentStatus(paymentID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "Payment not found: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, status)
}

// HandleNotification godoc
// @Summary Handle payment notification
// @Description Handle payment notification from payment gateway
// @Tags payments
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/payments/notification [post]
func (h *PaymentHandler) HandleNotification(c echo.Context) error {
	// Parse notification data
	var notificationData request.PaymentNotificationRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&notificationData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid notification data: " + err.Error(),
		})
	}

	// Process notification
	if err := h.paymentService.HandlePaymentNotification(notificationData); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to process notification: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
