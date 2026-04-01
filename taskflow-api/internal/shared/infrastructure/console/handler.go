package console

import (
	"context"
	"log"

	"taskflow-api/internal/shared/domain"
)

// Handle affiche dans la console le nom de l'event, l'ID de l'agrégat et l'horodatage.
func Handle(ctx context.Context, event domain.DomainEvent) error {
	log.Printf("[Event] %s | ID: %s | At: %s",
		event.EventName(),
		event.AggregateID(),
		event.OccurredAt().Format("2006-01-02 15:04:05"),
	)
	return nil
}
