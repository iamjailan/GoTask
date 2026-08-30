package task

import (
	"context"

	"gorm.io/gorm"
)

const recentStatisticsLimit = 20

type StatisticsRepository interface {
	Record(context.Context, Statistic) error
	Summary(context.Context, string) (StatisticsSummary, error)
}

type statisticsRepository struct{ db *gorm.DB }

func NewStatisticsRepository(db *gorm.DB) StatisticsRepository {
	return &statisticsRepository{db: db}
}

func (r *statisticsRepository) Record(ctx context.Context, statistic Statistic) error {
	return r.db.WithContext(ctx).Create(&statistic).Error
}

func (r *statisticsRepository) Summary(ctx context.Context, customerID string) (StatisticsSummary, error) {
	var summary StatisticsSummary
	if err := r.db.WithContext(ctx).Model(&Statistic{}).
		Select(`
			COUNT(*) FILTER (WHERE event_type = 'created') AS tasks_created,
			COUNT(*) FILTER (WHERE event_type = 'updated') AS tasks_updated,
			COUNT(*) FILTER (WHERE event_type = 'completed') AS tasks_completed,
			COUNT(*) FILTER (WHERE event_type = 'deleted') AS tasks_deleted
		`).
		Where("customer_id = ?", customerID).
		Scan(&summary.Counts).Error; err != nil {
		return StatisticsSummary{}, err
	}
	if err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).
		Order("occurred_at DESC, id DESC").Limit(recentStatisticsLimit).
		Find(&summary.RecentActivity).Error; err != nil {
		return StatisticsSummary{}, err
	}
	return summary, nil
}
