package task

import "context"

type Service interface {
	Create(context.Context, CreateInput) (Model, error)
	List(context.Context) ([]Model, error)
	Get(context.Context, uint) (Model, error)
	Update(context.Context, uint, UpdateInput) (Model, error)
	Delete(context.Context, uint) error
}

type service struct{ repo Repository }

type CreateInput struct {
	Title string
}

type UpdateInput struct {
	Title     string
	Completed bool
}

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) Create(ctx context.Context, input CreateInput) (Model, error) {
	model := Model{Title: input.Title}
	return model, s.repo.Create(ctx, &model)
}

func (s *service) List(ctx context.Context) ([]Model, error) { return s.repo.List(ctx) }

func (s *service) Get(ctx context.Context, id uint) (Model, error) { return s.repo.Get(ctx, id) }

func (s *service) Update(ctx context.Context, id uint, input UpdateInput) (Model, error) {
	return s.repo.Update(ctx, id, &Model{Title: input.Title, Completed: input.Completed})
}

func (s *service) Delete(ctx context.Context, id uint) error { return s.repo.Delete(ctx, id) }
