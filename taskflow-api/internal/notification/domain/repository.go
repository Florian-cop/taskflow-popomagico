package domain

import "context"

type NotificationRepository interface {
	Save(ctx context.Context, n *Notification) error
	Update(ctx context.Context, n *Notification) error
	FindByID(ctx context.Context, id string) (*Notification, error)
	ListByUser(ctx context.Context, userID string) ([]*Notification, error)
}

type PreferencesRepository interface {
	Get(ctx context.Context, userID string) (*Preferences, error)
	Save(ctx context.Context, p *Preferences) error
}

// MemberFinder abstrait la lecture des membres d'un projet.
// Permet au handler notification de ne pas dépendre du bounded context project.
type MemberFinder interface {
	FindMembers(ctx context.Context, projectID string) ([]string, error)
}
