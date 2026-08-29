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

type ChangeEmailRequestDoc struct {
	Email           string `json:"email" example:"new@example.com"`
	CurrentPassword string `json:"current_password" example:"current-password"`
}

type ChangePasswordRequestDoc struct {
	CurrentPassword string `json:"current_password" example:"current-password"`
	NewPassword     string `json:"new_password" example:"new-password"`
}
