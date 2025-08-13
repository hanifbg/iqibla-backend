package entity

import (
	"fmt"
	"strconv"
	"strings"
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
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodEWallet      PaymentMethod = "e_wallet"
	PaymentMethodQRIS         PaymentMethod = "qris"
	PaymentMethodRetailOutlet PaymentMethod = "retail_outlet"
)

// Payment represents a payment transaction
type Payment struct {
	ID              string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	OrderID         string         `gorm:"type:uuid;not null;index" json:"order_id"`
	Order           *Order         `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Amount          float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status          PaymentStatus  `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	PaymentMethod   PaymentMethod  `gorm:"type:varchar(50)" json:"payment_method,omitempty"`
	TransactionID   string         `gorm:"type:varchar(100);index" json:"transaction_id,omitempty"`
	TransactionTime *time.Time     `json:"transaction_time,omitempty"`
	PaymentToken    string         `gorm:"type:text" json:"payment_token,omitempty"`
	PaymentURL      string         `gorm:"type:text" json:"payment_url,omitempty"`
	ExpiryTime      *time.Time     `json:"expiry_time,omitempty"`
	PaymentDetails  JSONMap        `gorm:"type:jsonb" json:"payment_details,omitempty"`
	CreatedAt       time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// Order represents an order in the system
type Order struct {
	ID                          string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	OrderNumber                 string         `gorm:"type:varchar(50);uniqueIndex" json:"order_number"`
	CartID                      string         `gorm:"type:uuid;not null;index" json:"cart_id"`
	CustomerID                  *string        `gorm:"type:uuid;index" json:"customer_id,omitempty"`
	CustomerName                string         `gorm:"type:varchar(255)" json:"customer_name"`
	CustomerEmail               string         `gorm:"type:varchar(255)" json:"customer_email"`
	CustomerPhone               string         `gorm:"type:varchar(50)" json:"customer_phone"`
	ShippingStreetAddress       string         `gorm:"type:varchar(255)" json:"shipping_street_address"`
	ShippingCity                string         `gorm:"type:varchar(100)" json:"shipping_city"`
	ShippingProvince            string         `gorm:"type:varchar(100)" json:"shipping_province"`
	ShippingDistrict            string         `gorm:"type:varchar(100)" json:"shipping_district"`
	ShippingPostalCode          string         `gorm:"type:varchar(20)" json:"shipping_postal_code"`
	ShippingCountry             string         `gorm:"type:varchar(100)" json:"shipping_country"`
	ShippingCourier             string         `gorm:"type:varchar(100)" json:"shipping_courier"`
	ShippingService             string         `gorm:"type:varchar(100)" json:"shipping_service"`
	BillingAddressDetails       JSONMap        `gorm:"type:jsonb" json:"billing_address_details,omitempty"`
	OrderItems                  []OrderItem    `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
	Subtotal                    float64        `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	DiscountAmount              float64        `gorm:"type:decimal(10,2);default:0" json:"discount_amount"`
	DiscountCodeApplied         string         `gorm:"type:varchar(50)" json:"discount_code_applied,omitempty"`
	ShippingCost                float64        `gorm:"type:decimal(10,2);default:0" json:"shipping_cost"`
	TotalAmount                 float64        `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Currency                    string         `gorm:"type:varchar(3);default:'IDR'" json:"currency"`
	OrderStatus                 string         `gorm:"type:varchar(20);not null;default:'pending'" json:"order_status"`
	PaymentProcessor            string         `gorm:"type:varchar(50)" json:"payment_processor,omitempty"`
	PaymentGatewayTransactionID string         `gorm:"type:varchar(255)" json:"payment_gateway_transaction_id,omitempty"`
	SourceChannel               string         `gorm:"type:varchar(50);default:'web'" json:"source_channel"`
	Notes                       string         `gorm:"type:text" json:"notes,omitempty"`
	Payment                     *Payment       `gorm:"foreignKey:OrderID" json:"payment,omitempty"`
	CreatedAt                   time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt                   time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt                   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID               string          `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	OrderID          string          `gorm:"type:uuid;not null" json:"order_id"`
	ProductVariantID string          `gorm:"type:uuid;not null" json:"product_variant_id"`
	ProductVariant   *ProductVariant `gorm:"foreignKey:ProductVariantID" json:"product_variant,omitempty"`
	Quantity         int             `gorm:"not null" json:"quantity"`
	PriceAtPurchase  float64         `gorm:"type:decimal(10,2);not null" json:"price_at_purchase"`
	CreatedAt        time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt        time.Time       `gorm:"not null" json:"updated_at"`
	DeletedAt        gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
}

// FormatToIndonesianCurrency formats a float64 value to Indonesian currency format
// with dot (.) as thousand separator
// Example: 1000 -> 1.000, 100000 -> 100.000, 1234.56 -> 1.234,56
func FormatToIndonesianCurrency(amount float64) string {
	// Convert float to string with 2 decimal places
	amountStr := strconv.FormatFloat(amount, 'f', 2, 64)
	
	// Split the string into parts (before and after decimal point)
	parts := strings.Split(amountStr, ".")
	integerPart := parts[0]
	decimalPart := ""
	if len(parts) > 1 {
		decimalPart = parts[1]
	}
	
	// Format the integer part with dot as thousand separator
	formattedInteger := ""
	for i, c := range integerPart {
		if i > 0 && (len(integerPart)-i)%3 == 0 {
			formattedInteger += "."
		}
		formattedInteger += string(c)
	}
	
	// If there's a decimal part and it's not "00", add it back with comma
	if decimalPart != "" && decimalPart != "00" {
		// Trim trailing zeros
		decimalPart = strings.TrimRight(decimalPart, "0")
		return fmt.Sprintf("%s,%s", formattedInteger, decimalPart)
	}
	
	return formattedInteger
}
