package entity

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID         string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CustomerID *string        `gorm:"type:uuid;index" json:"customer_id,omitempty"` // Nullable for guest carts
	CartItems  []CartItem     `gorm:"foreignKey:CartID" json:"cart_items,omitempty"`
	CreatedAt  time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"not null" json:"updated_at"`
	ExpiresAt  *time.Time     `json:"expires_at,omitempty"`          // When the cart should expire (e.g., after 24 hours of inactivity)
	IsActive   bool           `gorm:"default:true" json:"is_active"` // False if converted to order or explicitly abandoned
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type CartItem struct {
	ID               string          `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CartID           string          `gorm:"type:uuid;not null" json:"cart_id"`
	ProductVariantID string          `gorm:"type:uuid;not null" json:"product_variant_id"`
	ProductVariant   *ProductVariant `gorm:"foreignKey:ProductVariantID" json:"product_variant,omitempty"`
	Quantity         int             `gorm:"not null" json:"quantity"`
	CreatedAt        time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt        time.Time       `gorm:"not null" json:"updated_at"`
	DeletedAt        gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
}
