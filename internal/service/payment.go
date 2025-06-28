package service

import (
	"github.com/hanifbg/landing_backend/internal/model/request"
	"github.com/hanifbg/landing_backend/internal/model/response"
)

type PaymentService interface {
	// Order operations
	CreateOrder(req request.CreateOrderRequest) (*response.OrderResponse, error)
	GetOrder(orderID string) (*response.OrderResponse, error)
	
	// Payment operations
	CreatePayment(orderID string) (*response.PaymentResponse, error)
	GetPaymentStatus(paymentID string) (*response.PaymentStatusResponse, error)
	HandlePaymentNotification(notificationData request.PaymentNotificationRequest) error
}