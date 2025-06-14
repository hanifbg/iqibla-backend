package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Category    string    `gorm:"type:varchar(100)" json:"category"`
	Brand       string    `gorm:"type:varchar(100)" json:"brand"`
	ImageURLs   JSONArray `gorm:"type:jsonb" json:"image_urls"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	Variants    []ProductVariant `gorm:"foreignKey:ProductID" json:"variants"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

type ProductVariant struct {
	ID              string          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProductID       string          `gorm:"type:uuid;not null" json:"product_id"`
	SKU             string          `gorm:"type:varchar(50);uniqueIndex" json:"sku"`
	Name            string          `gorm:"type:varchar(255)" json:"name"`
	Price           float64         `gorm:"type:decimal(10,2);not null" json:"price"`
	StockQuantity   int             `gorm:"not null" json:"stock_quantity"`
	ImageURL        string          `gorm:"type:varchar(255)" json:"image_url"`
	Weight          float64         `gorm:"type:decimal(10,2)" json:"weight"`
	Dimensions      *Dimensions     `gorm:"type:jsonb" json:"dimensions,omitempty"`
	AttributeValues JSONMap         `gorm:"type:jsonb" json:"attribute_values"`
	IsActive        bool            `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       *time.Time      `gorm:"index" json:"deleted_at,omitempty"`
}

type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Unit   string  `json:"unit"`
}

// Scan implements the sql.Scanner interface for Dimensions
func (d *Dimensions) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, d)
}

// Value implements the driver.Valuer interface for Dimensions
func (d Dimensions) Value() (driver.Value, error) {
	if d == (Dimensions{}) {
		return nil, nil
	}
	return json.Marshal(d)
}

// JSONArray is a custom type for handling string arrays in PostgreSQL jsonb
type JSONArray []string

func (j *JSONArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, j)
}

func (j JSONArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// JSONMap is a custom type for handling key-value pairs in PostgreSQL jsonb
type JSONMap map[string]interface{}

func (j *JSONMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, j)
}

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// BeforeCreate will set a UUID rather than numeric ID.
func (product *Product) BeforeCreate(tx *gorm.DB) error {
	if product.ID == "" {
		product.ID = uuid.New().String()
	}
	return nil
}

// BeforeCreate will set a UUID rather than numeric ID.
func (variant *ProductVariant) BeforeCreate(tx *gorm.DB) error {
	if variant.ID == "" {
		variant.ID = uuid.New().String()
	}
	return nil
}