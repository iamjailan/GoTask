package auth

import "time"

type RegisterInput struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Phone     string
	AvatarURL string
}

type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required,max=100"`
	LastName  string `json:"last_name" binding:"required,max=100"`
	Email     string `json:"email" binding:"required,email,max=255"`
	Password  string `json:"password" binding:"required,min=8,max=72"`
	Phone     string `json:"phone" binding:"omitempty,max=30"`
	AvatarURL string `json:"avatar_url" binding:"omitempty,url"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ConfirmEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6,numeric"`
}

type VerificationResponse struct {
	Message   string           `json:"message"`
	Customer  CustomerResponse `json:"customer"`
	Token     string           `json:"token"`
	ExpiresAt time.Time        `json:"expires_at"`
}

type RegistrationResponse struct {
	Email                string `json:"email"`
	VerificationRequired bool   `json:"verification_required"`
	Message              string `json:"message"`
}

type AuthResponse struct {
	Token     string           `json:"token"`
	ExpiresAt time.Time        `json:"expires_at"`
	Customer  CustomerResponse `json:"customer"`
}

type CustomerResponse struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
