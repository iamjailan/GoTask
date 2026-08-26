package me

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gotask/internal/api/customer/auth"
	"gotask/internal/api/customer/auth/models"
	metypes "gotask/internal/types/me"
	apiresponse "gotask/internal/utils/response"
)

type Handler struct{ service Service }

func NewHandler(service Service) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(router *gin.Engine, middleware gin.HandlerFunc) {
	routes := router.Group("/api/v1/me", middleware)
	routes.GET("", h.get)
	routes.PUT("", h.update)
	routes.PUT("/email", h.changeEmail)
	routes.PUT("/password", h.changePassword)
	routes.DELETE("", h.delete)
}

func currentUserID(c *gin.Context) (string, bool) {
	id := strings.TrimSpace(c.GetHeader(auth.CustomerIDHeader))
	if id == "" {
		apiresponse.Error(c, http.StatusUnauthorized, "unauthorized")
		c.Abort()
		return "", false
	}
	return id, true
}

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
	apiresponse.JSON(c, http.StatusOK, newResponse(model))
}

func (h *Handler) update(c *gin.Context) {
	id, ok := currentUserID(c)
	if !ok {
		return
	}
	var req metypes.ProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, http.StatusBadRequest, "invalid user profile")
		return
	}
	if req.DeprecatedEmail != nil || req.DeprecatedPassword != nil {
		apiresponse.Error(c, http.StatusBadRequest, "use the dedicated /api/v1/me/email or /api/v1/me/password endpoint for credential changes")
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
	apiresponse.JSON(c, http.StatusOK, newResponse(model))
}

func (h *Handler) changeEmail(c *gin.Context) {
	id, ok := currentUserID(c)
	if !ok {
		return
	}
	var req metypes.ChangeEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, http.StatusBadRequest, "a valid email and current password are required")
		return
	}
	model, err := h.service.ChangeEmail(c.Request.Context(), id, metypes.ChangeEmailInput{
		Email: req.Email, CurrentPassword: req.CurrentPassword,
	})
	if err != nil {
		h.writeError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, newResponse(model))
}

func (h *Handler) changePassword(c *gin.Context) {
	id, ok := currentUserID(c)
	if !ok {
		return
	}
	var req metypes.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, http.StatusBadRequest, "current_password and a new_password of 8-72 characters are required")
		return
	}
	if err := h.service.ChangePassword(c.Request.Context(), id, metypes.ChangePasswordInput{
		CurrentPassword: req.CurrentPassword, NewPassword: req.NewPassword,
	}); err != nil {
		h.writeError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, gin.H{"message": "password updated successfully"})
}

func (h *Handler) delete(c *gin.Context) {
	id, ok := currentUserID(c)
	if !ok {
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.writeError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, nil)
}

func (h *Handler) writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		apiresponse.Error(c, http.StatusUnauthorized, "unauthorized")
	case errors.Is(err, ErrEmailExist):
		apiresponse.Error(c, http.StatusConflict, "email already exists")
	case errors.Is(err, ErrEmailUnchanged):
		apiresponse.Error(c, http.StatusBadRequest, "new email must be different from the current email")
	case errors.Is(err, ErrPasswordUnchanged):
		apiresponse.Error(c, http.StatusBadRequest, "new password must be different from the current password")
	case errors.Is(err, ErrInvalidCurrentPassword):
		apiresponse.Error(c, http.StatusUnauthorized, "current password is incorrect")
	case errors.Is(err, ErrEmailNotificationUnavailable):
		apiresponse.Error(c, http.StatusServiceUnavailable, "email change notification is unavailable")
	case errors.Is(err, ErrEmailNotification):
		apiresponse.Error(c, http.StatusBadGateway, "email change notification failed")
	default:
		apiresponse.Error(c, http.StatusInternalServerError, "internal server error")
	}
}

func newResponse(model models.Model) metypes.UserResponse {
	return metypes.UserResponse{ID: model.ID, FirstName: model.FirstName, LastName: model.LastName,
		Email: model.Email, Phone: model.Phone, AvatarURL: model.AvatarURL, Role: model.Role,
		IsActive: model.IsActive, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}
}
