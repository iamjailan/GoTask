package task

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("task not found")

type Repository interface {
	Create(context.Context, *Model) error
	List(context.Context, string) ([]Model, error)
	Get(context.Context, string, string) (Model, error)
	Update(context.Context, string, string, *Model) (UpdateResult, error)
	Delete(context.Context, string, string) error
	Transaction(context.Context, func(Repository, StatisticsRepository) error) error
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) Create(ctx context.Context, model *Model) error {
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *repository) Transaction(ctx context.Context, fn func(Repository, StatisticsRepository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&repository{db: tx}, &statisticsRepository{db: tx})
	})
}

func (r *repository) List(ctx context.Context, customerID string) ([]Model, error) {
	var models []Model
	err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).Order("id asc").Find(&models).Error
	return models, err
}

func (r *repository) Get(ctx context.Context, customerID, id string) (Model, error) {
	var model Model
	if err := r.db.WithContext(ctx).Where("customer_id = ? AND id = ?", customerID, id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Model{}, ErrNotFound
		}
		return Model{}, err
	}
	return model, nil
}

func (r *repository) Update(ctx context.Context, customerID, id string, input *Model) (UpdateResult, error) {
	model, err := r.Get(ctx, customerID, id)
	if err != nil {
		return UpdateResult{}, err
	}
	previous := model
	model.Title = input.Title
	model.Description = input.Description
	if input.Status != "" {
		model.Status = input.Status
	}
	if input.Priority != "" {
		model.Priority = input.Priority
	}
	model.DueDate = input.DueDate
	model.Completed = input.Completed
	if input.Completed {
		if model.CompletedAt == nil {
			now := time.Now().UTC()
			model.CompletedAt = &now
		}
		model.Status = "completed"
	} else {
		model.CompletedAt = nil
		if model.Status == "completed" {
			model.Status = "pending"
		}
	}
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return UpdateResult{}, err
	}
	return UpdateResult{Previous: previous, Current: model}, nil
}

func (r *repository) Delete(ctx context.Context, customerID, id string) error {
	result := r.db.WithContext(ctx).Where("customer_id = ? AND id = ?", customerID, id).Delete(&Model{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
