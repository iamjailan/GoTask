package task

import "time"

type Request struct {
	Title       string     `json:"title" binding:"required,min=1,max=255"`
	Description string     `json:"description" binding:"omitempty,max=5000"`
	Status      string     `json:"status" binding:"omitempty,oneof=pending in_progress completed archived"`
	Priority    string     `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	DueDate     *time.Time `json:"due_date"`
	Completed   bool       `json:"completed"`
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
