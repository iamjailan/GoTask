package task

import (
	"encoding/json"
	"time"
)

type Request struct {
	Title       string     `json:"title" binding:"required,min=1,max=255"`
	Description string     `json:"description" binding:"omitempty,max=5000"`
	Status      string     `json:"status" binding:"omitempty,oneof=pending in_progress completed archived"`
	Priority    string     `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	DueDate     *time.Time `json:"due_date"`
	Completed   bool       `json:"completed"`
}

// UpdateRequest contains only the task fields that should be changed.
type UpdateRequest struct {
	Title       *string      `json:"title" binding:"omitempty,min=1,max=255"`
	Description *string      `json:"description" binding:"omitempty,max=5000"`
	Status      *string      `json:"status" binding:"omitempty,oneof=pending in_progress completed archived"`
	Priority    *string      `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	DueDate     OptionalTime `json:"due_date" swaggertype:"string" format:"date-time"`
	Completed   *bool        `json:"completed"`
}

func (r UpdateRequest) HasUpdates() bool {
	return r.Title != nil || r.Description != nil || r.Status != nil || r.Priority != nil || r.DueDate.Set || r.Completed != nil
}

// OptionalTime distinguishes an omitted due_date from an explicit null value.
type OptionalTime struct {
	Set   bool
	Value *time.Time
}

func (o *OptionalTime) UnmarshalJSON(value []byte) error {
	o.Set = true
	if string(value) == "null" {
		o.Value = nil
		return nil
	}
	var parsed time.Time
	if err := json.Unmarshal(value, &parsed); err != nil {
		return err
	}
	o.Value = &parsed
	return nil
}

type Response struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Completed   bool       `json:"completed"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
