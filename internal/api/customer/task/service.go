package task

import (
	"context"
	"time"
)

type Service interface {
	Create(context.Context, CreateInput) (Model, error)
	List(context.Context, string) ([]Model, error)
	Get(context.Context, string, string) (Model, error)
	Update(context.Context, string, string, UpdateInput) (Model, error)
	Delete(context.Context, string, string) error
}

type service struct{ repo Repository }

type CreateInput struct {
	CustomerID  string
	Title       string
	Description string
	Status      string
	Priority    string
	DueDate     *time.Time
	Completed   bool
}

type UpdateInput struct {
	Title       string
	Description string
	Status      string
	Priority    string
	DueDate     *time.Time
	Completed   bool
}

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) Create(ctx context.Context, input CreateInput) (Model, error) {
	status := input.Status
	if status == "" {
		status = "pending"
	}
	priority := input.Priority
	if priority == "" {
		priority = "medium"
	}
	model := Model{
		CustomerID: input.CustomerID, Title: input.Title, Description: input.Description,
		Status: status, Priority: priority, DueDate: input.DueDate, Completed: input.Completed,
	}
	if input.Completed {
		now := time.Now().UTC()
		model.Status = "completed"
		model.CompletedAt = &now
	}
	if err := s.repo.Transaction(ctx, func(taskRepo Repository, statisticsRepo StatisticsRepository) error {
		if err := taskRepo.Create(ctx, &model); err != nil {
			return err
		}
		if err := statisticsRepo.Record(ctx, Statistic{CustomerID: input.CustomerID, TaskID: model.ID, EventType: StatisticEventCreated}); err != nil {
			return err
		}
		if model.Completed {
			return statisticsRepo.Record(ctx, Statistic{CustomerID: input.CustomerID, TaskID: model.ID, EventType: StatisticEventCompleted})
		}
		return nil
	}); err != nil {
		return Model{}, err
	}
	return model, nil
}

func (s *service) List(ctx context.Context, customerID string) ([]Model, error) {
	return s.repo.List(ctx, customerID)
}

func (s *service) Get(ctx context.Context, customerID, id string) (Model, error) {
	return s.repo.Get(ctx, customerID, id)
}

func (s *service) Update(ctx context.Context, customerID, id string, input UpdateInput) (Model, error) {
	var updated Model
	err := s.repo.Transaction(ctx, func(taskRepo Repository, statisticsRepo StatisticsRepository) error {
		result, err := taskRepo.Update(ctx, customerID, id, &Model{
			Title: input.Title, Description: input.Description, Status: input.Status,
			Priority: input.Priority, DueDate: input.DueDate, Completed: input.Completed,
		})
		if err != nil {
			return err
		}
		updated = result.Current
		if err := statisticsRepo.Record(ctx, Statistic{CustomerID: customerID, TaskID: id, EventType: StatisticEventUpdated}); err != nil {
			return err
		}
		if !result.Previous.Completed && result.Current.Completed {
			return statisticsRepo.Record(ctx, Statistic{CustomerID: customerID, TaskID: id, EventType: StatisticEventCompleted})
		}
		return nil
	})
	if err != nil {
		return Model{}, err
	}
	return updated, nil
}

func (s *service) Delete(ctx context.Context, customerID, id string) error {
	return s.repo.Transaction(ctx, func(taskRepo Repository, statisticsRepo StatisticsRepository) error {
		if _, err := taskRepo.Get(ctx, customerID, id); err != nil {
			return err
		}
		if err := statisticsRepo.Record(ctx, Statistic{CustomerID: customerID, TaskID: id, EventType: StatisticEventDeleted}); err != nil {
			return err
		}
		return taskRepo.Delete(ctx, customerID, id)
	})
}
