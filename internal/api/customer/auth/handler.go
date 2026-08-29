package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"gotask/internal/api/customer/auth/models"
	"gotask/internal/api/customer/auth/utils"
	authtypes "gotask/internal/types/auth"
	response "gotask/internal/utils/response"

	"github.com/gin-gonic/gin"
)

type Handler struct{ service Service }

func NewHandler(service Service) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(router *gin.Engine, emailRateLimit gin.HandlerFunc) {
	routes := router.Group("/customer/auth")
	routes.POST("/register", emailRateLimit, h.register)
	routes.POST("/confirm-email", h.confirmEmail)
	routes.POST("/login", h.login)
}

// register godoc
// @Summary Start customer registration
// @Description Creates a pending customer registration and sends an email verification code.
// @Tags Customer authentication
// @Accept json
// @Produce json
// @Param request body authtypes.RegisterRequest true "Registration details"
// @Success 202 {object} response.SuccessEnvelope
// @Failure 400 {object} response.ErrorEnvelope
// @Failure 409 {object} response.ErrorEnvelope
// @Failure 429 {object} response.ErrorEnvelope
// @Failure 500 {object} response.ErrorEnvelope
// @Router /customer/auth/register [post]
func (h *Handler) register(c *gin.Context) {
	var req authtypes.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "first_name, last_name, a valid email, and a password of 8-72 characters are required")
		return
	}
	model, _, err := h.service.Register(c.Request.Context(), authtypes.RegisterInput{
		FirstName: strings.TrimSpace(req.FirstName), LastName: strings.TrimSpace(req.LastName),
		Email: req.Email, Password: req.Password, Phone: strings.TrimSpace(req.Phone), AvatarURL: req.AvatarURL,
	})
	if err != nil {
		h.writeError(c, err)
		return
	}
	response.JSON(c, http.StatusAccepted, authtypes.RegistrationResponse{
		Email: model.Email, VerificationRequired: true,
		Message: "verification code sent to your email; confirm it before logging in",
	})
}

// confirmEmail godoc
// @Summary Confirm customer email
// @Description Verifies the email code and returns a customer JWT.
// @Tags Customer authentication
// @Accept json
// @Produce json
// @Param request body authtypes.ConfirmEmailRequest true "Email verification details"
// @Success 200 {object} response.SuccessEnvelope
// @Failure 400 {object} response.ErrorEnvelope
// @Failure 429 {object} response.ErrorEnvelope
// @Failure 500 {object} response.ErrorEnvelope
// @Router /customer/auth/confirm-email [post]
func (h *Handler) confirmEmail(c *gin.Context) {
	var req authtypes.ConfirmEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "a valid email and 6-digit verification code are required")
		return
	}
	model, token, err := h.service.ConfirmEmail(c.Request.Context(), req.Email, req.Code)
	if err != nil {
		h.writeError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, authtypes.VerificationResponse{
		Message: "email verified successfully", Customer: newCustomerResponse(model),
		Token: token, ExpiresAt: time.Now().UTC().Add(utils.TokenLifetime),
	})
}

// login godoc
// @Summary Customer login
// @Description Authenticates a verified customer and returns a customer JWT.
// @Tags Customer authentication
// @Accept json
// @Produce json
// @Param request body authtypes.LoginRequest true "Customer credentials"
// @Success 200 {object} response.SuccessEnvelope
// @Failure 400 {object} response.ErrorEnvelope
// @Failure 401 {object} response.ErrorEnvelope
// @Failure 429 {object} response.ErrorEnvelope
// @Failure 500 {object} response.ErrorEnvelope
// @Router /customer/auth/login [post]
func (h *Handler) login(c *gin.Context) {
	var req authtypes.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "valid email and password are required")
		return
	}
	model, token, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		h.writeError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, newAuthResponse(model, token))
}

func (h *Handler) writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrEmailExists):
		response.Error(c, http.StatusConflict, "email already exists")
	case errors.Is(err, ErrInvalidCredentials), errors.Is(err, ErrInactive), errors.Is(err, ErrEmailNotVerified):
		response.Error(c, http.StatusUnauthorized, "invalid credentials")
	case errors.Is(err, ErrInvalidOTP), errors.Is(err, ErrOTPExpired):
		response.Error(c, http.StatusBadRequest, "invalid or expired verification code")
	default:
		response.Error(c, http.StatusInternalServerError, "internal server error")
	}
}

func newAuthResponse(model models.Model, token string) authtypes.AuthResponse {
	return authtypes.AuthResponse{Token: token, ExpiresAt: time.Now().UTC().Add(utils.TokenLifetime), Customer: authtypes.CustomerResponse{
		ID: model.ID, FirstName: model.FirstName, LastName: model.LastName, Email: model.Email,
		Phone: model.Phone, AvatarURL: model.AvatarURL, Role: model.Role, IsActive: model.IsActive,
		CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt,
	}}
}

func newCustomerResponse(model models.Model) authtypes.CustomerResponse {
	return authtypes.CustomerResponse{ID: model.ID, FirstName: model.FirstName, LastName: model.LastName,
		Email: model.Email, Phone: model.Phone, AvatarURL: model.AvatarURL, Role: model.Role,
		IsActive: model.IsActive, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}
}
