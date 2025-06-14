package entity

import (
	"time"

	"gorm.io/gorm"
)

type Discount struct {
	ID                string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Code              string         `gorm:"uniqueIndex;not null" json:"code"`
	Type              string         `gorm:"not null" json:"type"` // percentage, fixed_amount
	Value             float64        `gorm:"not null" json:"value"`
	MinimumOrderAmount float64        `gorm:"not null;default:0" json:"minimum_order_amount"`
	StartsAt          time.Time      `gorm:"not null" json:"starts_at"`
	ExpiresAt         time.Time      `json:"expires_at"`
	UsageLimit        int            `gorm:"not null;default:0" json:"usage_limit"`
	UsesCount         int            `gorm:"not null;default:0" json:"uses_count"`
	IsActive          bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt         time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}