package entity

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID        string      `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CartItems []CartItem  `gorm:"foreignKey:CartID" json:"cart_items,omitempty"`
	CreatedAt time.Time   `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time   `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type CartItem struct {
	ID              string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CartID          string         `gorm:"type:uuid;not null" json:"cart_id"`
	ProductVariantID string         `gorm:"type:uuid;not null" json:"product_variant_id"`
	ProductVariant   ProductVariant `gorm:"foreignKey:ProductVariantID" json:"product_variant,omitempty"`
	Quantity        int            `gorm:"not null" json:"quantity"`
	CreatedAt       time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}