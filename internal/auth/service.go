package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"gotask/internal/auth/models"
	"gotask/internal/auth/types"
	"gotask/internal/auth/utils"
	gotaskemail "gotask/internal/email"
)

var (
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInactive           = errors.New("customer is inactive")
	ErrEmailNotVerified   = errors.New("email is not verified")
	ErrInvalidOTP         = errors.New("invalid verification code")
	ErrOTPExpired         = errors.New("verification code expired")
)

type Service interface {
	Register(context.Context, types.RegisterInput) (models.Model, string, error)
	Login(context.Context, string, string) (models.Model, string, error)
	ConfirmEmail(context.Context, string, string) (models.Model, string, error)
}

type service struct {
	repo      Repository
	jwtSecret []byte
	email     gotaskemail.Service
}

func NewService(repo Repository, jwtSecret string, emailService gotaskemail.Service) Service {
	return &service{repo: repo, jwtSecret: []byte(jwtSecret), email: emailService}
}

func (s *service) Register(ctx context.Context, input types.RegisterInput) (models.Model, string, error) {
	email := utils.NormalizeEmail(input.Email)
	if _, err := s.repo.FindByEmail(ctx, email); err == nil {
		return models.Model{}, "", ErrEmailExists
	} else if !errors.Is(err, ErrNotFound) {
		return models.Model{}, "", err
	}
	if _, err := s.repo.FindPendingByEmail(ctx, email); err == nil {
		return models.Model{}, "", ErrEmailExists
	} else if !errors.Is(err, ErrNotFound) {
		return models.Model{}, "", err
	}

	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		return models.Model{}, "", err
	}
	code, err := gotaskemail.GenerateVerificationCode()
	if err != nil {
		return models.Model{}, "", err
	}
	codeHash, err := utils.HashPassword(code)
	if err != nil {
		return models.Model{}, "", err
	}
	expiresAt := time.Now().UTC().Add(10 * time.Minute)
	pending := models.PendingRegistration{
		FirstName: input.FirstName, LastName: input.LastName, Email: email,
		PasswordHash: hash, Phone: input.Phone, AvatarURL: input.AvatarURL,
		CodeHash: codeHash, CodeExpiresAt: expiresAt,
	}
	if err := s.repo.CreatePending(ctx, &pending); err != nil {
		return models.Model{}, "", err
	}
	if s.email == nil {
		return models.Model{}, "", gotaskemail.ErrNotConfigured
	}
	if err := s.email.SendVerificationCode(ctx, pending.Email, pending.FirstName, code); err != nil {
		_ = s.repo.DeletePending(ctx, pending.ID)
		return models.Model{}, "", err
	}
	return models.Model{FirstName: pending.FirstName, LastName: pending.LastName, Email: pending.Email,
		Phone: pending.Phone, AvatarURL: pending.AvatarURL}, "", nil
}

func (s *service) Login(ctx context.Context, email, password string) (models.Model, string, error) {
	model, err := s.repo.FindByEmail(ctx, utils.NormalizeEmail(email))
	if err != nil || !utils.CheckPassword(model.PasswordHash, password) {
		return models.Model{}, "", ErrInvalidCredentials
	}
	if model.EmailVerifiedAt == nil {
		return models.Model{}, "", ErrEmailNotVerified
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

func (s *service) ConfirmEmail(ctx context.Context, email, code string) (models.Model, string, error) {
	pending, err := s.repo.FindPendingByEmail(ctx, utils.NormalizeEmail(email))
	if err != nil {
		return models.Model{}, "", ErrInvalidOTP
	}
	if time.Now().UTC().After(pending.CodeExpiresAt) {
		return models.Model{}, "", ErrOTPExpired
	}
	if !utils.CheckPassword(pending.CodeHash, strings.TrimSpace(code)) {
		return models.Model{}, "", ErrInvalidOTP
	}
	if _, err := s.repo.FindByEmail(ctx, pending.Email); err == nil {
		return models.Model{}, "", ErrEmailExists
	} else if !errors.Is(err, ErrNotFound) {
		return models.Model{}, "", err
	}
	model := models.Model{FirstName: pending.FirstName, LastName: pending.LastName, Email: pending.Email,
		PasswordHash: pending.PasswordHash, Phone: pending.Phone, AvatarURL: pending.AvatarURL,
		Role: "user", IsActive: true}
	now := time.Now().UTC()
	model.EmailVerifiedAt = &now
	if err := s.repo.Create(ctx, &model); err != nil {
		return models.Model{}, "", err
	}
	if err := s.repo.DeletePending(ctx, pending.ID); err != nil {
		return models.Model{}, "", err
	}
	token, err := s.token(model)
	return model, token, err
}

func (s *service) token(model models.Model) (string, error) {
	return utils.CreateToken(string(s.jwtSecret), model.ID, model.Email, model.Role, time.Now().UTC())
}
