package application

import "time"

type AuditLogDTO struct {
	ID            string
	UserID        string
	EventName     string
	AggregateType string
	AggregateID   string
	Payload       string
	OccurredAt    time.Time
}

type QueryDTO struct {
	AggregateType string
	AggregateID   string
	UserID        string
	Limit         int
}
