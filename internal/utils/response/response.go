package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type envelope struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

// SuccessEnvelope documents the common shape used for successful API responses.
type SuccessEnvelope struct {
	Success bool `json:"success" example:"true"`
	Data    any  `json:"data"`
}

// ErrorEnvelope documents the common shape used for error API responses.
type ErrorEnvelope struct {
	Success bool      `json:"success" example:"false"`
	Data    ErrorData `json:"data"`
}

type UnauthorizedEnvelope struct {
	Success    bool   `json:"success" example:"false"`
	StatusCode int    `json:"statusCode" example:"401"`
	Error      string `json:"error" example:"unauthorized"`
}

// ErrorData contains a human-readable error message.
type ErrorData struct {
	Error string `json:"error" example:"unauthorized"`
}

func JSON(c *gin.Context, status int, data any) {
	c.JSON(status, envelope{Success: true, Data: data})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, envelope{Success: false, Data: gin.H{"error": message}})
}

func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, UnauthorizedEnvelope{
		Success:    false,
		StatusCode: http.StatusUnauthorized,
		Error:      "unauthorized",
	})
}
