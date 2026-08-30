package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
	tasktypes "gotask/internal/types/task"
	response "gotask/internal/utils/response"
)

type StatisticsHandler struct{ service StatisticsService }

func NewStatisticsHandler(service StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{service: service}
}

func (h *StatisticsHandler) RegisterRoutes(router *gin.Engine, middleware gin.HandlerFunc) {
	router.GET("/customer/tasks/statistics", middleware, h.summary)
}

// summary godoc
// @Summary Get task statistics
// @Description Returns historical task-event totals and the 20 most recent events for the authenticated customer.
// @Tags Customer tasks
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessEnvelope
// @Failure 401 {object} response.ErrorEnvelope
// @Failure 429 {object} response.ErrorEnvelope
// @Failure 500 {object} response.ErrorEnvelope
// @Router /customer/tasks/statistics [get]
func (h *StatisticsHandler) summary(c *gin.Context) {
	customerID, ok := currentCustomerID(c)
	if !ok {
		return
	}
	summary, err := h.service.Summary(c.Request.Context(), customerID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}
	response.JSON(c, http.StatusOK, newStatisticsResponse(summary))
}

func newStatisticsResponse(summary StatisticsSummary) tasktypes.StatisticsResponse {
	activity := make([]tasktypes.StatisticActivity, 0, len(summary.RecentActivity))
	for _, statistic := range summary.RecentActivity {
		activity = append(activity, tasktypes.StatisticActivity{
			Event: statistic.EventType, TaskID: statistic.TaskID, OccurredAt: statistic.OccurredAt,
		})
	}
	return tasktypes.StatisticsResponse{
		Counts: tasktypes.StatisticsCounts{
			TasksCreated: summary.Counts.TasksCreated, TasksUpdated: summary.Counts.TasksUpdated,
			TasksCompleted: summary.Counts.TasksCompleted, TasksDeleted: summary.Counts.TasksDeleted,
		},
		RecentActivity: activity,
	}
}
