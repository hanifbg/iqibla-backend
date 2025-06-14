package repository

import "github.com/hanifbg/landing_backend/internal/model/entity"

type ProductRepository interface {
	GetAllProducts() ([]entity.Product, error)
	GetProductByID(id string) (*entity.Product, error)
}
