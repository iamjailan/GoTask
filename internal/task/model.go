package task

import "time"

type Model struct {
	ID         string    `gorm:"primaryKey;type:varchar(50);default:concat('tsk_', gen_random_uuid())" json:"id"`
	CustomerID string    `gorm:"type:varchar(50);index" json:"-"`
	Title      string    `gorm:"not null" json:"title"`
	Completed  bool      `gorm:"not null;default:false" json:"completed"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Model) TableName() string { return "tasks" }
