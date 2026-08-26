package database

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open(databaseURL string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
}

func Migrate(databaseURL, migrationsPath string) error {
	migrator, err := newMigrator(databaseURL, migrationsPath)
	if err != nil {
		return err
	}
	defer migrator.Close()

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}

func Rollback(databaseURL, migrationsPath string) error {
	migrator, err := newMigrator(databaseURL, migrationsPath)
	if err != nil {
		return err
	}
	defer migrator.Close()

	if err := migrator.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("rollback migration: %w", err)
	}
	return nil
}

func Version(databaseURL, migrationsPath string) (uint, bool, error) {
	migrator, err := newMigrator(databaseURL, migrationsPath)
	if err != nil {
		return 0, false, err
	}
	defer migrator.Close()

	version, dirty, err := migrator.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, fmt.Errorf("read migration version: %w", err)
	}
	return version, dirty, nil
}

func newMigrator(databaseURL, migrationsPath string) (*migrate.Migrate, error) {
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("resolve migrations path: %w", err)
	}
	migrator, err := migrate.New("file://"+absPath, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create migrator: %w", err)
	}
	return migrator, nil
}
