package infrastructure

import (
	"context"

	"gorm.io/gorm"

	"taskflow-api/internal/audit/domain"
)

type GormAuditRepository struct {
	db *gorm.DB
}

func NewGormAuditRepository(db *gorm.DB) *GormAuditRepository {
	return &GormAuditRepository{db: db}
}

func (r *GormAuditRepository) Save(ctx context.Context, entry *domain.AuditLog) error {
	return r.db.WithContext(ctx).Create(&AuditLogModel{
		ID: entry.ID, UserID: entry.UserID, EventName: entry.EventName,
		AggregateType: entry.AggregateType, AggregateID: entry.AggregateID,
		Payload: entry.Payload, OccurredAt: entry.OccurredAt,
	}).Error
}

func (r *GormAuditRepository) Query(ctx context.Context, f domain.Filter) ([]*domain.AuditLog, error) {
	q := r.db.WithContext(ctx).Order("occurred_at DESC")
	if f.AggregateType != "" {
		q = q.Where("aggregate_type = ?", f.AggregateType)
	}
	if f.AggregateID != "" {
		q = q.Where("aggregate_id = ?", f.AggregateID)
	}
	if f.UserID != "" {
		q = q.Where("user_id = ?", f.UserID)
	}
	if f.Limit > 0 {
		q = q.Limit(f.Limit)
	} else {
		q = q.Limit(100)
	}

	var models []AuditLogModel
	if err := q.Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]*domain.AuditLog, len(models))
	for i, m := range models {
		out[i] = &domain.AuditLog{
			ID: m.ID, UserID: m.UserID, EventName: m.EventName,
			AggregateType: m.AggregateType, AggregateID: m.AggregateID,
			Payload: m.Payload, OccurredAt: m.OccurredAt,
		}
	}
	return out, nil
}
