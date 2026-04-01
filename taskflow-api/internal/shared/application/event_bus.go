package application

import (
	"context"

	"taskflow-api/internal/shared/domain"
)

// EventHandler est une fonction qui réagit à un événement métier.
type EventHandler func(ctx context.Context, event domain.DomainEvent) error

// EventBus est le port de publication et d'abonnement aux événements métier.
type EventBus interface {
	Publish(ctx context.Context, event domain.DomainEvent) error
	Subscribe(eventName string, handler EventHandler)
}
