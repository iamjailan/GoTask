package me

type ChangeEmailRequest struct {
	Email           string `json:"email" binding:"required,email,max=255"`
	CurrentPassword string `json:"current_password" binding:"required"`
}

type ChangeEmailInput struct {
	Email           string
	CurrentPassword string
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=72"`
}

type ChangePasswordInput struct {
	CurrentPassword string
	NewPassword     string
}
