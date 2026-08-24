package task

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("task not found")

type Repository interface {
	Create(context.Context, *Model) error
	List(context.Context) ([]Model, error)
	Get(context.Context, uint) (Model, error)
	Update(context.Context, uint, *Model) (Model, error)
	Delete(context.Context, uint) error
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) Create(ctx context.Context, model *Model) error {
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *repository) List(ctx context.Context) ([]Model, error) {
	var models []Model
	err := r.db.WithContext(ctx).Order("id asc").Find(&models).Error
	return models, err
}

func (r *repository) Get(ctx context.Context, id uint) (Model, error) {
	var model Model
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Model{}, ErrNotFound
		}
		return Model{}, err
	}
	return model, nil
}

func (r *repository) Update(ctx context.Context, id uint, input *Model) (Model, error) {
	model, err := r.Get(ctx, id)
	if err != nil {
		return Model{}, err
	}
	model.Title = input.Title
	model.Completed = input.Completed
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return Model{}, err
	}
	return model, nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&Model{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
