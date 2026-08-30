package task

import "context"

type StatisticsService interface {
	Summary(context.Context, string) (StatisticsSummary, error)
}

type statisticsService struct{ repo StatisticsRepository }

func NewStatisticsService(repo StatisticsRepository) StatisticsService {
	return &statisticsService{repo: repo}
}

func (s *statisticsService) Summary(ctx context.Context, customerID string) (StatisticsSummary, error) {
	return s.repo.Summary(ctx, customerID)
}
