package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gotask/internal/task"
)

func Open(databaseURL string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&task.Model{})
}
