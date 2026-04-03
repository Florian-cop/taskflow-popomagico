package infrastructure

import (
	"context"

	auditDomain "taskflow-api/internal/audit/domain"
	"gorm.io/gorm"
)

type GormAuditRepository struct {
	db *gorm.DB
}

func NewGormAuditRepository(db *gorm.DB) *GormAuditRepository {
	return &GormAuditRepository{db: db}
}

func (r *GormAuditRepository) Save(ctx context.Context, entry *auditDomain.AuditEntry) error {
	model := toModel(entry)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *GormAuditRepository) FindByEntityID(ctx context.Context, entityID string) ([]*auditDomain.AuditEntry, error) {
	var models []AuditEntryModel
	if err := r.db.WithContext(ctx).Where("entity_id = ?", entityID).Order("occurred_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	entries := make([]*auditDomain.AuditEntry, len(models))
	for i := range models {
		entries[i] = toDomain(&models[i])
	}
	return entries, nil
}

func (r *GormAuditRepository) FindAll(ctx context.Context) ([]*auditDomain.AuditEntry, error) {
	var models []AuditEntryModel
	if err := r.db.WithContext(ctx).Order("occurred_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	entries := make([]*auditDomain.AuditEntry, len(models))
	for i := range models {
		entries[i] = toDomain(&models[i])
	}
	return entries, nil
}
