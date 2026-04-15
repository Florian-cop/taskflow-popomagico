package infrastructure

import (
	"context"
	"encoding/json"
	"strings"

	auditApp "taskflow-api/internal/audit/application"
	auditDomain "taskflow-api/internal/audit/domain"
	sharedApp "taskflow-api/internal/shared/application"
	sharedDomain "taskflow-api/internal/shared/domain"
)

// EventHandlers enregistre un handler universel sur tous les events métier.
// Consommateur pur : aucun service métier ne connaît son existence.
type EventHandlers struct {
	service *auditApp.AuditService
}

func NewEventHandlers(s *auditApp.AuditService) *EventHandlers {
	return &EventHandlers{service: s}
}

// Register branche le handler sur tous les events que l'on souhaite tracer.
func (h *EventHandlers) Register(bus sharedApp.EventBus, eventNames ...string) {
	for _, name := range eventNames {
		bus.Subscribe(name, h.Handle)
	}
}

func (h *EventHandlers) Handle(ctx context.Context, e sharedDomain.DomainEvent) error {
	userID, _ := sharedDomain.UserIDFromContext(ctx)

	aggregateType := ""
	if parts := strings.SplitN(e.EventName(), ".", 2); len(parts) > 0 {
		aggregateType = parts[0]
	}

	payload, _ := json.Marshal(e)

	return h.service.Record(ctx, &auditDomain.AuditLog{
		ID:            sharedDomain.NewID(),
		UserID:        userID,
		EventName:     e.EventName(),
		AggregateType: aggregateType,
		AggregateID:   e.AggregateID(),
		Payload:       string(payload),
		OccurredAt:    e.OccurredAt(),
	})
}
