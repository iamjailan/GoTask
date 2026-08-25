package auth

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gotask/internal/auth/models"
)

var ErrNotFound = errors.New("customer not found")

type Repository interface {
	Create(context.Context, *models.Model) error
	CreatePending(context.Context, *models.PendingRegistration) error
	Exists(context.Context, string) (bool, error)
	FindByEmail(context.Context, string) (models.Model, error)
	FindPendingByEmail(context.Context, string) (models.PendingRegistration, error)
	DeletePending(context.Context, string) error
	UpdateLastLogin(context.Context, string, time.Time) error
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) Create(ctx context.Context, model *models.Model) error {
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *repository) CreatePending(ctx context.Context, model *models.PendingRegistration) error {
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *repository) Exists(ctx context.Context, id string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Model{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (models.Model, error) {
	var model models.Model
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Model{}, ErrNotFound
		}
		return models.Model{}, err
	}
	return model, nil
}

func (r *repository) FindPendingByEmail(ctx context.Context, email string) (models.PendingRegistration, error) {
	var model models.PendingRegistration
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.PendingRegistration{}, ErrNotFound
		}
		return models.PendingRegistration{}, err
	}
	return model, nil
}

func (r *repository) DeletePending(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.PendingRegistration{}).Error
}

func (r *repository) UpdateLastLogin(ctx context.Context, id string, at time.Time) error {
	return r.db.WithContext(ctx).Model(&models.Model{}).Where("id = ?", id).Update("last_login_at", at).Error
}
