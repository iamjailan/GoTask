package database

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var migrationFilePattern = regexp.MustCompile(`^(\d+)_(.+)\.up\.sql$`)

type migrationDefinition struct {
	version uint
	name    string
}

type migrationHistory struct {
	db    *gorm.DB
	sqlDB *sql.DB
}

func Open(databaseURL string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
}

func Migrate(databaseURL, migrationsPath string) error {
	migrator, err := newMigrator(databaseURL, migrationsPath)
	if err != nil {
		return err
	}
	defer migrator.Close()

	history, _, err := newMigrationHistory(databaseURL, migrationsPath)
	if err != nil {
		return err
	}
	defer history.Close()

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}

	version, dirty, err := migrationVersion(migrator)
	if err != nil {
		return err
	}
	if dirty {
		return fmt.Errorf("apply migrations: version %d is dirty", version)
	}
	if err := history.markAppliedThrough(version); err != nil {
		return err
	}
	return nil
}

func Rollback(databaseURL, migrationsPath string) error {
	migrator, err := newMigrator(databaseURL, migrationsPath)
	if err != nil {
		return err
	}
	defer migrator.Close()

	history, _, err := newMigrationHistory(databaseURL, migrationsPath)
	if err != nil {
		return err
	}
	defer history.Close()

	version, dirty, err := migrationVersion(migrator)
	if err != nil {
		return err
	}
	if dirty {
		return fmt.Errorf("rollback migration: version %d is dirty", version)
	}
	if err := history.markAppliedThrough(version); err != nil {
		return err
	}

	if err := migrator.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("rollback migration: %w", err)
	}
	if err := history.markRolledBack(version); err != nil {
		return err
	}
	return nil
}

// Reset removes every object in PostgreSQL's public schema.
// This deletes all tables, data, indexes, sequences, and migration records in that schema.
func Reset(databaseURL string) error {
	db, err := Open(databaseURL)
	if err != nil {
		return fmt.Errorf("open database for reset: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get database connection for reset: %w", err)
	}
	defer sqlDB.Close()

	if err := db.Exec("DROP SCHEMA IF EXISTS public CASCADE").Error; err != nil {
		return fmt.Errorf("reset database: drop public schema: %w", err)
	}
	if err := db.Exec("CREATE SCHEMA public").Error; err != nil {
		return fmt.Errorf("reset database: create public schema: %w", err)
	}
	return nil
}

func Version(databaseURL, migrationsPath string) (uint, bool, error) {
	migrator, err := newMigrator(databaseURL, migrationsPath)
	if err != nil {
		return 0, false, err
	}
	defer migrator.Close()

	return migrationVersion(migrator)
}

func newMigrationHistory(databaseURL, migrationsPath string) (*migrationHistory, []migrationDefinition, error) {
	definitions, err := migrationDefinitions(migrationsPath)
	if err != nil {
		return nil, nil, err
	}

	db, err := Open(databaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("open database for migration history: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("get database connection for migration history: %w", err)
	}

	history := &migrationHistory{db: db, sqlDB: sqlDB}
	if err := history.ensureTable(); err != nil {
		sqlDB.Close()
		return nil, nil, err
	}
	if err := history.register(definitions); err != nil {
		sqlDB.Close()
		return nil, nil, err
	}
	return history, definitions, nil
}

func (h *migrationHistory) Close() {
	_ = h.sqlDB.Close()
}

func (h *migrationHistory) ensureTable() error {
	return h.db.Exec(`
		CREATE TABLE IF NOT EXISTS migration_history (
			version BIGINT PRIMARY KEY,
			name TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			applied_at TIMESTAMPTZ,
			is_applied BOOLEAN NOT NULL DEFAULT FALSE
		)
	`).Error
}

func (h *migrationHistory) register(definitions []migrationDefinition) error {
	for _, definition := range definitions {
		if err := h.db.Exec(`
			INSERT INTO migration_history (version, name)
			VALUES (?, ?)
			ON CONFLICT (version) DO UPDATE SET name = EXCLUDED.name
		`, definition.version, definition.name).Error; err != nil {
			return fmt.Errorf("register migration %d: %w", definition.version, err)
		}
	}
	return nil
}

func (h *migrationHistory) markAppliedThrough(version uint) error {
	if version == 0 {
		return nil
	}
	if err := h.db.Exec(`
		UPDATE migration_history
		SET is_applied = TRUE,
			applied_at = COALESCE(applied_at, CURRENT_TIMESTAMP)
		WHERE version <= ?
	`, version).Error; err != nil {
		return fmt.Errorf("mark migrations through %d as applied: %w", version, err)
	}
	return nil
}

func (h *migrationHistory) markRolledBack(version uint) error {
	if version == 0 {
		return nil
	}
	if err := h.db.Exec(`
		UPDATE migration_history
		SET is_applied = FALSE
		WHERE version = ?
	`, version).Error; err != nil {
		return fmt.Errorf("mark migration %d as rolled back: %w", version, err)
	}
	return nil
}

func migrationDefinitions(migrationsPath string) ([]migrationDefinition, error) {
	files, err := filepath.Glob(filepath.Join(migrationsPath, "*.up.sql"))
	if err != nil {
		return nil, fmt.Errorf("list migration files: %w", err)
	}

	definitions := make([]migrationDefinition, 0, len(files))
	for _, file := range files {
		matches := migrationFilePattern.FindStringSubmatch(filepath.Base(file))
		if matches == nil {
			continue
		}
		version, err := strconv.ParseUint(matches[1], 10, 0)
		if err != nil {
			return nil, fmt.Errorf("parse migration version in %q: %w", file, err)
		}
		definitions = append(definitions, migrationDefinition{
			version: uint(version),
			name:    matches[2],
		})
	}
	sort.Slice(definitions, func(i, j int) bool { return definitions[i].version < definitions[j].version })
	return definitions, nil
}

func migrationVersion(migrator *migrate.Migrate) (uint, bool, error) {
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
