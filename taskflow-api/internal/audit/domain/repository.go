package domain

import "context"

// AuditRepository est le port de sortie pour le stockage de l'audit trail.
// Changer de backend de stockage = implementer cette interface.
type AuditRepository interface {
	Save(ctx context.Context, entry *AuditEntry) error
	FindByEntityID(ctx context.Context, entityID string) ([]*AuditEntry, error)
	FindAll(ctx context.Context) ([]*AuditEntry, error)
}
