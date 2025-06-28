package repository

import "github.com/hanifbg/landing_backend/internal/model/entity"

type PaymentRepository interface {
	// Order operations
	CreateOrder(order *entity.Order) error
	FindOrderByID(orderID string) (*entity.Order, error)
	UpdateOrderStatus(orderID, status string) error
	GetOrderWithItems(orderID string) (*entity.Order, error)
	
	// Order item operations
	CreateOrderItem(item *entity.OrderItem) error
	
	// Payment operations
	CreatePayment(payment *entity.Payment) error
	FindPaymentByID(paymentID string) (*entity.Payment, error)
	FindPaymentByOrderID(orderID string) (*entity.Payment, error)
	FindPaymentByTransactionID(transactionID string) (*entity.Payment, error)
	UpdatePaymentStatus(paymentID string, status entity.PaymentStatus) error
	UpdatePayment(payment *entity.Payment) error
}