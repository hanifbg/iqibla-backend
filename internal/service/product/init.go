package product

import (
	"github.com/hanifbg/landing_backend/config"
	"github.com/hanifbg/landing_backend/internal/repository"
	"github.com/hanifbg/landing_backend/internal/repository/util"
)

type ProductService struct {
	productRepo repository.ProductRepository
}

func New(cfg *config.AppConfig, repo *util.RepoWrapper) *ProductService {
	return &ProductService{
		productRepo: repo.ProductRepo,
	}
}
