package postgres

import (

	"github.com/hanifbg/landing_backend/internal/model/entity"
	"gorm.io/gorm"
)

type PaymentRepositoryImpl struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepositoryImpl {
	return &PaymentRepositoryImpl{db: db}
}

// Order operations
func (r *PaymentRepositoryImpl) CreateOrder(order *entity.Order) error {
	return r.db.Create(order).Error
}

func (r *PaymentRepositoryImpl) FindOrderByID(orderID string) (*entity.Order, error) {
	var order entity.Order
	if err := r.db.Where("id = ?", orderID).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *PaymentRepositoryImpl) UpdateOrderStatus(orderID, status string) error {
	return r.db.Model(&entity.Order{}).Where("id = ?", orderID).Update("order_status", status).Error
}

func (r *PaymentRepositoryImpl) GetOrderWithItems(orderID string) (*entity.Order, error) {
	var order entity.Order
	if err := r.db.Preload("OrderItems.ProductVariant").Where("id = ?", orderID).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// Order item operations
func (r *PaymentRepositoryImpl) CreateOrderItem(item *entity.OrderItem) error {
	return r.db.Create(item).Error
}

// Payment operations
func (r *PaymentRepositoryImpl) CreatePayment(payment *entity.Payment) error {
	return r.db.Create(payment).Error
}

func (r *PaymentRepositoryImpl) FindPaymentByID(paymentID string) (*entity.Payment, error) {
	var payment entity.Payment
	if err := r.db.Where("id = ?", paymentID).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepositoryImpl) FindPaymentByOrderID(orderID string) (*entity.Payment, error) {
	var payment entity.Payment
	if err := r.db.Where("order_id = ?", orderID).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepositoryImpl) FindPaymentByTransactionID(transactionID string) (*entity.Payment, error) {
	var payment entity.Payment
	if err := r.db.Where("transaction_id = ?", transactionID).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepositoryImpl) UpdatePaymentStatus(paymentID string, status entity.PaymentStatus) error {
	return r.db.Model(&entity.Payment{}).Where("id = ?", paymentID).Update("status", status).Error
}

func (r *PaymentRepositoryImpl) UpdatePayment(payment *entity.Payment) error {
	return r.db.Save(payment).Error
}