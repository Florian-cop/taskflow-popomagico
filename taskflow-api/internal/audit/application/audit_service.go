package application

import (
	"context"

	auditDomain "taskflow-api/internal/audit/domain"
)

type AuditService struct {
	repo auditDomain.AuditRepository
}

func NewAuditService(repo auditDomain.AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) Record(ctx context.Context, entry *auditDomain.AuditLog) error {
	return s.repo.Save(ctx, entry)
}

func (s *AuditService) Query(ctx context.Context, q QueryDTO) ([]*AuditLogDTO, error) {
	entries, err := s.repo.Query(ctx, auditDomain.Filter{
		AggregateType: q.AggregateType,
		AggregateID:   q.AggregateID,
		UserID:        q.UserID,
		Limit:         q.Limit,
	})
	if err != nil {
		return nil, err
	}
	out := make([]*AuditLogDTO, len(entries))
	for i, e := range entries {
		out[i] = &AuditLogDTO{
			ID: e.ID, UserID: e.UserID, EventName: e.EventName,
			AggregateType: e.AggregateType, AggregateID: e.AggregateID,
			Payload: e.Payload, OccurredAt: e.OccurredAt,
		}
	}
	return out, nil
}
