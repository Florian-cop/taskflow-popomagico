package domain

import "time"

// AuditEntry represente une entree dans le journal d'audit.
type AuditEntry struct {
	ID         string
	Actor      string
	Action     string
	EntityType string
	EntityID   string
	Details    map[string]string
	OccurredAt time.Time
}
