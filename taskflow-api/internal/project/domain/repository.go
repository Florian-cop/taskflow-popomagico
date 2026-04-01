package domain

import "context"

// ProjectRepository est le port sortant pour la persistance des projets.
type ProjectRepository interface {
	FindByID(ctx context.Context, id string) (*Project, error)
	FindAll(ctx context.Context) ([]*Project, error)
	Save(ctx context.Context, project *Project) error
	Delete(ctx context.Context, id string) error
}
