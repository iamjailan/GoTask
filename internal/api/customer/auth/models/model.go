package models

import "time"

type Model struct {
	ID              string     `gorm:"primaryKey;type:varchar(50);default:concat('cus_', gen_random_uuid())"`
	FirstName       string     `gorm:"type:varchar(100);not null"`
	LastName        string     `gorm:"type:varchar(100);not null"`
	Email           string     `gorm:"type:varchar(255);not null;uniqueIndex"`
	PasswordHash    string     `gorm:"type:text;not null"`
	Phone           string     `gorm:"type:varchar(30)"`
	AvatarURL       string     `gorm:"type:text"`
	Role            string     `gorm:"type:varchar(50);not null;default:user"`
	IsActive        bool       `gorm:"not null;default:true"`
	EmailVerifiedAt *time.Time `gorm:"type:timestamp"`
	LastLoginAt     *time.Time `gorm:"type:timestamp"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (Model) TableName() string { return "customers" }
