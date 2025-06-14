package postgres

import (
	"fmt"
	"log"

	"github.com/hanifbg/landing_backend/config"
	"github.com/hanifbg/landing_backend/internal/model/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type RepoDatabase struct {
	DB *gorm.DB
}

func Init(config *config.AppConfig) (*RepoDatabase, error) {
	repo := &RepoDatabase{}
	db, err := getConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	repo.DB = db
	//migration
	if err := repo.MigrateDB(); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	// Seed test data
	if err := repo.SeedTestData(); err != nil {
		log.Printf("Warning: Failed to seed test data: %v", err)
	}
	return repo, nil
}

func getConnection(cfg *config.AppConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.DbHost,
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbName,
		cfg.DbPort,
		cfg.DbSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (repo *RepoDatabase) MigrateDB() error {
	return repo.DB.AutoMigrate(
		&entity.Product{},
		&entity.ProductVariant{},
	)
}
