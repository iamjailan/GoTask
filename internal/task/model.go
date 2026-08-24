package task

import "time"

type Model struct {
	ID          string     `gorm:"primaryKey;type:varchar(50);default:concat('tsk_', gen_random_uuid())" json:"id"`
	CustomerID  string     `gorm:"type:varchar(50);index" json:"-"`
	Title       string     `gorm:"not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Status      string     `gorm:"type:varchar(30);not null;default:pending;index" json:"status"`
	Priority    string     `gorm:"type:varchar(20);not null;default:medium" json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Completed   bool       `gorm:"not null;default:false" json:"completed"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (Model) TableName() string { return "tasks" }
