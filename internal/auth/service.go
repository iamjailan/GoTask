package auth

import (
	"context"
	"errors"
	"time"

	"gotask/internal/auth/models"
	"gotask/internal/auth/types"
	"gotask/internal/auth/utils"
)

var (
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInactive           = errors.New("customer is inactive")
)

type Service interface {
	Register(context.Context, types.RegisterInput) (models.Model, string, error)
	Login(context.Context, string, string) (models.Model, string, error)
}

type service struct {
	repo      Repository
	jwtSecret []byte
}

func NewService(repo Repository, jwtSecret string) Service {
	return &service{repo: repo, jwtSecret: []byte(jwtSecret)}
}

func (s *service) Register(ctx context.Context, input types.RegisterInput) (models.Model, string, error) {
	email := utils.NormalizeEmail(input.Email)
	if _, err := s.repo.FindByEmail(ctx, email); err == nil {
		return models.Model{}, "", ErrEmailExists
	} else if !errors.Is(err, ErrNotFound) {
		return models.Model{}, "", err
	}

	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		return models.Model{}, "", err
	}
	model := models.Model{
		FirstName: input.FirstName, LastName: input.LastName, Email: email,
		PasswordHash: hash, Phone: input.Phone, AvatarURL: input.AvatarURL,
		Role: "user", IsActive: true,
	}
	if err := s.repo.Create(ctx, &model); err != nil {
		return models.Model{}, "", err
	}
	token, err := s.token(model)
	return model, token, err
}

func (s *service) Login(ctx context.Context, email, password string) (models.Model, string, error) {
	model, err := s.repo.FindByEmail(ctx, utils.NormalizeEmail(email))
	if err != nil || !utils.CheckPassword(model.PasswordHash, password) {
		return models.Model{}, "", ErrInvalidCredentials
	}
	if !model.IsActive {
		return models.Model{}, "", ErrInactive
	}
	now := time.Now().UTC()
	if err := s.repo.UpdateLastLogin(ctx, model.ID, now); err != nil {
		return models.Model{}, "", err
	}
	model.LastLoginAt = &now
	token, err := s.token(model)
	return model, token, err
}

func (s *service) token(model models.Model) (string, error) {
	return utils.CreateToken(string(s.jwtSecret), model.ID, model.Email, model.Role, time.Now().UTC())
}
