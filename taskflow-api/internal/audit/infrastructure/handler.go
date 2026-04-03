package infrastructure

import (
	"context"
	"log"

	auditDomain "taskflow-api/internal/audit/domain"
	projectEvents "taskflow-api/internal/project/domain"
	sharedDomain "taskflow-api/internal/shared/domain"
	taskEvents "taskflow-api/internal/task/domain"
)

// AuditHandler ecoute les domain events et persiste les entrees d'audit.
type AuditHandler struct {
	repo auditDomain.AuditRepository
}

func NewAuditHandler(repo auditDomain.AuditRepository) *AuditHandler {
	return &AuditHandler{repo: repo}
}

// Handle traite un domain event et cree l'entree d'audit correspondante.
func (h *AuditHandler) Handle(ctx context.Context, event sharedDomain.DomainEvent) error {
	entry := h.buildEntry(ctx, event)

	if err := h.repo.Save(ctx, entry); err != nil {
		log.Printf("[Audit] erreur sauvegarde: %v", err)
		return err
	}

	log.Printf("[Audit] %s | %s %s:%s | by %s",
		entry.OccurredAt.Format("2006-01-02 15:04:05"),
		entry.Action, entry.EntityType, entry.EntityID, entry.Actor,
	)
	return nil
}

func (h *AuditHandler) buildEntry(ctx context.Context, event sharedDomain.DomainEvent) *auditDomain.AuditEntry {
	entry := &auditDomain.AuditEntry{
		ID:         sharedDomain.NewID(),
		Action:     event.EventName(),
		EntityID:   event.AggregateID(),
		OccurredAt: event.OccurredAt(),
		Details:    make(map[string]string),
	}

	// Extraire l'acteur et les details selon le type d'evenement
	switch e := event.(type) {
	case taskEvents.TaskCreatedEvent:
		entry.EntityType = "Task"
		entry.Actor = actorFromContext(ctx)
		entry.Details["title"] = e.Title()
		entry.Details["projectId"] = e.ProjectID()

	case taskEvents.TaskMovedEvent:
		entry.EntityType = "Task"
		entry.Actor = actorFromContext(ctx)
		entry.Details["from"] = string(e.FromStatus())
		entry.Details["to"] = string(e.ToStatus())
		entry.Details["projectId"] = e.ProjectID()

	case taskEvents.TaskAssignedEvent:
		entry.EntityType = "Task"
		entry.Actor = actorFromContext(ctx)
		entry.Details["assigneeId"] = e.AssigneeID()
		entry.Details["projectId"] = e.ProjectID()

	case projectEvents.ProjectCreatedEvent:
		entry.EntityType = "Project"
		entry.Actor = e.OwnerID()
		entry.Details["name"] = e.Name()

	case projectEvents.MemberAddedEvent:
		entry.EntityType = "Project"
		entry.Actor = actorFromContext(ctx)
		entry.Details["addedUserId"] = e.UserID()

	default:
		entry.EntityType = "Unknown"
		entry.Actor = actorFromContext(ctx)
	}

	return entry
}

type contextKey string

const UserIDKey contextKey = "userID"

func actorFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		return userID
	}
	return "system"
}
