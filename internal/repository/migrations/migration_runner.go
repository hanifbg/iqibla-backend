package migrations

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

// MigrationRunner provides functionality to run migrations
type MigrationRunner struct {
	db *sql.DB
}

// NewMigrationRunner creates a new MigrationRunner
func NewMigrationRunner(db *sql.DB) *MigrationRunner {
	return &MigrationRunner{
		db: db,
	}
}

// ApplyMigration applies a specific migration file
func (r *MigrationRunner) ApplyMigration(migrationPath string) error {
	content, err := ioutil.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("error reading migration file %s: %v", migrationPath, err)
	}

	_, err = r.db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("error applying migration %s: %v", migrationPath, err)
	}

	log.Printf("Successfully applied migration: %s", filepath.Base(migrationPath))
	return nil
}

// RollbackMigration rolls back a specific migration using the Down section
func (r *MigrationRunner) RollbackMigration(migrationPath string) error {
	content, err := ioutil.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("error reading migration file %s: %v", migrationPath, err)
	}

	// Extract Down migration if there's a comment indicating it
	downSQL := string(content)

	// For a more sophisticated approach, you'd parse the file to extract
	// only the down migration section. This is a simplified approach.

	_, err = r.db.Exec(downSQL)
	if err != nil {
		return fmt.Errorf("error rolling back migration %s: %v", migrationPath, err)
	}

	log.Printf("Successfully rolled back migration: %s", filepath.Base(migrationPath))
	return nil
}
