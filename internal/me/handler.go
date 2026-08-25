package me

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gotask/internal/auth"
	"gotask/internal/auth/models"
	apiresponse "gotask/internal/response"
)

type Handler struct{ service Service }

type updateRequest struct {
	FirstName string `json:"first_name" binding:"omitempty,max=100"`
	LastName  string `json:"last_name" binding:"omitempty,max=100"`
	Email     string `json:"email" binding:"omitempty,email,max=255"`
	Phone     string `json:"phone" binding:"omitempty,max=30"`
	AvatarURL string `json:"avatar_url" binding:"omitempty,url"`
}

type response struct {
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

func NewHandler(service Service) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(router *gin.Engine, middleware gin.HandlerFunc) {
	routes := router.Group("/api/v1/me", middleware)
	routes.GET("", h.get)
	routes.PUT("", h.update)
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
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, http.StatusBadRequest, "invalid user profile")
		return
	}
	model, err := h.service.Update(c.Request.Context(), id, UpdateInput{
		FirstName: req.FirstName, LastName: req.LastName, Email: req.Email,
		Phone: req.Phone, AvatarURL: req.AvatarURL,
	})
	if err != nil {
		h.writeError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, newResponse(model))
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
	default:
		apiresponse.Error(c, http.StatusInternalServerError, "internal server error")
	}
}

func newResponse(model models.Model) response {
	return response{ID: model.ID, FirstName: model.FirstName, LastName: model.LastName,
		Email: model.Email, Phone: model.Phone, AvatarURL: model.AvatarURL, Role: model.Role,
		IsActive: model.IsActive, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}
}
