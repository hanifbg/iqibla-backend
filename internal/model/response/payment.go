package response

import (
	"time"

	"github.com/hanifbg/landing_backend/internal/model/entity"
)

type OrderItemResponse struct {
	ID               string  `json:"id"`
	ProductVariantID string  `json:"product_variant_id"`
	ProductName      string  `json:"product_name"`
	ProductImage     string  `json:"product_image"`
	Quantity         int     `json:"quantity"`
	PriceAtPurchase  float64 `json:"price_at_purchase"`
}

type CreateOrderResponse struct {
	OrderID     string  `json:"id"`
	OrderNumber string  `json:"order_number"`
	TotalAmount float64 `json:"total_amount"`
	Message     string  `json:"message"`
}

type OrderResponse struct {
	ID                          string              `json:"id"`
	OrderNumber                 string              `json:"order_number"`
	CartID                      string              `json:"cart_id"`
	CustomerName                string              `json:"customer_name"`
	CustomerEmail               string              `json:"customer_email"`
	CustomerPhone               string              `json:"customer_phone"`
	ShippingAddress             string              `json:"shipping_address"`
	ShippingCityID              string              `json:"shipping_city_id"`
	ShippingProvinceID          string              `json:"shipping_province_id"`
	ShippingPostalCode          string              `json:"shipping_postal_code"`
	Subtotal                    float64             `json:"subtotal"`
	DiscountAmount              float64             `json:"discount_amount"`
	DiscountCodeApplied         string              `json:"discount_code_applied,omitempty"`
	ShippingCost                float64             `json:"shipping_cost"`
	TotalAmount                 float64             `json:"total_amount"`
	Currency                    string              `json:"currency"`
	OrderStatus                 string              `json:"order_status"`
	PaymentProcessor            string              `json:"payment_processor,omitempty"`
	PaymentGatewayTransactionID string              `json:"payment_gateway_transaction_id,omitempty"`
	SourceChannel               string              `json:"source_channel"`
	Notes                       string              `json:"notes,omitempty"`
	OrderItems                  []OrderItemResponse `json:"order_items,omitempty"`
	Payment                     *PaymentResponse    `json:"payment,omitempty"`
	CreatedAt                   time.Time           `json:"created_at"`
	UpdatedAt                   time.Time           `json:"updated_at"`
}

type PaymentResponse struct {
	ID            string               `json:"id"`
	OrderID       string               `json:"order_id"`
	Amount        float64              `json:"amount"`
	Status        entity.PaymentStatus `json:"status"`
	PaymentMethod string               `json:"payment_method,omitempty"`
	TransactionID string               `json:"transaction_id,omitempty"`
	PaymentToken  string               `json:"payment_token,omitempty"`
	PaymentURL    string               `json:"payment_url,omitempty"`
	ExpiryTime    *time.Time           `json:"expiry_time,omitempty"`
	CreatedAt     time.Time            `json:"created_at"`
}

type PaymentStatusResponse struct {
	ID              string               `json:"id"`
	OrderID         string               `json:"order_id"`
	Status          entity.PaymentStatus `json:"status"`
	TransactionID   string               `json:"transaction_id,omitempty"`
	TransactionTime *time.Time           `json:"transaction_time,omitempty"`
	PaymentMethod   string               `json:"payment_method,omitempty"`
	Amount          float64              `json:"amount"`
	UpdatedAt       time.Time            `json:"updated_at"`
}
