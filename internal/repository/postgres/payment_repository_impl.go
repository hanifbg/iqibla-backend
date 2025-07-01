package postgres

import (
	"github.com/hanifbg/landing_backend/internal/model/entity"
	"gorm.io/gorm"
)

// Order operations
func (r *RepoDatabase) CreateOrder(order *entity.Order) error {
	return r.DB.Create(order).Error
}

func (r *RepoDatabase) FindOrderByID(orderID string) (*entity.Order, error) {
	var order entity.Order
	if err := r.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *RepoDatabase) UpdateOrderStatus(orderID, status string) error {
	return r.DB.Model(&entity.Order{}).Where("id = ?", orderID).Update("order_status", status).Error
}

func (r *RepoDatabase) GetOrderWithItems(orderID string) (*entity.Order, error) {
	var order entity.Order
	if err := r.DB.Preload("OrderItems.ProductVariant").Where("id = ?", orderID).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// Order item operations
func (r *RepoDatabase) CreateOrderItem(item *entity.OrderItem) error {
	return r.DB.Create(item).Error
}

// Payment operations
func (r *RepoDatabase) CreatePayment(payment *entity.Payment) error {
	return r.DB.Create(payment).Error
}

func (r *RepoDatabase) FindPaymentByID(paymentID string) (*entity.Payment, error) {
	var payment entity.Payment
	if err := r.DB.Where("id = ?", paymentID).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *RepoDatabase) FindPaymentByOrderID(orderID string) (*entity.Payment, error) {
	var payment entity.Payment
	if err := r.DB.Where("order_id = ?", orderID).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *RepoDatabase) FindPaymentByTransactionID(transactionID string) (*entity.Payment, error) {
	var payment entity.Payment
	if err := r.DB.Where("transaction_id = ?", transactionID).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *RepoDatabase) UpdatePaymentStatus(paymentID string, status entity.PaymentStatus) error {
	return r.DB.Model(&entity.Payment{}).Where("id = ?", paymentID).Update("status", status).Error
}

func (r *RepoDatabase) UpdatePayment(payment *entity.Payment) error {
	return r.DB.Save(payment).Error
}

// Transaction operations
func (r *RepoDatabase) CreateOrderWithItems(order *entity.Order, items []entity.OrderItem) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Create order
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// Create order items
		for i := range items {
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *RepoDatabase) UpdatePaymentAndOrderStatus(payment *entity.Payment, orderID, orderStatus string) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// Update payment
		if err := tx.Save(payment).Error; err != nil {
			return err
		}

		// Update order status
		if err := tx.Model(&entity.Order{}).Where("id = ?", orderID).Update("order_status", orderStatus).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *RepoDatabase) GetSeq() (int64, error) {
	var nextSeq int64
	err := r.DB.Raw("SELECT nextval('order_number_seq')").Scan(&nextSeq).Error
	if err != nil {
		return 0, err
	}

	return nextSeq, nil
}
