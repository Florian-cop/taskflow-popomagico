package domain

import "context"

type Filter struct {
	AggregateType string
	AggregateID   string
	UserID        string
	Limit         int
}

type AuditRepository interface {
	Save(ctx context.Context, entry *AuditLog) error
	Query(ctx context.Context, f Filter) ([]*AuditLog, error)
}
