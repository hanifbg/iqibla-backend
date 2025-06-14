package cart

import (
	"github.com/hanifbg/landing_backend/config"
	"github.com/hanifbg/landing_backend/internal/repository"
	"github.com/hanifbg/landing_backend/internal/repository/util"
)

type CartService struct {
	cartRepo repository.CartRepository
}

func New(cfg *config.AppConfig, repo *util.RepoWrapper) *CartService {
	return &CartService{
		cartRepo: repo.CartRepo,
	}
}
