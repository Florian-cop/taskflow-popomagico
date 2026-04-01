package application

import (
	"context"

	"taskflow-api/internal/project/domain"
	sharedDomain "taskflow-api/internal/shared/domain"
)

// MockProjectRepository est un repository en mémoire pour les tests.
type MockProjectRepository struct {
	Projects map[string]*domain.Project
}

func NewMockProjectRepository() *MockProjectRepository {
	return &MockProjectRepository{Projects: make(map[string]*domain.Project)}
}

func (m *MockProjectRepository) FindByID(ctx context.Context, id string) (*domain.Project, error) {
	project, ok := m.Projects[id]
	if !ok {
		return nil, sharedDomain.ErrNotFound
	}
	return project, nil
}

func (m *MockProjectRepository) FindAll(ctx context.Context) ([]*domain.Project, error) {
	var projects []*domain.Project
	for _, p := range m.Projects {
		projects = append(projects, p)
	}
	return projects, nil
}

func (m *MockProjectRepository) Save(ctx context.Context, project *domain.Project) error {
	m.Projects[project.ID] = project
	return nil
}

func (m *MockProjectRepository) Delete(ctx context.Context, id string) error {
	delete(m.Projects, id)
	return nil
}
