package task

import (
	"context"
	"errors"
	"testing"
)

func TestServiceRecordsTaskLifecycleEvents(t *testing.T) {
	ctx := context.Background()
	repo := &fakeTaskRepository{models: map[string]Model{}}
	service := NewService(repo)

	created, err := service.Create(ctx, CreateInput{CustomerID: "cus_123", Title: "Write tests", Completed: true})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}
	if got, want := eventTypes(repo.events), []string{StatisticEventCreated, StatisticEventCompleted}; !sameStrings(got, want) {
		t.Errorf("create events = %v, want %v", got, want)
	}

	_, err = service.Update(ctx, "cus_123", created.ID, UpdateInput{Title: stringPointer("Write tests"), Completed: boolPointer(false)})
	if err != nil {
		t.Fatalf("reopen task: %v", err)
	}
	_, err = service.Update(ctx, "cus_123", created.ID, UpdateInput{Title: stringPointer("Write tests"), Completed: boolPointer(true)})
	if err != nil {
		t.Fatalf("complete task again: %v", err)
	}
	if err := service.Delete(ctx, "cus_123", created.ID); err != nil {
		t.Fatalf("delete task: %v", err)
	}

	want := []string{
		StatisticEventCreated, StatisticEventCompleted,
		StatisticEventUpdated,
		StatisticEventUpdated, StatisticEventCompleted,
		StatisticEventDeleted,
	}
	if got := eventTypes(repo.events); !sameStrings(got, want) {
		t.Errorf("events = %v, want %v", got, want)
	}
	for _, event := range repo.events {
		if event.CustomerID != "cus_123" || event.TaskID != created.ID {
			t.Errorf("event has incorrect ownership or task ID: %#v", event)
		}
	}
}

func TestServiceDoesNotRecordCompletionWhenTaskStaysCompleted(t *testing.T) {
	ctx := context.Background()
	repo := &fakeTaskRepository{models: map[string]Model{
		"tsk_123": {ID: "tsk_123", CustomerID: "cus_123", Title: "Done", Completed: true, Status: "completed"},
	}}
	service := NewService(repo)

	if _, err := service.Update(ctx, "cus_123", "tsk_123", UpdateInput{Title: stringPointer("Still done"), Completed: boolPointer(true)}); err != nil {
		t.Fatalf("update task: %v", err)
	}
	if got, want := eventTypes(repo.events), []string{StatisticEventUpdated}; !sameStrings(got, want) {
		t.Errorf("events = %v, want %v", got, want)
	}
}

func TestServicePartialUpdatePreservesUnspecifiedFields(t *testing.T) {
	repo := &fakeTaskRepository{models: map[string]Model{
		"tsk_123": {
			ID: "tsk_123", CustomerID: "cus_123", Title: "Keep title", Description: "Keep description",
			Status: "in_progress", Priority: "high", Completed: false,
		},
	}}

	updated, err := NewService(repo).Update(context.Background(), "cus_123", "tsk_123", UpdateInput{Completed: boolPointer(true)})
	if err != nil {
		t.Fatalf("update task: %v", err)
	}
	if updated.Title != "Keep title" || updated.Description != "Keep description" || updated.Priority != "high" {
		t.Errorf("unspecified fields changed: %#v", updated)
	}
	if !updated.Completed || updated.Status != "completed" {
		t.Errorf("completion update was not applied: %#v", updated)
	}
}

func TestServiceRollsBackTaskWhenStatisticsWriteFails(t *testing.T) {
	repo := &fakeTaskRepository{
		models:    map[string]Model{},
		recordErr: errors.New("statistics unavailable"),
	}

	_, err := NewService(repo).Create(context.Background(), CreateInput{
		CustomerID: "cus_123",
		Title:      "Do not persist",
	})
	if !errors.Is(err, repo.recordErr) {
		t.Fatalf("create error = %v, want statistics error", err)
	}
	if len(repo.models) != 0 || len(repo.events) != 0 {
		t.Errorf("transaction committed after event error: models=%v events=%v", repo.models, repo.events)
	}
}

type fakeTaskRepository struct {
	models    map[string]Model
	events    []Statistic
	recordErr error
}

func (r *fakeTaskRepository) Transaction(_ context.Context, fn func(Repository, StatisticsRepository) error) error {
	tx := &fakeTaskRepository{models: cloneModels(r.models), events: append([]Statistic(nil), r.events...), recordErr: r.recordErr}
	if err := fn(tx, tx); err != nil {
		return err
	}
	r.models = tx.models
	r.events = tx.events
	return nil
}

func (r *fakeTaskRepository) Create(_ context.Context, model *Model) error {
	if model.ID == "" {
		model.ID = "tsk_123"
	}
	r.models[model.ID] = *model
	return nil
}

func (r *fakeTaskRepository) List(_ context.Context, _ string) ([]Model, error) { return nil, nil }

func (r *fakeTaskRepository) Get(_ context.Context, customerID, id string) (Model, error) {
	model, ok := r.models[id]
	if !ok || model.CustomerID != customerID {
		return Model{}, ErrNotFound
	}
	return model, nil
}

func (r *fakeTaskRepository) Update(ctx context.Context, customerID, id string, input UpdateInput) (UpdateResult, error) {
	model, err := r.Get(ctx, customerID, id)
	if err != nil {
		return UpdateResult{}, err
	}
	previous := model
	if input.Title != nil {
		model.Title = *input.Title
	}
	if input.Description != nil {
		model.Description = *input.Description
	}
	if input.Status != nil {
		model.Status = *input.Status
	}
	if input.Priority != nil {
		model.Priority = *input.Priority
	}
	if input.DueDateSet {
		model.DueDate = input.DueDate
	}
	if input.Completed != nil {
		model.Completed = *input.Completed
		if model.Completed {
			model.Status = "completed"
		} else if model.Status == "completed" {
			model.Status = "pending"
		}
	}
	r.models[id] = model
	return UpdateResult{Previous: previous, Current: model}, nil
}

func (r *fakeTaskRepository) Delete(_ context.Context, customerID, id string) error {
	if _, err := r.Get(context.Background(), customerID, id); err != nil {
		return err
	}
	delete(r.models, id)
	return nil
}

func (r *fakeTaskRepository) Record(_ context.Context, statistic Statistic) error {
	if r.recordErr != nil {
		return r.recordErr
	}
	r.events = append(r.events, statistic)
	return nil
}

func (r *fakeTaskRepository) Summary(_ context.Context, _ string) (StatisticsSummary, error) {
	return StatisticsSummary{}, nil
}

func eventTypes(events []Statistic) []string {
	types := make([]string, 0, len(events))
	for _, event := range events {
		types = append(types, event.EventType)
	}
	return types
}

func sameStrings(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

func cloneModels(models map[string]Model) map[string]Model {
	clone := make(map[string]Model, len(models))
	for id, model := range models {
		clone[id] = model
	}
	return clone
}

func stringPointer(value string) *string { return &value }

func boolPointer(value bool) *bool { return &value }

var _ Repository = (*fakeTaskRepository)(nil)
var _ StatisticsRepository = (*fakeTaskRepository)(nil)
