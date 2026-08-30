package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gotask/internal/api/customer/auth"
)

func TestSwaggerRoutesRequireBasicAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	registerSwaggerRoutes(router, "docs-user", "docs-password")

	unauthenticated := httptest.NewRecorder()
	router.ServeHTTP(unauthenticated, httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil))
	if unauthenticated.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated Swagger request returned %d, want %d", unauthenticated.Code, http.StatusUnauthorized)
	}

	authenticatedRequest := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
	authenticatedRequest.SetBasicAuth("docs-user", "docs-password")
	authenticated := httptest.NewRecorder()
	router.ServeHTTP(authenticated, authenticatedRequest)
	if authenticated.Code != http.StatusOK {
		t.Fatalf("authenticated Swagger request returned %d, want %d", authenticated.Code, http.StatusOK)
	}
}

func TestMissingRouteRequiresAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.NoRoute(auth.JWTMiddleware("test-secret"), notFound)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/missing", nil))

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
	const want = `{"success":false,"statusCode":401,"error":"unauthorized"}`
	if got := recorder.Body.String(); got != want {
		t.Errorf("response = %s, want %s", got, want)
	}
}

func TestAuthenticatedMissingRouteReturnsJSONNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.NoRoute(auth.JWTMiddleware("test-secret"), notFound)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": "cus_123"})
	tokenString, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
	const want = `{"success":false,"statusCode":404,"error":"route not found"}`
	if got := recorder.Body.String(); got != want {
		t.Errorf("response = %s, want %s", got, want)
	}
}
