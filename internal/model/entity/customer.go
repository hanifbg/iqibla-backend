package entity

import (
	"time"

	"gorm.io/gorm"
)

// Customer represents a customer in the system
type Customer struct {
	ID                string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	FirstName         string         `gorm:"type:varchar(255)" json:"first_name"`
	LastName          string         `gorm:"type:varchar(255)" json:"last_name"`
	Email             string         `gorm:"type:varchar(255);uniqueIndex" json:"email"`
	PhoneNumber       string         `gorm:"type:varchar(50)" json:"phone_number"`
	PasswordHash      string         `gorm:"type:varchar(255)" json:"password_hash"`
	Salt              string         `gorm:"type:varchar(255)" json:"salt"`
	AddressDetails    JSONMap        `gorm:"type:jsonb" json:"address_details,omitempty"`
	LastLoginAt       *time.Time     `json:"last_login_at,omitempty"`
	IsEmailVerified   bool           `gorm:"default:false" json:"is_email_verified"`
	CreatedAt         time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}