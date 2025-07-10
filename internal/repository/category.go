package repository

import "github.com/hanifbg/landing_backend/internal/model/entity"

type CategoryRepository interface {
	GetCategoryBySlug(slug string) (*entity.Category, error)
}
