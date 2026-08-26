package me

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gotask/internal/api/customer/auth"
	"gotask/internal/api/customer/auth/models"
	metypes "gotask/internal/types/me"
	response "gotask/internal/utils/response"
)

type Handler struct{ service Service }

// profileUpdateRequestDoc documents the accepted customer profile fields.
type profileUpdateRequestDoc struct {
	FirstName string `json:"first_name" example:"Ada"`
	LastName  string `json:"last_name" example:"Lovelace"`
	Phone     string `json:"phone" example:"+12025550123"`
	AvatarURL string `json:"avatar_url" example:"https://example.com/avatar.png"`
}

// changeEmailRequestDoc documents the email change request body.
type changeEmailRequestDoc struct {
	Email           string `json:"email" example:"new@example.com"`
	CurrentPassword string `json:"current_password" example:"current-password"`
}

// changePasswordRequestDoc documents the password change request body.
type changePasswordRequestDoc struct {
	CurrentPassword string `json:"current_password" example:"current-password"`
	NewPassword     string `json:"new_password" example:"new-password"`
}

func NewHandler(service Service) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(router *gin.Engine, middleware gin.HandlerFunc) {
	routes := router.Group("/customer/me", middleware)
	routes.GET("", h.get)
	routes.PUT("", h.update)
	routes.PUT("/email", h.changeEmail)
	routes.PUT("/password", h.changePassword)
	routes.DELETE("", h.delete)
}

func currentUserID(c *gin.Context) (string, bool) {
	id := strings.TrimSpace(c.GetHeader(auth.CustomerIDHeader))
	if id == "" {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		c.Abort()
		return "", false
	}
	return id, true
}

// get godoc
// @Summary Get the authenticated customer
// @Tags Customer profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessEnvelope
// @Failure 401 {object} response.ErrorEnvelope
// @Failure 500 {object} response.ErrorEnvelope
// @Router /customer/me [get]
func (h *Handler) get(c *gin.Context) {
	id, ok := currentUserID(c)
	if !ok {
		return
	}
	model, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		h.writeError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, newResponse(model))
}

// update godoc
// @Summary Update the authenticated customer profile
// @Tags Customer profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body profileUpdateRequestDoc true "Profile details"
// @Success 200 {object} response.SuccessEnvelope
// @Failure 400 {object} response.ErrorEnvelope
// @Failure 401 {object} response.ErrorEnvelope
// @Failure 409 {object} response.ErrorEnvelope
// @Failure 500 {object} response.ErrorEnvelope
// @Router /customer/me [put]
func (h *Handler) update(c *gin.Context) {
	id, ok := currentUserID(c)
	if !ok {
		return
	}
	var req metypes.ProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user profile")
		return
	}
	if req.DeprecatedEmail != nil || req.DeprecatedPassword != nil {
		response.Error(c, http.StatusBadRequest, "use the dedicated /customer/me/email or /customer/me/password endpoint for credential changes")
		return
	}
	model, err := h.service.UpdateProfile(c.Request.Context(), id, metypes.ProfileUpdateInput{
		FirstName: req.FirstName, LastName: req.LastName,
		Phone: req.Phone, AvatarURL: req.AvatarURL,
	})
	if err != nil {
		h.writeError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, newResponse(model))
}

// changeEmail godoc
// @Summary Change the authenticated customer email
// @Tags Customer profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body changeEmailRequestDoc true "New email and current password"
// @Success 200 {object} response.SuccessEnvelope
// @Failure 400 {object} response.ErrorEnvelope
// @Failure 401 {object} response.ErrorEnvelope
// @Failure 409 {object} response.ErrorEnvelope
// @Failure 502 {object} response.ErrorEnvelope
// @Failure 503 {object} response.ErrorEnvelope
// @Router /customer/me/email [put]
func (h *Handler) changeEmail(c *gin.Context) {
	id, ok := currentUserID(c)
	if !ok {
		return
	}
	var req metypes.ChangeEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "a valid email and current password are required")
		return
	}
	model, err := h.service.ChangeEmail(c.Request.Context(), id, metypes.ChangeEmailInput{
		Email: req.Email, CurrentPassword: req.CurrentPassword,
	})
	if err != nil {
		h.writeError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, newResponse(model))
}

// changePassword godoc
// @Summary Change the authenticated customer password
// @Tags Customer profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body changePasswordRequestDoc true "Current and new password"
// @Success 200 {object} response.SuccessEnvelope
// @Failure 400 {object} response.ErrorEnvelope
// @Failure 401 {object} response.ErrorEnvelope
// @Router /customer/me/password [put]
func (h *Handler) changePassword(c *gin.Context) {
	id, ok := currentUserID(c)
	if !ok {
		return
	}
	var req metypes.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "current_password and a new_password of 8-72 characters are required")
		return
	}
	if err := h.service.ChangePassword(c.Request.Context(), id, metypes.ChangePasswordInput{
		CurrentPassword: req.CurrentPassword, NewPassword: req.NewPassword,
	}); err != nil {
		h.writeError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"message": "password updated successfully"})
}

// delete godoc
// @Summary Delete the authenticated customer
// @Tags Customer profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessEnvelope
// @Failure 401 {object} response.ErrorEnvelope
// @Failure 500 {object} response.ErrorEnvelope
// @Router /customer/me [delete]
func (h *Handler) delete(c *gin.Context) {
	id, ok := currentUserID(c)
	if !ok {
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.writeError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, nil)
}

func (h *Handler) writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		response.Error(c, http.StatusUnauthorized, "unauthorized")
	case errors.Is(err, ErrEmailExist):
		response.Error(c, http.StatusConflict, "email already exists")
	case errors.Is(err, ErrEmailUnchanged):
		response.Error(c, http.StatusBadRequest, "new email must be different from the current email")
	case errors.Is(err, ErrPasswordUnchanged):
		response.Error(c, http.StatusBadRequest, "new password must be different from the current password")
	case errors.Is(err, ErrInvalidCurrentPassword):
		response.Error(c, http.StatusUnauthorized, "current password is incorrect")
	case errors.Is(err, ErrEmailNotificationUnavailable):
		response.Error(c, http.StatusServiceUnavailable, "email change notification is unavailable")
	case errors.Is(err, ErrEmailNotification):
		response.Error(c, http.StatusBadGateway, "email change notification failed")
	default:
		response.Error(c, http.StatusInternalServerError, "internal server error")
	}
}

func newResponse(model models.Model) metypes.UserResponse {
	return metypes.UserResponse{ID: model.ID, FirstName: model.FirstName, LastName: model.LastName,
		Email: model.Email, Phone: model.Phone, AvatarURL: model.AvatarURL, Role: model.Role,
		IsActive: model.IsActive, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}
}
