package infrastructure

import (
	"context"

	"gorm.io/gorm"

	sharedDomain "taskflow-api/internal/shared/domain"
	"taskflow-api/internal/task/domain"
)

type GormTaskRepository struct {
	db *gorm.DB
}

func NewGormTaskRepository(db *gorm.DB) *GormTaskRepository {
	return &GormTaskRepository{db: db}
}

func (r *GormTaskRepository) FindByID(ctx context.Context, id string) (*domain.Task, error) {
	var model TaskModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedDomain.ErrNotFound
		}
		return nil, err
	}
	return toDomain(&model), nil
}

func (r *GormTaskRepository) FindByProjectID(ctx context.Context, projectID string) ([]*domain.Task, error) {
	var models []TaskModel
	if err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Find(&models).Error; err != nil {
		return nil, err
	}

	tasks := make([]*domain.Task, len(models))
	for i, model := range models {
		tasks[i] = toDomain(&model)
	}
	return tasks, nil
}

func (r *GormTaskRepository) Save(ctx context.Context, task *domain.Task) error {
	model := toModel(task)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *GormTaskRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&TaskModel{}, "id = ?", id).Error
}
