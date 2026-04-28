package application

import (
	"context"
	"log"
	"time"

	notifDomain "taskflow-api/internal/notification/domain"
	sharedDomain "taskflow-api/internal/shared/domain"
)

// Dispatcher route une notification vers les canaux activés par l'utilisateur.
// Strategy pattern : []Channel est une slice d'interfaces ; ajouter un canal
// ne nécessite aucune modification du code existant.
//
// Résilience (chantier 1 disruption #2) :
//   - tout échec d'un canal est isolé : les autres canaux continuent ;
//   - les échecs sont persistés via FailedNotificationRepository pour retraitement ;
//   - la méthode Dispatch ne retourne jamais une erreur "canal" — seuls les
//     incidents structurels (impossible de lire les prefs, etc.) remontent en log.
type Dispatcher struct {
	channels   []notifDomain.Channel
	prefsRepo  notifDomain.PreferencesRepository
	failedRepo notifDomain.FailedNotificationRepository
}

func NewDispatcher(
	channels []notifDomain.Channel,
	prefsRepo notifDomain.PreferencesRepository,
	failedRepo notifDomain.FailedNotificationRepository,
) *Dispatcher {
	return &Dispatcher{channels: channels, prefsRepo: prefsRepo, failedRepo: failedRepo}
}

func (d *Dispatcher) Dispatch(ctx context.Context, userID, notifType, title, body string) error {
	prefs, err := d.prefsRepo.Get(ctx, userID)
	if err != nil {
		prefs = notifDomain.DefaultPreferences(userID)
	}

	n := notifDomain.NewNotification(sharedDomain.NewID(), userID, notifType, title, body)

	for _, c := range d.channels {
		if !prefs.IsEnabled(c.Name()) {
			continue
		}
		if err := c.Send(ctx, n); err != nil {
			log.Printf("[Dispatcher] canal %s a échoué pour %s: %v", c.Name(), userID, err)
			d.recordFailure(ctx, c.Name(), n, err)
		}
	}
	return nil
}

func (d *Dispatcher) recordFailure(ctx context.Context, channel string, n *notifDomain.Notification, sendErr error) {
	if d.failedRepo == nil {
		return
	}
	failed := &notifDomain.FailedNotification{
		ID:             sharedDomain.NewID(),
		NotificationID: n.ID,
		UserID:         n.UserID,
		Channel:        channel,
		Type:           n.Type,
		Title:          n.Title,
		Body:           n.Body,
		Error:          sendErr.Error(),
		Status:         notifDomain.FailedStatusPending,
		OccurredAt:     time.Now(),
	}
	if err := d.failedRepo.Save(ctx, failed); err != nil {
		log.Printf("[Dispatcher] impossible de persister l'échec %s: %v", channel, err)
	}
}

// Channels expose la liste des canaux configurés (utile pour les routes admin).
func (d *Dispatcher) Channels() []notifDomain.Channel {
	return d.channels
}

// RetryFailed retente l'envoi d'un message échoué via son canal d'origine.
// Approche simple et synchrone : adaptée au pilote, à remplacer par un worker
// asynchrone avec back-off exponentiel en production.
func (d *Dispatcher) RetryFailed(ctx context.Context, failedID string) error {
	failed, err := d.failedRepo.FindByID(ctx, failedID)
	if err != nil {
		return err
	}
	channel := d.findChannel(failed.Channel)
	if channel == nil {
		return ErrChannelUnknown
	}

	n := &notifDomain.Notification{
		ID:        failed.NotificationID,
		UserID:    failed.UserID,
		Type:      failed.Type,
		Title:     failed.Title,
		Body:      failed.Body,
		CreatedAt: failed.OccurredAt,
	}

	failed.RetryCount++
	now := time.Now()
	failed.LastRetriedAt = &now

	if sendErr := channel.Send(ctx, n); sendErr != nil {
		failed.Error = sendErr.Error()
		_ = d.failedRepo.Update(ctx, failed)
		return sendErr
	}

	failed.Status = notifDomain.FailedStatusRetried
	return d.failedRepo.Update(ctx, failed)
}

func (d *Dispatcher) findChannel(name string) notifDomain.Channel {
	for _, c := range d.channels {
		if c.Name() == name {
			return c
		}
	}
	return nil
}
