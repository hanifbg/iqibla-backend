package util

import (
	"log"
	"net/http"
	"time"

	"github.com/hanifbg/landing_backend/config"
	"github.com/hanifbg/landing_backend/internal/repository"
	"github.com/hanifbg/landing_backend/internal/repository/external"
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
	WhatsAppRepo repository.WhatsApp
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

	httpClient := &http.Client{
		Timeout: time.Duration(cfg.HttpTimeout) * time.Second,
	}

	// Initialize RajaOngkir repository with caching configuration
	rajaOngkirRepo := rajaongkir.NewRepository(rajaongkir.Config{
		APIKey:  cfg.RajaOngkirAPIKey,
		BaseURL: cfg.RajaOngkirBaseURL,
		Client:  httpClient,
		// Cache configuration
		CacheEnabled:      cfg.RajaOngkirCacheEnabled,
		CacheTTLHours:     cfg.RajaOngkirCacheTTLHours,
		WarmupOnStartup:   cfg.RajaOngkirWarmupOnStartup,
		WarmupTimeoutSecs: cfg.RajaOngkirWarmupTimeoutSecs,
	})

	externalRepo := external.New(cfg, httpClient)

	repoWrapper = &RepoWrapper{
		ProductRepo:  dbConnection,
		CartRepo:     dbConnection,
		PaymentRepo:  dbConnection,
		CategoryRepo: dbConnection,
		ShippingRepo: rajaOngkirRepo,
		MailRepo:     mailer,
		WhatsAppRepo: externalRepo.WAApi,
	}

	return repoWrapper, nil
}
