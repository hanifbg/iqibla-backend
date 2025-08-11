package util

import (
	"github.com/hanifbg/landing_backend/config"
	"github.com/hanifbg/landing_backend/internal/repository"
	"github.com/hanifbg/landing_backend/internal/repository/mail"
	db "github.com/hanifbg/landing_backend/internal/repository/postgres"
)

type RepoWrapper struct {
	ProductRepo  repository.ProductRepository
	CartRepo     repository.CartRepository
	PaymentRepo  repository.PaymentRepository
	CategoryRepo repository.CategoryRepository
	MailRepo     repository.Mailer
}

func New(cfg *config.AppConfig) (repoWrapper *RepoWrapper, err error) {

	var dbConnection *db.RepoDatabase

	dbConnection, err = db.Init(cfg)
	if err != nil {
		return nil, err
	}

	mailer, err := mail.Init(cfg)
	if err != nil {
		return nil, err
	}

	// httpClient := &http.Client{
	// 	Timeout: 10 * time.Second,
	// }

	repoWrapper = &RepoWrapper{
		ProductRepo:  dbConnection,
		CartRepo:     dbConnection,
		PaymentRepo:  dbConnection,
		CategoryRepo: dbConnection,
		MailRepo:     mailer,
	}

	return
}
