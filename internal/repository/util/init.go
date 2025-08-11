package util

import (
	"log"
	"time"

	"github.com/hanifbg/landing_backend/config"
	"github.com/hanifbg/landing_backend/internal/repository"
	"github.com/hanifbg/landing_backend/internal/repository/mail"
	db "github.com/hanifbg/landing_backend/internal/repository/postgres"
	"github.com/hanifbg/landing_backend/internal/repository/rajaongkir"
)

type RepoWrapper struct {
	ProductRepo  repository.ProductRepository
	CartRepo     repository.CartRepository
	PaymentRepo  repository.PaymentRepository
	CategoryRepo repository.CategoryRepository
	ShippingRepo repository.ShippingRepository
	MailRepo     repository.Mailer
}

func New(cfg *config.AppConfig) (repoWrapper *RepoWrapper, err error) {
	// Initialize database connection
	var dbConnection *db.RepoDatabase
	dbConnection, err = db.Init(cfg)
	if err != nil {
		log.Printf("Warning: Database initialization failed: %v", err)
		return nil, err
	}

	mailer, err := mail.Init(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize RajaOngkir repository with caching configuration
	rajaOngkirRepo := rajaongkir.NewRepository(rajaongkir.Config{
		APIKey:  cfg.RajaOngkirAPIKey,
		BaseURL: cfg.RajaOngkirBaseURL,
		Timeout: 30 * time.Second,
		// Cache configuration
		CacheEnabled:      cfg.RajaOngkirCacheEnabled,
		CacheTTLHours:     cfg.RajaOngkirCacheTTLHours,
		WarmupOnStartup:   cfg.RajaOngkirWarmupOnStartup,
		WarmupTimeoutSecs: cfg.RajaOngkirWarmupTimeoutSecs,
	})

	repoWrapper = &RepoWrapper{
		ProductRepo:  dbConnection,
		CartRepo:     dbConnection,
		PaymentRepo:  dbConnection,
		CategoryRepo: dbConnection,
		ShippingRepo: rajaOngkirRepo,
		MailRepo:     mailer,
	}

	return repoWrapper, nil
}
