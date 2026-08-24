package task

import "context"

type Service interface {
	Create(context.Context, CreateInput) (Model, error)
	List(context.Context, string) ([]Model, error)
	Get(context.Context, string, string) (Model, error)
	Update(context.Context, string, string, UpdateInput) (Model, error)
	Delete(context.Context, string, string) error
}

type service struct{ repo Repository }

type CreateInput struct {
	CustomerID string
	Title      string
}

type UpdateInput struct {
	Title     string
	Completed bool
}

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) Create(ctx context.Context, input CreateInput) (Model, error) {
	model := Model{Title: input.Title}
	model.CustomerID = input.CustomerID
	return model, s.repo.Create(ctx, &model)
}

func (s *service) List(ctx context.Context, customerID string) ([]Model, error) {
	return s.repo.List(ctx, customerID)
}

func (s *service) Get(ctx context.Context, customerID, id string) (Model, error) {
	return s.repo.Get(ctx, customerID, id)
}

func (s *service) Update(ctx context.Context, customerID, id string, input UpdateInput) (Model, error) {
	return s.repo.Update(ctx, customerID, id, &Model{Title: input.Title, Completed: input.Completed})
}

func (s *service) Delete(ctx context.Context, customerID, id string) error {
	return s.repo.Delete(ctx, customerID, id)
}
