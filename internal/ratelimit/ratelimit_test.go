package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestMiddlewareLimitsEachClientIndependently(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	if err := router.SetTrustedProxies(nil); err != nil {
		t.Fatal(err)
	}
	router.Use(New(2, time.Second).Middleware())
	router.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	request := func(client string) *httptest.ResponseRecorder {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = client
		router.ServeHTTP(recorder, req)
		return recorder
	}

	for i := 0; i < 2; i++ {
		if recorder := request("198.51.100.10:12345"); recorder.Code != http.StatusOK {
			t.Fatalf("request %d returned %d, want %d", i+1, recorder.Code, http.StatusOK)
		}
	}

	limited := request("198.51.100.10:12345")
	if limited.Code != http.StatusTooManyRequests {
		t.Fatalf("limited request returned %d, want %d", limited.Code, http.StatusTooManyRequests)
	}
	if retryAfter := limited.Header().Get("Retry-After"); retryAfter != "1" {
		t.Fatalf("Retry-After = %q, want %q", retryAfter, "1")
	}

	if recorder := request("198.51.100.11:12345"); recorder.Code != http.StatusOK {
		t.Fatalf("different client returned %d, want %d", recorder.Code, http.StatusOK)
	}
}
