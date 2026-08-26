package me

import (
	"context"
	"errors"
	"testing"

	"gotask/internal/auth/models"
	"gotask/internal/auth/utils"
	gotaskemail "gotask/internal/email"
	metypes "gotask/internal/types/me"
)

func TestChangeEmailUpdatesAddressAndNotifiesPreviousAddress(t *testing.T) {
	passwordHash, err := utils.HashPassword("current-password")
	if err != nil {
		t.Fatal(err)
	}
	repo := &fakeRepository{user: models.Model{
		ID: "cus_123", FirstName: "Ada", Email: "old@example.com", PasswordHash: passwordHash,
	}}
	email := &fakeEmailService{}
	service := NewService(repo, email)

	user, err := service.ChangeEmail(context.Background(), "cus_123", metypes.ChangeEmailInput{
		Email: " NEW@Example.com ", CurrentPassword: "current-password",
	})
	if err != nil {
		t.Fatalf("ChangeEmail() error = %v", err)
	}
	if user.Email != "new@example.com" || repo.user.Email != "new@example.com" {
		t.Fatalf("email = %q, want %q", repo.user.Email, "new@example.com")
	}
	if email.recipient != "old@example.com" || email.newEmail != "new@example.com" {
		t.Fatalf("notification = recipient %q, new email %q", email.recipient, email.newEmail)
	}
}

func TestChangeEmailRestoresAddressWhenNotificationFails(t *testing.T) {
	passwordHash, err := utils.HashPassword("current-password")
	if err != nil {
		t.Fatal(err)
	}
	repo := &fakeRepository{user: models.Model{
		ID: "cus_123", FirstName: "Ada", Email: "old@example.com", PasswordHash: passwordHash,
	}}
	service := NewService(repo, &fakeEmailService{sendErr: errors.New("provider unavailable")})

	_, err = service.ChangeEmail(context.Background(), "cus_123", metypes.ChangeEmailInput{
		Email: "new@example.com", CurrentPassword: "current-password",
	})
	if !errors.Is(err, ErrEmailNotification) {
		t.Fatalf("ChangeEmail() error = %v, want ErrEmailNotification", err)
	}
	if repo.user.Email != "old@example.com" {
		t.Fatalf("email = %q, want original address restored", repo.user.Email)
	}
}

func TestChangePasswordRequiresCurrentPasswordAndStoresNewHash(t *testing.T) {
	passwordHash, err := utils.HashPassword("current-password")
	if err != nil {
		t.Fatal(err)
	}
	repo := &fakeRepository{user: models.Model{ID: "cus_123", PasswordHash: passwordHash}}
	service := NewService(repo, &fakeEmailService{})

	err = service.ChangePassword(context.Background(), "cus_123", metypes.ChangePasswordInput{
		CurrentPassword: "wrong-password", NewPassword: "new-password",
	})
	if !errors.Is(err, ErrInvalidCurrentPassword) {
		t.Fatalf("ChangePassword() error = %v, want ErrInvalidCurrentPassword", err)
	}

	err = service.ChangePassword(context.Background(), "cus_123", metypes.ChangePasswordInput{
		CurrentPassword: "current-password", NewPassword: "new-password",
	})
	if err != nil {
		t.Fatalf("ChangePassword() error = %v", err)
	}
	if !utils.CheckPassword(repo.user.PasswordHash, "new-password") {
		t.Fatal("stored password hash does not match new password")
	}
}

type fakeRepository struct {
	user models.Model
}

func (r *fakeRepository) Get(_ context.Context, id string) (models.Model, error) {
	if r.user.ID != id {
		return models.Model{}, ErrNotFound
	}
	return r.user, nil
}

func (r *fakeRepository) UpdateProfile(_ context.Context, _ string, _ metypes.ProfileUpdateInput) (models.Model, error) {
	return r.user, nil
}

func (r *fakeRepository) EmailExists(_ context.Context, email string) (bool, error) {
	return r.user.Email == email, nil
}

func (r *fakeRepository) UpdateEmail(_ context.Context, id, email string) (models.Model, error) {
	if r.user.ID != id {
		return models.Model{}, ErrNotFound
	}
	r.user.Email = email
	return r.user, nil
}

func (r *fakeRepository) UpdatePassword(_ context.Context, id, passwordHash string) error {
	if r.user.ID != id {
		return ErrNotFound
	}
	r.user.PasswordHash = passwordHash
	return nil
}

func (r *fakeRepository) Delete(_ context.Context, _ string) error { return nil }

type fakeEmailService struct {
	recipient string
	newEmail  string
	sendErr   error
}

func (s *fakeEmailService) SendVerificationCode(context.Context, string, string, string) error {
	return nil
}

func (s *fakeEmailService) SendEmailChangedNotification(_ context.Context, recipient, _ string, newEmail string) error {
	s.recipient = recipient
	s.newEmail = newEmail
	return s.sendErr
}

var _ gotaskemail.Service = (*fakeEmailService)(nil)
