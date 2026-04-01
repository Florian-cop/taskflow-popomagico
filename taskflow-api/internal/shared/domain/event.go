package domain

import "time"

// DomainEvent est le contrat commun à tous les événements métier.
type DomainEvent interface {
	EventID() string
	EventName() string
	AggregateID() string
	OccurredAt() time.Time
}
