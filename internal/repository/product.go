package repository

//go:generate mockgen -source=product.go -destination=../service/product/mocks/product_repository_mock.go -package=mocks

import "github.com/hanifbg/landing_backend/internal/model/entity"

type ProductRepository interface {
	GetAllProducts() ([]entity.Product, error)
	GetProductByID(id string) (*entity.Product, error)
	GetAllProductsByCategory(category string) ([]entity.Product, error)
	GetAllProductsByCategorySlug(categorySlug string) ([]entity.Product, error)
}
