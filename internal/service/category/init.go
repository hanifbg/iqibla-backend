package category

import "github.com/hanifbg/landing_backend/internal/repository"

type CategoryService struct {
	categoryRepo repository.CategoryRepository
	productRepo  repository.ProductRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository, productRepo repository.ProductRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		productRepo:  productRepo,
	}
}
