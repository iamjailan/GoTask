package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gotask/internal/auth/models"
	"gotask/internal/task"
)

func Open(databaseURL string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
}

func Migrate(db *gorm.DB) error {
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto").Error; err != nil {
		return err
	}
	return db.AutoMigrate(&task.Model{}, &models.Model{})
}
