package infrastructure

import (
	"context"

	"gorm.io/gorm"

	sharedDomain "taskflow-api/internal/shared/domain"
	"taskflow-api/internal/user/domain"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedDomain.ErrNotFound
		}
		return nil, err
	}
	return toDomain(&model), nil
}

func (r *GormUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).First(&model, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedDomain.ErrNotFound
		}
		return nil, err
	}
	return toDomain(&model), nil
}

func (r *GormUserRepository) SearchByEmail(ctx context.Context, query string, limit int) ([]*domain.User, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	pattern := "%" + query + "%"

	var models []UserModel
	if err := r.db.WithContext(ctx).
		Where("email ILIKE ?", pattern).
		Order("email ASC").
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}

	out := make([]*domain.User, len(models))
	for i, m := range models {
		out[i] = toDomain(&m)
	}
	return out, nil
}

func (r *GormUserRepository) Save(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(toModel(user)).Error
}

func (r *GormUserRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(toModel(user)).Error
}
