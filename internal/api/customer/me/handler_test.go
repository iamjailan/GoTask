package me

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"gotask/internal/api/customer/auth/models"
	metypes "gotask/internal/types/me"
)

func TestRegisterRoutesAddsCredentialEndpoints(t *testing.T) {
	router := gin.New()
	NewHandler(&stubService{}).RegisterRoutes(router, func(c *gin.Context) { c.Next() }, func(c *gin.Context) { c.Next() })

	routes := make(map[string]bool)
	for _, route := range router.Routes() {
		routes[route.Method+" "+route.Path] = true
	}
	for _, expected := range []string{
		http.MethodPut + " /customer/me/email",
		http.MethodPut + " /customer/me/password",
	} {
		if !routes[expected] {
			t.Errorf("missing route %s", expected)
		}
	}
}

func TestProfileUpdateRejectsCredentialFields(t *testing.T) {
	router := gin.New()
	NewHandler(&stubService{}).RegisterRoutes(router, func(c *gin.Context) { c.Next() }, func(c *gin.Context) { c.Next() })

	req := httptest.NewRequest(http.MethodPut, "/customer/me", strings.NewReader(`{"email":"new@example.com"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "cus_123")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

type stubService struct{}

func (*stubService) Get(context.Context, string) (models.Model, error) {
	return models.Model{}, nil
}

func (*stubService) UpdateProfile(context.Context, string, metypes.ProfileUpdateInput) (models.Model, error) {
	return models.Model{}, nil
}

func (*stubService) ChangeEmail(context.Context, string, metypes.ChangeEmailInput) (models.Model, error) {
	return models.Model{}, nil
}

func (*stubService) ChangePassword(context.Context, string, metypes.ChangePasswordInput) error {
	return nil
}

func (*stubService) Delete(context.Context, string) error { return nil }

var _ Service = (*stubService)(nil)
