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
	Success    bool   `json:"success" example:"false"`
	StatusCode int    `json:"statusCode" example:"400"`
	Error      string `json:"error" example:"title is required"`
}

func JSON(c *gin.Context, status int, data any) {
	c.JSON(status, envelope{Success: true, Data: data})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, ErrorEnvelope{Success: false, StatusCode: status, Error: message})
}

func Unauthorized(c *gin.Context) {
	Error(c, http.StatusUnauthorized, "unauthorized")
}
