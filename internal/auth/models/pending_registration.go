package models

import "time"

type PendingRegistration struct {
	ID            string    `gorm:"primaryKey;type:varchar(50);default:concat('reg_', gen_random_uuid())"`
	FirstName     string    `gorm:"type:varchar(100);not null"`
	LastName      string    `gorm:"type:varchar(100);not null"`
	Email         string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	PasswordHash  string    `gorm:"type:text;not null"`
	Phone         string    `gorm:"type:varchar(30)"`
	AvatarURL     string    `gorm:"type:text"`
	CodeHash      string    `gorm:"type:text;not null"`
	CodeExpiresAt time.Time `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (PendingRegistration) TableName() string { return "pending_registrations" }
