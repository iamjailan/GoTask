package me

import (
	"context"
	"errors"
	"strings"

	"gotask/internal/api/customer/auth/models"
	"gotask/internal/api/customer/auth/utils"
	gotaskemail "gotask/internal/email"
	metypes "gotask/internal/types/me"
)

type Service interface {
	Get(context.Context, string) (models.Model, error)
	UpdateProfile(context.Context, string, metypes.ProfileUpdateInput) (models.Model, error)
	ChangeEmail(context.Context, string, metypes.ChangeEmailInput) (models.Model, error)
	ChangePassword(context.Context, string, metypes.ChangePasswordInput) error
	Delete(context.Context, string) error
}

type service struct {
	repo  Repository
	email gotaskemail.Service
}

func NewService(repo Repository, emailService gotaskemail.Service) Service {
	return &service{repo: repo, email: emailService}
}

func (s *service) Get(ctx context.Context, id string) (models.Model, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) UpdateProfile(ctx context.Context, id string, input metypes.ProfileUpdateInput) (models.Model, error) {
	input.FirstName = strings.TrimSpace(input.FirstName)
	input.LastName = strings.TrimSpace(input.LastName)
	input.Phone = strings.TrimSpace(input.Phone)
	input.AvatarURL = strings.TrimSpace(input.AvatarURL)
	return s.repo.UpdateProfile(ctx, id, input)
}

func (s *service) ChangeEmail(ctx context.Context, id string, input metypes.ChangeEmailInput) (models.Model, error) {
	model, err := s.repo.Get(ctx, id)
	if err != nil {
		return models.Model{}, err
	}
	if !utils.CheckPassword(model.PasswordHash, input.CurrentPassword) {
		return models.Model{}, ErrInvalidCurrentPassword
	}

	email := utils.NormalizeEmail(input.Email)
	if email == model.Email {
		return models.Model{}, ErrEmailUnchanged
	}
	exists, err := s.repo.EmailExists(ctx, email)
	if err != nil {
		return models.Model{}, err
	}
	if exists {
		return models.Model{}, ErrEmailExist
	}
	if s.email == nil {
		return models.Model{}, ErrEmailNotificationUnavailable
	}

	updated, err := s.repo.UpdateEmail(ctx, id, email)
	if err != nil {
		return models.Model{}, err
	}
	if err := s.email.SendEmailChangedNotification(ctx, model.Email, model.FirstName, updated.Email); err != nil {
		if _, restoreErr := s.repo.UpdateEmail(ctx, id, model.Email); restoreErr != nil {
			return models.Model{}, errors.Join(mapEmailNotificationError(err), restoreErr)
		}
		return models.Model{}, mapEmailNotificationError(err)
	}
	return updated, nil
}

func (s *service) ChangePassword(ctx context.Context, id string, input metypes.ChangePasswordInput) error {
	model, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if !utils.CheckPassword(model.PasswordHash, input.CurrentPassword) {
		return ErrInvalidCurrentPassword
	}
	if utils.CheckPassword(model.PasswordHash, input.NewPassword) {
		return ErrPasswordUnchanged
	}
	hash, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(ctx, id, hash)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func mapEmailNotificationError(err error) error {
	if errors.Is(err, gotaskemail.ErrNotConfigured) {
		return ErrEmailNotificationUnavailable
	}
	return ErrEmailNotification
}
