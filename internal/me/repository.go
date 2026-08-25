package me

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gotask/internal/auth/models"
)

var (
	ErrNotFound   = errors.New("user not found")
	ErrEmailExist = errors.New("email already exists")
)

type Repository interface {
	Get(context.Context, string) (models.Model, error)
	Update(context.Context, string, UpdateInput) (models.Model, error)
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

func (r *repository) Update(ctx context.Context, id string, input UpdateInput) (models.Model, error) {
	model, err := r.Get(ctx, id)
	if err != nil {
		return models.Model{}, err
	}

	if input.Email != "" {
		var count int64
		if err := r.db.WithContext(ctx).Model(&models.Model{}).
			Where("email = ? AND id <> ?", input.Email, id).Count(&count).Error; err != nil {
			return models.Model{}, err
		}
		if count > 0 {
			return models.Model{}, ErrEmailExist
		}
		model.Email = input.Email
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
