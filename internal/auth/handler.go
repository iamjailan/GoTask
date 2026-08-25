package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gotask/internal/auth/models"
	"gotask/internal/auth/types"
	"gotask/internal/auth/utils"
	apiresponse "gotask/internal/response"
)

type Handler struct{ service Service }

type registerRequest struct {
	FirstName string `json:"first_name" binding:"required,max=100"`
	LastName  string `json:"last_name" binding:"required,max=100"`
	Email     string `json:"email" binding:"required,email,max=255"`
	Password  string `json:"password" binding:"required,min=8,max=72"`
	Phone     string `json:"phone" binding:"omitempty,max=30"`
	AvatarURL string `json:"avatar_url" binding:"omitempty,url"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type authResponse struct {
	Token     string           `json:"token"`
	ExpiresAt time.Time        `json:"expires_at"`
	Customer  customerResponse `json:"customer"`
}

type customerResponse struct {
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

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	routes := router.Group("/api/v1/auth")
	routes.POST("/register", h.register)
	routes.POST("/login", h.login)
}

func (h *Handler) register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, http.StatusBadRequest, "first_name, last_name, a valid email, and a password of 8-72 characters are required")
		return
	}
	model, token, err := h.service.Register(c.Request.Context(), types.RegisterInput{
		FirstName: strings.TrimSpace(req.FirstName), LastName: strings.TrimSpace(req.LastName),
		Email: req.Email, Password: req.Password, Phone: strings.TrimSpace(req.Phone), AvatarURL: req.AvatarURL,
	})
	if err != nil {
		h.writeError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusCreated, newAuthResponse(model, token))
}

func (h *Handler) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresponse.Error(c, http.StatusBadRequest, "valid email and password are required")
		return
	}
	model, token, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		h.writeError(c, err)
		return
	}
	apiresponse.JSON(c, http.StatusOK, newAuthResponse(model, token))
}

func (h *Handler) writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrEmailExists):
		apiresponse.Error(c, http.StatusConflict, "email already exists")
	case errors.Is(err, ErrInvalidCredentials), errors.Is(err, ErrInactive):
		apiresponse.Error(c, http.StatusUnauthorized, "invalid credentials")
	default:
		apiresponse.Error(c, http.StatusInternalServerError, "internal server error")
	}
}

func newAuthResponse(model models.Model, token string) authResponse {
	return authResponse{Token: token, ExpiresAt: time.Now().UTC().Add(utils.TokenLifetime), Customer: customerResponse{
		ID: model.ID, FirstName: model.FirstName, LastName: model.LastName, Email: model.Email,
		Phone: model.Phone, AvatarURL: model.AvatarURL, Role: model.Role, IsActive: model.IsActive,
		CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt,
	}}
}
