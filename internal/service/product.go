package service

import "github.com/hanifbg/landing_backend/internal/model/entity"

type ProductService interface {
	GetAllProducts(category string) ([]entity.Product, error)
	GetProductByID(id string) (*entity.Product, error)
}
