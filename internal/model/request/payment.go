package request

type CreateOrderRequest struct {
	CartID             string  `json:"cart_id" validate:"required"`
	CustomerName       string  `json:"customer_name" validate:"required"`
	CustomerEmail      string  `json:"customer_email" validate:"required,email"`
	CustomerPhone      string  `json:"customer_phone" validate:"required"`
	ShippingAddress    string  `json:"shipping_address" validate:"required"`
	ShippingCityID     string  `json:"shipping_city_id" validate:"required"`
	ShippingProvinceID string  `json:"shipping_province_id" validate:"required"`
	ShippingPostalCode string  `json:"shipping_postal_code" validate:"required"`
	ShippingCourier    string  `json:"shipping_courier" validate:"required"`
	ShippingService    string  `json:"shipping_service" validate:"required"`
	ShippingCost       float64 `json:"shipping_cost" validate:"required"`
	TotalWeight        int     `json:"total_weight" validate:"required"`
	Notes              string  `json:"notes,omitempty"`
}

type PaymentNotificationRequest struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	Currency          string `json:"currency"`
}
