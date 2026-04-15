package dto

import "time"

type AuditLogResponse struct {
	ID            string    `json:"id"`
	UserID        string    `json:"userId"`
	EventName     string    `json:"eventName"`
	AggregateType string    `json:"aggregateType"`
	AggregateID   string    `json:"aggregateId"`
	Payload       string    `json:"payload"`
	OccurredAt    time.Time `json:"occurredAt"`
}
