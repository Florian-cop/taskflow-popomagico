package infrastructure

import (
	"context"

	"gorm.io/gorm"

	"taskflow-api/internal/notification/domain"
	sharedDomain "taskflow-api/internal/shared/domain"
)

type GormNotificationRepository struct {
	db *gorm.DB
}

func NewGormNotificationRepository(db *gorm.DB) *GormNotificationRepository {
	return &GormNotificationRepository{db: db}
}

func (r *GormNotificationRepository) Save(ctx context.Context, n *domain.Notification) error {
	return r.db.WithContext(ctx).Create(toModel(n)).Error
}

func (r *GormNotificationRepository) Update(ctx context.Context, n *domain.Notification) error {
	return r.db.WithContext(ctx).Save(toModel(n)).Error
}

func (r *GormNotificationRepository) FindByID(ctx context.Context, id string) (*domain.Notification, error) {
	var m NotificationModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedDomain.ErrNotFound
		}
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *GormNotificationRepository) ListByUser(ctx context.Context, userID string) ([]*domain.Notification, error) {
	var models []NotificationModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]*domain.Notification, len(models))
	for i, m := range models {
		out[i] = toDomain(&m)
	}
	return out, nil
}

func toModel(n *domain.Notification) *NotificationModel {
	return &NotificationModel{
		ID: n.ID, UserID: n.UserID, Type: n.Type,
		Title: n.Title, Body: n.Body,
		ReadAt: n.ReadAt, CreatedAt: n.CreatedAt,
	}
}

func toDomain(m *NotificationModel) *domain.Notification {
	return &domain.Notification{
		ID: m.ID, UserID: m.UserID, Type: m.Type,
		Title: m.Title, Body: m.Body,
		ReadAt: m.ReadAt, CreatedAt: m.CreatedAt,
	}
}
