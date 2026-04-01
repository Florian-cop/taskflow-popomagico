package application

import (
	"context"

	sharedDomain "taskflow-api/internal/shared/domain"
	"taskflow-api/internal/task/domain"
)

// MockTaskRepository est un repository en mémoire pour les tests.
type MockTaskRepository struct {
	Tasks map[string]*domain.Task
}

func NewMockTaskRepository() *MockTaskRepository {
	return &MockTaskRepository{Tasks: make(map[string]*domain.Task)}
}

func (m *MockTaskRepository) FindByID(ctx context.Context, id string) (*domain.Task, error) {
	task, ok := m.Tasks[id]
	if !ok {
		return nil, sharedDomain.ErrNotFound
	}
	return task, nil
}

func (m *MockTaskRepository) FindByProjectID(ctx context.Context, projectID string) ([]*domain.Task, error) {
	var tasks []*domain.Task
	for _, t := range m.Tasks {
		if t.ProjectID == projectID {
			tasks = append(tasks, t)
		}
	}
	return tasks, nil
}

func (m *MockTaskRepository) Save(ctx context.Context, task *domain.Task) error {
	m.Tasks[task.ID] = task
	return nil
}

func (m *MockTaskRepository) Delete(ctx context.Context, id string) error {
	delete(m.Tasks, id)
	return nil
}
