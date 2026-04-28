package domain

import "time"

type FailedStatus string

const (
	FailedStatusPending FailedStatus = "pending"
	FailedStatusRetried FailedStatus = "retried"
)

// FailedNotification représente une tentative d'envoi qui a échoué.
// Persistée pour retraitement manuel (chantier 1 disruption #2).
type FailedNotification struct {
	ID            string
	NotificationID string
	UserID        string
	Channel       string
	Type          string
	Title         string
	Body          string
	Error         string
	RetryCount    int
	Status        FailedStatus
	OccurredAt    time.Time
	LastRetriedAt *time.Time
}
