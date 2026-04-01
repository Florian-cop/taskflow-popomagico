package infrastructure

import "time"

type TaskModel struct {
	ID          string  `gorm:"primaryKey"`
	Title       string
	Description string
	Status      string
	AssigneeID  *string
	ProjectID   string `gorm:"index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (TaskModel) TableName() string { return "tasks" }
