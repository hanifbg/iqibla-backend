package entity

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a product category in the database.
type Category struct {
	ID                  string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name                string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Slug                string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"slug"`
	HeroHeadline        *string        `gorm:"type:varchar(255)" json:"hero_headline,omitempty"`
	HeroSubheadline     *string        `gorm:"type:text" json:"hero_subheadline,omitempty"`
	HeroImageUrl        *string        `gorm:"type:varchar(255)" json:"hero_image_url,omitempty"`
	Section3Headline    *string        `gorm:"type:varchar(255)" json:"section3_headline,omitempty"`
	Section3Subheadline *string        `gorm:"type:text" json:"section3_subheadline,omitempty"`
	Section3ImageUrl    *string        `gorm:"type:varchar(255)" json:"section3_image_url,omitempty"`
	MetaTitle           *string        `gorm:"type:varchar(255)" json:"meta_title,omitempty"`
	MetaDescription     *string        `gorm:"type:text" json:"meta_description,omitempty"`
	DisplayOrder        int            `gorm:"default:0" json:"display_order"`
	IsActive            bool           `gorm:"default:true" json:"is_active"`
	CreatedAt           time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Products            []Product      `gorm:"-" json:"products,omitempty"`
}
