package infrastructure

import (
	"context"
	"time"

	"gorm.io/gorm"

	"taskflow-api/internal/notification/domain"
	sharedDomain "taskflow-api/internal/shared/domain"
)

type FailedNotificationModel struct {
	ID             string `gorm:"primaryKey"`
	NotificationID string
	UserID         string `gorm:"index"`
	Channel        string `gorm:"index"`
	Type           string
	Title          string
	Body           string
	Error          string
	RetryCount     int
	Status         string `gorm:"index"`
	OccurredAt     time.Time
	LastRetriedAt  *time.Time
}

func (FailedNotificationModel) TableName() string { return "failed_notifications" }

type GormFailedNotificationRepository struct {
	db *gorm.DB
}

func NewGormFailedNotificationRepository(db *gorm.DB) *GormFailedNotificationRepository {
	return &GormFailedNotificationRepository{db: db}
}

func (r *GormFailedNotificationRepository) Save(ctx context.Context, f *domain.FailedNotification) error {
	return r.db.WithContext(ctx).Create(toFailedModel(f)).Error
}

func (r *GormFailedNotificationRepository) Update(ctx context.Context, f *domain.FailedNotification) error {
	return r.db.WithContext(ctx).Save(toFailedModel(f)).Error
}

func (r *GormFailedNotificationRepository) FindByID(ctx context.Context, id string) (*domain.FailedNotification, error) {
	var m FailedNotificationModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, sharedDomain.ErrNotFound
		}
		return nil, err
	}
	return toFailedDomain(&m), nil
}

func (r *GormFailedNotificationRepository) ListPending(ctx context.Context, limit int) ([]*domain.FailedNotification, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	var models []FailedNotificationModel
	if err := r.db.WithContext(ctx).
		Where("status = ?", string(domain.FailedStatusPending)).
		Order("occurred_at DESC").
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]*domain.FailedNotification, len(models))
	for i, m := range models {
		out[i] = toFailedDomain(&m)
	}
	return out, nil
}

func toFailedModel(f *domain.FailedNotification) *FailedNotificationModel {
	return &FailedNotificationModel{
		ID: f.ID, NotificationID: f.NotificationID, UserID: f.UserID,
		Channel: f.Channel, Type: f.Type, Title: f.Title, Body: f.Body,
		Error: f.Error, RetryCount: f.RetryCount, Status: string(f.Status),
		OccurredAt: f.OccurredAt, LastRetriedAt: f.LastRetriedAt,
	}
}

func toFailedDomain(m *FailedNotificationModel) *domain.FailedNotification {
	return &domain.FailedNotification{
		ID: m.ID, NotificationID: m.NotificationID, UserID: m.UserID,
		Channel: m.Channel, Type: m.Type, Title: m.Title, Body: m.Body,
		Error: m.Error, RetryCount: m.RetryCount, Status: domain.FailedStatus(m.Status),
		OccurredAt: m.OccurredAt, LastRetriedAt: m.LastRetriedAt,
	}
}
