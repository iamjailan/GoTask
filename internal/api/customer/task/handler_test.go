package task

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"gotask/internal/api/customer/auth"
)

func TestUpdateAcceptsCompletedWithoutTitle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &partialUpdateTaskService{}
	router := gin.New()
	NewHandler(service).RegisterRoutes(router, func(c *gin.Context) {
		c.Set(auth.CustomerIDKey, "cus_123")
		c.Next()
	})

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/customer/tasks/tsk_123", strings.NewReader(`{"completed":true}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if service.input.Title != nil || service.input.Completed == nil || !*service.input.Completed {
		t.Errorf("update input = %#v, want only completed=true", service.input)
	}
}

type partialUpdateTaskService struct{ input UpdateInput }

func (*partialUpdateTaskService) Create(context.Context, CreateInput) (Model, error) {
	return Model{}, nil
}

func (*partialUpdateTaskService) List(context.Context, string) ([]Model, error) { return nil, nil }

func (*partialUpdateTaskService) Get(context.Context, string, string) (Model, error) {
	return Model{}, nil
}

func (s *partialUpdateTaskService) Update(_ context.Context, _ string, _ string, input UpdateInput) (Model, error) {
	s.input = input
	return Model{ID: "tsk_123", Title: "Existing title", Completed: true, Status: "completed"}, nil
}

func (*partialUpdateTaskService) Delete(context.Context, string, string) error { return nil }

var _ Service = (*partialUpdateTaskService)(nil)
