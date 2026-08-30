package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMiddlewareHandlesAllowedOriginPreflight(t *testing.T) {
	t.Parallel()

	middleware, err := New(Config{
		AllowedOrigins:   "https://app.example.com",
		AllowedMethods:   "GET, POST",
		AllowedHeaders:   "Authorization, Content-Type",
		ExposedHeaders:   "X-Request-ID",
		AllowCredentials: "true",
		MaxAge:           "1h",
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	router := gin.New()
	router.Use(middleware)
	router.POST("/customer/tasks", func(c *gin.Context) { c.Status(http.StatusCreated) })

	request := httptest.NewRequest(http.MethodOptions, "/customer/tasks", nil)
	request.Header.Set("Origin", "https://app.example.com")
	request.Header.Set("Access-Control-Request-Method", http.MethodPost)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
	if value := response.Header().Get("Access-Control-Allow-Origin"); value != "https://app.example.com" {
		t.Fatalf("Access-Control-Allow-Origin = %q", value)
	}
	if value := response.Header().Get("Access-Control-Allow-Credentials"); value != "true" {
		t.Fatalf("Access-Control-Allow-Credentials = %q", value)
	}
	if value := response.Header().Get("Access-Control-Max-Age"); value != "3600" {
		t.Fatalf("Access-Control-Max-Age = %q", value)
	}
}

func TestMiddlewareDoesNotAllowUnknownOrigin(t *testing.T) {
	t.Parallel()

	middleware, err := New(Config{
		AllowedOrigins:   "https://app.example.com",
		AllowedMethods:   "GET",
		AllowedHeaders:   "Authorization",
		AllowCredentials: "false",
		MaxAge:           "1h",
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	router := gin.New()
	router.Use(middleware)
	router.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	request.Header.Set("Origin", "https://untrusted.example.com")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if value := response.Header().Get("Access-Control-Allow-Origin"); value != "" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want empty", value)
	}
}

func TestNewRejectsWildcardWithCredentials(t *testing.T) {
	t.Parallel()

	_, err := New(Config{
		AllowedOrigins:   "*",
		AllowCredentials: "true",
		MaxAge:           "1h",
	})
	if err == nil {
		t.Fatal("New() error = nil, want error")
	}
}
