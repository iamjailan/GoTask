package task

import "time"

const (
	StatisticEventCreated   = "created"
	StatisticEventUpdated   = "updated"
	StatisticEventCompleted = "completed"
	StatisticEventDeleted   = "deleted"
)

// Statistic is an immutable record of a task lifecycle event.
type Statistic struct {
	ID         string    `gorm:"primaryKey;type:varchar(50);default:concat('tst_', gen_random_uuid())"`
	CustomerID string    `gorm:"type:varchar(50);not null;index"`
	TaskID     string    `gorm:"type:varchar(50);not null;index"`
	EventType  string    `gorm:"type:varchar(20);not null"`
	OccurredAt time.Time `gorm:"not null;autoCreateTime"`
}

func (Statistic) TableName() string { return "task_statistics" }

type StatisticsCounts struct {
	TasksCreated   int64 `json:"tasks_created"`
	TasksUpdated   int64 `json:"tasks_updated"`
	TasksCompleted int64 `json:"tasks_completed"`
	TasksDeleted   int64 `json:"tasks_deleted"`
}

type StatisticsSummary struct {
	Counts         StatisticsCounts
	RecentActivity []Statistic
}

type UpdateResult struct {
	Previous Model
	Current  Model
}
