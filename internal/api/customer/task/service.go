package task

import (
	"context"
	"time"
)

type Service interface {
	Create(context.Context, CreateInput) (Model, error)
	List(context.Context, string) ([]Model, error)
	Get(context.Context, string, string) (Model, error)
	Update(context.Context, string, string, UpdateInput) (Model, error)
	Delete(context.Context, string, string) error
}

type service struct{ repo Repository }

type CreateInput struct {
	CustomerID  string
	Title       string
	Description string
	Status      string
	Priority    string
	DueDate     *time.Time
	Completed   bool
}

type UpdateInput struct {
	Title       string
	Description string
	Status      string
	Priority    string
	DueDate     *time.Time
	Completed   bool
}

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) Create(ctx context.Context, input CreateInput) (Model, error) {
	status := input.Status
	if status == "" {
		status = "pending"
	}
	priority := input.Priority
	if priority == "" {
		priority = "medium"
	}
	model := Model{
		CustomerID: input.CustomerID, Title: input.Title, Description: input.Description,
		Status: status, Priority: priority, DueDate: input.DueDate, Completed: input.Completed,
	}
	if input.Completed {
		now := time.Now().UTC()
		model.Status = "completed"
		model.CompletedAt = &now
	}
	return model, s.repo.Create(ctx, &model)
}

func (s *service) List(ctx context.Context, customerID string) ([]Model, error) {
	return s.repo.List(ctx, customerID)
}

func (s *service) Get(ctx context.Context, customerID, id string) (Model, error) {
	return s.repo.Get(ctx, customerID, id)
}

func (s *service) Update(ctx context.Context, customerID, id string, input UpdateInput) (Model, error) {
	return s.repo.Update(ctx, customerID, id, &Model{
		Title: input.Title, Description: input.Description, Status: input.Status,
		Priority: input.Priority, DueDate: input.DueDate, Completed: input.Completed,
	})
}

func (s *service) Delete(ctx context.Context, customerID, id string) error {
	return s.repo.Delete(ctx, customerID, id)
}
