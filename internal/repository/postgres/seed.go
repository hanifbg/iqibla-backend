package postgres

import (
	"github.com/google/uuid"
	"github.com/hanifbg/landing_backend/internal/model/entity"
	"gorm.io/gorm"
)

// SeedTestData inserts test data into the database
func (repo *RepoDatabase) SeedTestData() error {
	// Check if test data already exists
	var count int64
	if err := repo.DB.Model(&entity.Product{}).Count(&count).Error; err != nil {
		return err
	}

	// Skip seeding if data already exists
	if count > 0 {
		return nil
	}

	// Create test product
	product := entity.Product{
		ID:          uuid.New().String(),
		Name:        "Test Prayer Robe",
		Description: "A beautiful and comfortable prayer robe for daily prayers",
		Category:    "Prayer Robes",
		Brand:       "iQibla",
		ImageURLs:   []string{"https://example.com/robe1.jpg", "https://example.com/robe2.jpg"},
		IsActive:    true,
	}

	// Create test variants
	variants := []entity.ProductVariant{
		{
			ID:            uuid.New().String(),
			ProductID:     product.ID,
			SKU:           "PR-001-S",
			Name:          "Small Black Robe",
			Price:         99.99,
			StockQuantity: 50,
			ImageURL:      "https://example.com/robe1-small.jpg",
			Weight:        0.5,
			Dimensions: &entity.Dimensions{
				Length: 100,
				Width:  60,
				Height: 2,
				Unit:   "cm",
			},
			AttributeValues: map[string]interface{}{
				"size":  "S",
				"color": "Black",
			},
			IsActive: true,
		},
		{
			ID:            uuid.New().String(),
			ProductID:     product.ID,
			SKU:           "PR-001-M",
			Name:          "Medium Black Robe",
			Price:         99.99,
			StockQuantity: 40,
			ImageURL:      "https://example.com/robe1-medium.jpg",
			Weight:        0.6,
			Dimensions: &entity.Dimensions{
				Length: 110,
				Width:  65,
				Height: 2,
				Unit:   "cm",
			},
			AttributeValues: map[string]interface{}{
				"size":  "M",
				"color": "Black",
			},
			IsActive: true,
		},
	}

	// Use transaction to ensure data consistency
	return repo.DB.Transaction(func(tx *gorm.DB) error {
		// Insert product
		if err := tx.Create(&product).Error; err != nil {
			return err
		}

		// Insert variants
		for _, variant := range variants {
			if err := tx.Create(&variant).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
