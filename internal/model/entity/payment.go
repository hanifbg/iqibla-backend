package entity

import (
	"time"

	"gorm.io/gorm"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

// Payment status constants
const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSuccess   PaymentStatus = "success"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusExpired   PaymentStatus = "expired"
	PaymentStatusCancelled PaymentStatus = "cancelled"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// PaymentMethod represents the payment method used
type PaymentMethod string

// Payment method constants
const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodEWallet PaymentMethod = "e_wallet"
	PaymentMethodQRIS PaymentMethod = "qris"
	PaymentMethodRetailOutlet PaymentMethod = "retail_outlet"
)

// Payment represents a payment transaction
type Payment struct {
	ID               string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	OrderID          string         `gorm:"type:uuid;not null;index" json:"order_id"`
	Order            *Order         `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Amount           float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status           PaymentStatus  `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	PaymentMethod    PaymentMethod  `gorm:"type:varchar(50)" json:"payment_method,omitempty"`
	TransactionID    string         `gorm:"type:varchar(100);index" json:"transaction_id,omitempty"`
	TransactionTime  *time.Time     `json:"transaction_time,omitempty"`
	PaymentToken     string         `gorm:"type:text" json:"payment_token,omitempty"`
	PaymentURL       string         `gorm:"type:text" json:"payment_url,omitempty"`
	ExpiryTime       *time.Time     `json:"expiry_time,omitempty"`
	PaymentDetails   JSONMap        `gorm:"type:jsonb" json:"payment_details,omitempty"`
	CreatedAt        time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// Order represents an order in the system
type Order struct {
	ID              string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CartID          string         `gorm:"type:uuid;not null;index" json:"cart_id"`
	CustomerID      *string        `gorm:"type:uuid;index" json:"customer_id,omitempty"`
	CustomerName    string         `gorm:"type:varchar(255)" json:"customer_name"`
	CustomerEmail   string         `gorm:"type:varchar(255)" json:"customer_email"`
	CustomerPhone   string         `gorm:"type:varchar(20)" json:"customer_phone"`
	ShippingAddress string         `gorm:"type:text" json:"shipping_address"`
	OrderItems      []OrderItem    `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
	Subtotal        float64        `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	DiscountAmount  float64        `gorm:"type:decimal(10,2);default:0" json:"discount_amount"`
	DiscountCode    string         `gorm:"type:varchar(50)" json:"discount_code,omitempty"`
	ShippingCost    float64        `gorm:"type:decimal(10,2);default:0" json:"shipping_cost"`
	TotalAmount     float64        `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	OrderStatus     string         `gorm:"type:varchar(20);not null;default:'pending'" json:"order_status"`
	Notes           string         `gorm:"type:text" json:"notes,omitempty"`
	Payment         *Payment       `gorm:"foreignKey:OrderID" json:"payment,omitempty"`
	CreatedAt       time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID               string          `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	OrderID          string          `gorm:"type:uuid;not null" json:"order_id"`
	ProductVariantID string          `gorm:"type:uuid;not null" json:"product_variant_id"`
	ProductVariant   *ProductVariant `gorm:"foreignKey:ProductVariantID" json:"product_variant,omitempty"`
	Quantity         int             `gorm:"not null" json:"quantity"`
	UnitPrice        float64         `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Subtotal         float64         `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	CreatedAt        time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt        time.Time       `gorm:"not null" json:"updated_at"`
	DeletedAt        gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
}