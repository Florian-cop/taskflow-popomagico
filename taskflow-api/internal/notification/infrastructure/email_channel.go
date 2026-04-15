package infrastructure

import (
	"context"
	"log"

	"taskflow-api/internal/notification/domain"
)

// EmailChannel simule l'envoi d'email via un log console.
// Raccourci assumé : en production, remplacer par un vrai client SMTP / SendGrid / SES.
// Aucun service métier ne connaît cette implémentation — ils voient uniquement l'interface Channel.
type EmailChannel struct{}

func NewEmailChannel() *EmailChannel { return &EmailChannel{} }

func (c *EmailChannel) Name() string { return "email" }

func (c *EmailChannel) Send(_ context.Context, n *domain.Notification) error {
	log.Printf("[EMAIL] to=%s type=%s title=%q body=%q", n.UserID, n.Type, n.Title, n.Body)
	return nil
}
