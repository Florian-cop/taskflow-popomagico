package infrastructure

import (
	"taskflow-api/internal/task/domain"
)

func toModel(t *domain.Task) *TaskModel {
	return &TaskModel{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		AssigneeID:  t.AssigneeID,
		ProjectID:   t.ProjectID,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func toDomain(model *TaskModel) *domain.Task {
	return &domain.Task{
		ID:          model.ID,
		Title:       model.Title,
		Description: model.Description,
		Status:      domain.TaskStatus(model.Status),
		AssigneeID:  model.AssigneeID,
		ProjectID:   model.ProjectID,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}
