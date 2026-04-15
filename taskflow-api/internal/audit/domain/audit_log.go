package domain

import "time"

type AuditLog struct {
	ID            string
	UserID        string
	EventName     string
	AggregateType string
	AggregateID   string
	Payload       string // JSON
	OccurredAt    time.Time
}
