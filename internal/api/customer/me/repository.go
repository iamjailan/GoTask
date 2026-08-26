package me

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gotask/internal/api/customer/auth/models"
	metypes "gotask/internal/types/me"
)

var (
	ErrNotFound                     = errors.New("user not found")
	ErrEmailExist                   = errors.New("email already exists")
	ErrEmailUnchanged               = errors.New("new email matches current email")
	ErrPasswordUnchanged            = errors.New("new password matches current password")
	ErrInvalidCurrentPassword       = errors.New("invalid current password")
	ErrEmailNotification            = errors.New("email change notification failed")
	ErrEmailNotificationUnavailable = errors.New("email change notification is not configured")
)

type Repository interface {
	Get(context.Context, string) (models.Model, error)
	UpdateProfile(context.Context, string, metypes.ProfileUpdateInput) (models.Model, error)
	EmailExists(context.Context, string) (bool, error)
	UpdateEmail(context.Context, string, string) (models.Model, error)
	UpdatePassword(context.Context, string, string) error
	Delete(context.Context, string) error
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) Get(ctx context.Context, id string) (models.Model, error) {
	var model models.Model
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Model{}, ErrNotFound
		}
		return models.Model{}, err
	}
	return model, nil
}

func (r *repository) UpdateProfile(ctx context.Context, id string, input metypes.ProfileUpdateInput) (models.Model, error) {
	model, err := r.Get(ctx, id)
	if err != nil {
		return models.Model{}, err
	}

	if input.FirstName != "" {
		model.FirstName = input.FirstName
	}
	if input.LastName != "" {
		model.LastName = input.LastName
	}
	if input.Phone != "" {
		model.Phone = input.Phone
	}
	if input.AvatarURL != "" {
		model.AvatarURL = input.AvatarURL
	}

	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return models.Model{}, err
	}
	return model, nil
}

func (r *repository) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Model{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repository) UpdateEmail(ctx context.Context, id, email string) (models.Model, error) {
	model, err := r.Get(ctx, id)
	if err != nil {
		return models.Model{}, err
	}
	model.Email = email
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return models.Model{}, err
	}
	return model, nil
}

func (r *repository) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	result := r.db.WithContext(ctx).Model(&models.Model{}).Where("id = ?", id).Update("password_hash", passwordHash)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Model{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
