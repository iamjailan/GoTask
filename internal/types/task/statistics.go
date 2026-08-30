package task

import "time"

type StatisticsResponse struct {
	Counts         StatisticsCounts    `json:"counts"`
	RecentActivity []StatisticActivity `json:"recent_activity"`
}

type StatisticsCounts struct {
	TasksCreated   int64 `json:"tasks_created"`
	TasksUpdated   int64 `json:"tasks_updated"`
	TasksCompleted int64 `json:"tasks_completed"`
	TasksDeleted   int64 `json:"tasks_deleted"`
}

type StatisticActivity struct {
	Event      string    `json:"event"`
	TaskID     string    `json:"task_id"`
	OccurredAt time.Time `json:"occurred_at"`
}
