package response

import (
	"time"

	"github.com/hanifbg/landing_backend/internal/model/entity"
)

type OrderItemResponse struct {
	ID               string                 `json:"id"`
	ProductVariantID string                 `json:"product_variant_id"`
	VariantName      string                 `json:"variant_name"`
	Quantity         int                    `json:"quantity"`
	UnitPrice        float64                `json:"unit_price"`
	Subtotal         float64                `json:"subtotal"`
	ImageURL         string                 `json:"image_url"`
	Attributes       map[string]interface{} `json:"attributes,omitempty"`
}

type OrderResponse struct {
	ID              string              `json:"id"`
	CartID          string              `json:"cart_id"`
	CustomerName    string              `json:"customer_name"`
	CustomerEmail   string              `json:"customer_email"`
	CustomerPhone   string              `json:"customer_phone"`
	ShippingAddress string              `json:"shipping_address"`
	Subtotal        float64             `json:"subtotal"`
	DiscountAmount  float64             `json:"discount_amount"`
	DiscountCode    string              `json:"discount_code,omitempty"`
	ShippingCost    float64             `json:"shipping_cost"`
	TotalAmount     float64             `json:"total_amount"`
	OrderStatus     string              `json:"order_status"`
	Notes           string              `json:"notes,omitempty"`
	Items           []OrderItemResponse `json:"items"`
	CreatedAt       time.Time           `json:"created_at"`
}

type PaymentResponse struct {
	ID              string             `json:"id"`
	OrderID         string             `json:"order_id"`
	Amount          float64            `json:"amount"`
	Status          entity.PaymentStatus `json:"status"`
	PaymentMethod   string             `json:"payment_method,omitempty"`
	TransactionID   string             `json:"transaction_id,omitempty"`
	PaymentToken    string             `json:"payment_token,omitempty"`
	PaymentURL      string             `json:"payment_url,omitempty"`
	ExpiryTime      *time.Time         `json:"expiry_time,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
}

type PaymentStatusResponse struct {
	ID              string             `json:"id"`
	OrderID         string             `json:"order_id"`
	Status          entity.PaymentStatus `json:"status"`
	TransactionID   string             `json:"transaction_id,omitempty"`
	TransactionTime *time.Time         `json:"transaction_time,omitempty"`
	PaymentMethod   string             `json:"payment_method,omitempty"`
	Amount          float64            `json:"amount"`
	UpdatedAt       time.Time          `json:"updated_at"`
}