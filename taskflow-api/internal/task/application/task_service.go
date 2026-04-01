package application

import (
	"context"

	shared "taskflow-api/internal/shared/application"
	"taskflow-api/internal/shared/domain"
	taskDomain "taskflow-api/internal/task/domain"
)

type TaskService struct {
	repo     taskDomain.TaskRepository
	eventBus shared.EventBus
}

func NewTaskService(repo taskDomain.TaskRepository, eventBus shared.EventBus) *TaskService {
	return &TaskService{repo: repo, eventBus: eventBus}
}

func (s *TaskService) CreateTask(ctx context.Context, dto CreateTaskDTO) (*TaskDTO, error) {
	task := taskDomain.NewTask(domain.NewID(), dto.Title, dto.Description, dto.ProjectID)

	if err := s.repo.Save(ctx, task); err != nil {
		return nil, err
	}

	event := taskDomain.NewTaskCreatedEvent(task.ID, task.ProjectID, task.Title)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		return nil, err
	}

	return toDTO(task), nil
}

func (s *TaskService) MoveTask(ctx context.Context, dto MoveTaskDTO) (*TaskDTO, error) {
	task, err := s.repo.FindByID(ctx, dto.TaskID)
	if err != nil {
		return nil, err
	}

	event, err := task.MoveTo(taskDomain.TaskStatus(dto.NewStatus))
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, task); err != nil {
		return nil, err
	}

	if err := s.eventBus.Publish(ctx, event); err != nil {
		return nil, err
	}

	return toDTO(task), nil
}

func (s *TaskService) GetTasksByProject(ctx context.Context, projectID string) ([]*TaskDTO, error) {
	tasks, err := s.repo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	dtos := make([]*TaskDTO, len(tasks))
	for i, t := range tasks {
		dtos[i] = toDTO(t)
	}
	return dtos, nil
}

func toDTO(t *taskDomain.Task) *TaskDTO {
	return &TaskDTO{
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
