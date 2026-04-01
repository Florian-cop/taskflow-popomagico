package domain

import "context"

// TaskRepository est le port sortant pour la persistance des tâches.
// L'implémentation concrète (GORM, mémoire, etc.) vit dans infrastructure/.
type TaskRepository interface {
	FindByID(ctx context.Context, id string) (*Task, error)
	FindByProjectID(ctx context.Context, projectID string) ([]*Task, error)
	Save(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id string) error
}
