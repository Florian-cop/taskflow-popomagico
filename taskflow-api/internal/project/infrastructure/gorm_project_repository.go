package infrastructure

import (
	"context"

	"gorm.io/gorm"

	"taskflow-api/internal/project/domain"
	sharedDomain "taskflow-api/internal/shared/domain"
)

type GormProjectRepository struct {
	db *gorm.DB
}

func NewGormProjectRepository(db *gorm.DB) *GormProjectRepository {
	return &GormProjectRepository{db: db}
}

func (r *GormProjectRepository) FindByID(ctx context.Context, id string) (*domain.Project, error) {
	var model ProjectModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedDomain.ErrNotFound
		}
		return nil, err
	}

	var memberModels []MemberModel
	if err := r.db.WithContext(ctx).Where("project_id = ?", id).Find(&memberModels).Error; err != nil {
		return nil, err
	}

	return toDomain(&model, memberModels), nil
}

func (r *GormProjectRepository) FindAll(ctx context.Context) ([]*domain.Project, error) {
	var models []ProjectModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	projects := make([]*domain.Project, len(models))
	for i, model := range models {
		var memberModels []MemberModel
		if err := r.db.WithContext(ctx).Where("project_id = ?", model.ID).Find(&memberModels).Error; err != nil {
			return nil, err
		}
		projects[i] = toDomain(&model, memberModels)
	}

	return projects, nil
}

func (r *GormProjectRepository) Save(ctx context.Context, project *domain.Project) error {
	model, memberModels := toModel(project)

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Upsert du projet
		if err := tx.Save(model).Error; err != nil {
			return err
		}

		// Remplace tous les membres
		if err := tx.Where("project_id = ?", project.ID).Delete(&MemberModel{}).Error; err != nil {
			return err
		}
		if len(memberModels) > 0 {
			if err := tx.Create(&memberModels).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *GormProjectRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("project_id = ?", id).Delete(&MemberModel{}).Error; err != nil {
			return err
		}
		return tx.Delete(&ProjectModel{}, "id = ?", id).Error
	})
}
