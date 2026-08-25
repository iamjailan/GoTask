package me

import (
	"context"
	"strings"

	"gotask/internal/auth/models"
	"gotask/internal/auth/utils"
)

type Service interface {
	Get(context.Context, string) (models.Model, error)
	Update(context.Context, string, UpdateInput) (models.Model, error)
	Delete(context.Context, string) error
}

type service struct{ repo Repository }

type UpdateInput struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
	AvatarURL string
}

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) Get(ctx context.Context, id string) (models.Model, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, id string, input UpdateInput) (models.Model, error) {
	input.FirstName = strings.TrimSpace(input.FirstName)
	input.LastName = strings.TrimSpace(input.LastName)
	input.Email = utils.NormalizeEmail(input.Email)
	input.Phone = strings.TrimSpace(input.Phone)
	return s.repo.Update(ctx, id, input)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
