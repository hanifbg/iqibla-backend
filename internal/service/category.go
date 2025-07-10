package service

import "github.com/hanifbg/landing_backend/internal/model/entity"

type CategoryService interface {
	GetCategoryBySlug(slug string) (*entity.Category, error)
}
