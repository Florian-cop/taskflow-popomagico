package domain

import "context"

// Channel est le port sortant pour un canal de notification.
// Ajouter un canal (Slack, Teams, SMS) = implémenter cette interface
// sans rien modifier aux services métier (task, project).
type Channel interface {
	Name() string
	Send(ctx context.Context, n *Notification) error
}
