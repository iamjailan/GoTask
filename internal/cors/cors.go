package cors

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Config struct {
	AllowedOrigins   string
	AllowedMethods   string
	AllowedHeaders   string
	ExposedHeaders   string
	AllowCredentials string
	MaxAge           string
}

type policy struct {
	allowedOrigins   map[string]struct{}
	allowAnyOrigin   bool
	allowedMethods   string
	allowedHeaders   string
	exposedHeaders   string
	allowCredentials bool
	maxAge           string
}

// New creates CORS middleware from comma-separated configuration values.
// With no allowed origins configured, it does not emit CORS response headers.
func New(cfg Config) (gin.HandlerFunc, error) {
	allowCredentials, err := strconv.ParseBool(cfg.AllowCredentials)
	if err != nil {
		return nil, fmt.Errorf("parse CORS_ALLOW_CREDENTIALS: %w", err)
	}

	maxAge, err := time.ParseDuration(cfg.MaxAge)
	if err != nil {
		return nil, fmt.Errorf("parse CORS_MAX_AGE: %w", err)
	}
	if maxAge < 0 {
		return nil, fmt.Errorf("CORS_MAX_AGE must not be negative")
	}

	origins := values(cfg.AllowedOrigins)
	p := policy{
		allowedOrigins:   make(map[string]struct{}, len(origins)),
		allowedMethods:   strings.Join(values(cfg.AllowedMethods), ", "),
		allowedHeaders:   strings.Join(values(cfg.AllowedHeaders), ", "),
		exposedHeaders:   strings.Join(values(cfg.ExposedHeaders), ", "),
		allowCredentials: allowCredentials,
		maxAge:           strconv.FormatInt(int64(maxAge/time.Second), 10),
	}

	for _, origin := range origins {
		if origin == "*" {
			p.allowAnyOrigin = true
			continue
		}
		p.allowedOrigins[origin] = struct{}{}
	}
	if p.allowAnyOrigin && p.allowCredentials {
		return nil, fmt.Errorf("CORS_ALLOWED_ORIGINS cannot contain * when CORS_ALLOW_CREDENTIALS is true")
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" || !p.allows(origin) {
			c.Next()
			return
		}

		c.Header("Vary", "Origin")
		if p.allowAnyOrigin {
			c.Header("Access-Control-Allow-Origin", "*")
		} else {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		if p.allowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if p.exposedHeaders != "" {
			c.Header("Access-Control-Expose-Headers", p.exposedHeaders)
		}

		if c.Request.Method == http.MethodOptions && c.GetHeader("Access-Control-Request-Method") != "" {
			c.Header("Access-Control-Allow-Methods", p.allowedMethods)
			c.Header("Access-Control-Allow-Headers", p.allowedHeaders)
			c.Header("Access-Control-Max-Age", p.maxAge)
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}, nil
}

func (p policy) allows(origin string) bool {
	if p.allowAnyOrigin {
		return true
	}
	_, ok := p.allowedOrigins[origin]
	return ok
}

func values(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if item := strings.TrimSpace(part); item != "" {
			result = append(result, item)
		}
	}
	return result
}
