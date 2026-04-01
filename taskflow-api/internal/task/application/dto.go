package application

import "time"

type CreateTaskDTO struct {
	Title       string
	Description string
	ProjectID   string
}

type MoveTaskDTO struct {
	TaskID    string
	NewStatus string
}

type TaskDTO struct {
	ID          string
	Title       string
	Description string
	Status      string
	AssigneeID  *string
	ProjectID   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
