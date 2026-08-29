package me

import "time"

type ProfileUpdateRequest struct {
	FirstName          string  `json:"first_name" binding:"omitempty,max=100"`
	LastName           string  `json:"last_name" binding:"omitempty,max=100"`
	Phone              string  `json:"phone" binding:"omitempty,max=30"`
	AvatarURL          string  `json:"avatar_url" binding:"omitempty,url"`
	DeprecatedEmail    *string `json:"email" binding:"-"`
	DeprecatedPassword *string `json:"password" binding:"-"`
}

type ProfileUpdateInput struct {
	FirstName string
	LastName  string
	Phone     string
	AvatarURL string
}

type UserResponse struct {
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

type ProfileUpdateRequestDoc struct {
	FirstName string `json:"first_name" example:"Ada"`
	LastName  string `json:"last_name" example:"Lovelace"`
	Phone     string `json:"phone" example:"+12025550123"`
	AvatarURL string `json:"avatar_url" example:"https://example.com/avatar.png"`
}
