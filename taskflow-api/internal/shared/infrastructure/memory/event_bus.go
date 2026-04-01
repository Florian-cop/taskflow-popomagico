package memory

import (
	"context"
	"log"

	"taskflow-api/internal/shared/application"
	"taskflow-api/internal/shared/domain"
)

// InMemoryEventBus est un event bus synchrone en mémoire.
// Il implémente l'interface application.EventBus.
// En Phase 1, c'est suffisant. Plus tard on pourra le remplacer par NATS
// sans toucher aux services.
type InMemoryEventBus struct {
	handlers map[string][]application.EventHandler
}

func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		handlers: make(map[string][]application.EventHandler),
	}
}

// Subscribe enregistre un handler pour un type d'événement donné.
func (b *InMemoryEventBus) Subscribe(eventName string, handler application.EventHandler) {
	b.handlers[eventName] = append(b.handlers[eventName], handler)
}

// Publish envoie l'événement à tous les handlers abonnés.
// Si un handler échoue, on log l'erreur mais on continue les autres.
func (b *InMemoryEventBus) Publish(ctx context.Context, event domain.DomainEvent) error {
	for _, handler := range b.handlers[event.EventName()] {
		if err := handler(ctx, event); err != nil {
			log.Printf("[EventBus] erreur handler pour %s: %v", event.EventName(), err)
		}
	}
	return nil
}
