package dto

import "time"

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type MoveTaskRequest struct {
	Status string `json:"status"`
}

type TaskResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	AssigneeID  *string   `json:"assigneeId"`
	ProjectID   string    `json:"projectId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
