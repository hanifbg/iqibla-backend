package postgres

import (
	"fmt"

	"github.com/hanifbg/landing_backend/internal/model/entity"
	"gorm.io/gorm"
)

func (repo *RepoDatabase) GetCategoryBySlug(slug string) (*entity.Category, error) {
	var category entity.Category

	result := repo.DB.Where("slug = ? AND is_active = ?", slug, true).First(&category)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch category: %v", result.Error)
	}

	return &category, nil
}
