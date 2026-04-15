package infrastructure

import (
	"context"
	"fmt"

	notifApp "taskflow-api/internal/notification/application"
	notifDomain "taskflow-api/internal/notification/domain"
	sharedApp "taskflow-api/internal/shared/application"
	sharedDomain "taskflow-api/internal/shared/domain"
	taskDomain "taskflow-api/internal/task/domain"
)

// EventHandlers branche la notification au bus d'events.
// Consommateur pur : les services task/project ignorent son existence.
type EventHandlers struct {
	dispatcher *notifApp.Dispatcher
	members    notifDomain.MemberFinder
}

func NewEventHandlers(d *notifApp.Dispatcher, m notifDomain.MemberFinder) *EventHandlers {
	return &EventHandlers{dispatcher: d, members: m}
}

// Register abonne les handlers aux events métier.
func (h *EventHandlers) Register(bus sharedApp.EventBus) {
	bus.Subscribe("task.assigned", h.HandleTaskAssigned)
	bus.Subscribe("task.moved", h.HandleTaskMoved)
}

func (h *EventHandlers) HandleTaskAssigned(ctx context.Context, e sharedDomain.DomainEvent) error {
	event, ok := e.(taskDomain.TaskAssignedEvent)
	if !ok {
		return nil
	}
	title := "Une tâche vous a été assignée"
	body := fmt.Sprintf("Tâche %s sur le projet %s", event.AggregateID(), event.ProjectID())
	return h.dispatcher.Dispatch(ctx, event.AssigneeID(), "task.assigned", title, body)
}

func (h *EventHandlers) HandleTaskMoved(ctx context.Context, e sharedDomain.DomainEvent) error {
	event, ok := e.(taskDomain.TaskMovedEvent)
	if !ok {
		return nil
	}
	members, err := h.members.FindMembers(ctx, event.ProjectID())
	if err != nil {
		return err
	}
	actorID, _ := sharedDomain.UserIDFromContext(ctx)
	title := "Une tâche a été déplacée"
	body := fmt.Sprintf("%s → %s", event.FromStatus(), event.ToStatus())
	for _, userID := range members {
		if userID == actorID {
			continue
		}
		_ = h.dispatcher.Dispatch(ctx, userID, "task.moved", title, body)
	}
	return nil
}
