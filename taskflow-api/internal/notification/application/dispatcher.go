package application

import (
	"context"
	"log"

	notifDomain "taskflow-api/internal/notification/domain"
	sharedDomain "taskflow-api/internal/shared/domain"
)

// Dispatcher route une notification vers les canaux activés par l'utilisateur.
// Strategy pattern : []Channel est une slice d'interfaces ; ajouter un canal
// ne nécessite aucune modification du code existant.
type Dispatcher struct {
	channels  []notifDomain.Channel
	prefsRepo notifDomain.PreferencesRepository
}

func NewDispatcher(channels []notifDomain.Channel, prefsRepo notifDomain.PreferencesRepository) *Dispatcher {
	return &Dispatcher{channels: channels, prefsRepo: prefsRepo}
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
		}
	}
	return nil
}
