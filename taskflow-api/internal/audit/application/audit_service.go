package application

import (
	"context"
	"time"

	auditDomain "taskflow-api/internal/audit/domain"
)

type AuditService struct {
	repo auditDomain.AuditRepository
}

func NewAuditService(repo auditDomain.AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

type AuditEntryDTO struct {
	ID         string            `json:"id"`
	Actor      string            `json:"actor"`
	Action     string            `json:"action"`
	EntityType string            `json:"entityType"`
	EntityID   string            `json:"entityId"`
	Details    map[string]string `json:"details"`
	OccurredAt time.Time         `json:"occurredAt"`
}

func (s *AuditService) GetByEntityID(ctx context.Context, entityID string) ([]AuditEntryDTO, error) {
	entries, err := s.repo.FindByEntityID(ctx, entityID)
	if err != nil {
		return nil, err
	}
	return toDTOs(entries), nil
}

func (s *AuditService) GetAll(ctx context.Context) ([]AuditEntryDTO, error) {
	entries, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return toDTOs(entries), nil
}

func toDTOs(entries []*auditDomain.AuditEntry) []AuditEntryDTO {
	dtos := make([]AuditEntryDTO, len(entries))
	for i, e := range entries {
		dtos[i] = AuditEntryDTO{
			ID:         e.ID,
			Actor:      e.Actor,
			Action:     e.Action,
			EntityType: e.EntityType,
			EntityID:   e.EntityID,
			Details:    e.Details,
			OccurredAt: e.OccurredAt,
		}
	}
	return dtos
}
