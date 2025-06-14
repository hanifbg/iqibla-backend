package postgres

import (
	"fmt"

	"github.com/hanifbg/landing_backend/internal/model/entity"
	"gorm.io/gorm"
)

func (repo *RepoDatabase) GetAllProducts() ([]entity.Product, error) {
	var products []entity.Product

	result := repo.DB.Preload("Variants", "is_active = ?", true).
		Where("is_active = ?", true).
		Find(&products)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch products: %v", result.Error)
	}

	return products, nil
}

func (repo *RepoDatabase) GetProductByID(id string) (*entity.Product, error) {
	var product entity.Product

	result := repo.DB.Preload("Variants", "is_active = ?", true).
		Where("id = ? AND is_active = ?", id, true).
		First(&product)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch product: %v", result.Error)
	}

	return &product, nil
}
