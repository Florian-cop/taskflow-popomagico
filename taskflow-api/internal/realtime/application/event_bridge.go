package application

import (
	"context"
	"encoding/json"
	"log"

	realtimeDomain "taskflow-api/internal/realtime/domain"
	sharedApp "taskflow-api/internal/shared/application"
	sharedDomain "taskflow-api/internal/shared/domain"
)

// ProjectAware est satisfait par tout event qui expose un ProjectID.
// Tous les events métier concernés (task.*, project.*, member.added) l'implémentent.
type ProjectAware interface {
	ProjectID() string
}

// EventBridge consomme les events domain et les relaye vers le Broadcaster,
// scopés par projet. C'est un adaptateur sortant : aucun service métier
// ne le connaît, il s'abonne simplement au bus.
type EventBridge struct {
	broadcaster realtimeDomain.Broadcaster
}

func NewEventBridge(b realtimeDomain.Broadcaster) *EventBridge {
	return &EventBridge{broadcaster: b}
}

// Register abonne le bridge aux events configurés.
func (eb *EventBridge) Register(bus sharedApp.EventBus, eventNames ...string) {
	for _, name := range eventNames {
		bus.Subscribe(name, eb.Handle)
	}
}

func (eb *EventBridge) Handle(ctx context.Context, event sharedDomain.DomainEvent) error {
	pa, ok := event.(ProjectAware)
	if !ok {
		return nil
	}

	payload := map[string]any{
		"type":        event.EventName(),
		"aggregateId": event.AggregateID(),
		"projectId":   pa.ProjectID(),
		"occurredAt":  event.OccurredAt(),
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[EventBridge] erreur marshal: %v", err)
		return err
	}

	return eb.broadcaster.Broadcast(ctx, pa.ProjectID(), encoded)
}
