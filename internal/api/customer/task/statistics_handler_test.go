package task

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gotask/internal/api/customer/auth"
)

func TestStatisticsHandlerReturnsSummary(t *testing.T) {
	gin.SetMode(gin.TestMode)
	occurredAt := time.Date(2026, time.August, 30, 12, 0, 0, 0, time.UTC)
	service := &stubStatisticsService{summary: StatisticsSummary{
		Counts:         StatisticsCounts{TasksCreated: 3, TasksUpdated: 2, TasksCompleted: 1, TasksDeleted: 1},
		RecentActivity: []Statistic{{EventType: StatisticEventCompleted, TaskID: "tsk_123", OccurredAt: occurredAt}},
	}}
	router := gin.New()
	NewStatisticsHandler(service).RegisterRoutes(router, customerMiddleware)
	NewHandler(&stubTaskService{}).RegisterRoutes(router, customerMiddleware)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/customer/tasks/statistics", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if service.customerID != "cus_123" {
		t.Errorf("customer ID = %q, want cus_123", service.customerID)
	}
	var body struct {
		Success bool `json:"success"`
		Data    struct {
			Counts StatisticsCounts `json:"counts"`
			Recent []struct {
				Event  string `json:"event"`
				TaskID string `json:"task_id"`
			} `json:"recent_activity"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !body.Success || body.Data.Counts.TasksCreated != 3 || len(body.Data.Recent) != 1 {
		t.Fatalf("unexpected response: %s", recorder.Body.String())
	}
	if body.Data.Recent[0].Event != StatisticEventCompleted || body.Data.Recent[0].TaskID != "tsk_123" {
		t.Errorf("activity = %#v", body.Data.Recent[0])
	}
}

func TestStatisticsHandlerRequiresAuthenticatedCustomer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewStatisticsHandler(&stubStatisticsService{}).RegisterRoutes(router, func(c *gin.Context) { c.Next() })

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/customer/tasks/statistics", nil))

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
	if got, want := recorder.Body.String(), `{"success":false,"statusCode":401,"error":"unauthorized"}`; got != want {
		t.Errorf("response = %s, want %s", got, want)
	}
}

func TestStatisticsHandlerReturnsServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewStatisticsHandler(&stubStatisticsService{err: errors.New("database unavailable")}).RegisterRoutes(router, customerMiddleware)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/customer/tasks/statistics", nil))

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
}

func customerMiddleware(c *gin.Context) {
	c.Set(auth.CustomerIDKey, "cus_123")
	c.Next()
}

type stubStatisticsService struct {
	summary    StatisticsSummary
	err        error
	customerID string
}

func (s *stubStatisticsService) Summary(_ context.Context, customerID string) (StatisticsSummary, error) {
	s.customerID = customerID
	return s.summary, s.err
}

var _ StatisticsService = (*stubStatisticsService)(nil)

type stubTaskService struct{}

func (*stubTaskService) Create(context.Context, CreateInput) (Model, error) { return Model{}, nil }

func (*stubTaskService) List(context.Context, string) ([]Model, error) { return nil, nil }

func (*stubTaskService) Get(context.Context, string, string) (Model, error) { return Model{}, nil }

func (*stubTaskService) Update(context.Context, string, string, UpdateInput) (Model, error) {
	return Model{}, nil
}

func (*stubTaskService) Delete(context.Context, string, string) error { return nil }

var _ Service = (*stubTaskService)(nil)
