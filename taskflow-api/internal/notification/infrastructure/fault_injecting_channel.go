package infrastructure

import (
	"context"
	"errors"
	"sync/atomic"

	"taskflow-api/internal/notification/domain"
)

// FaultInjectingChannel est un décorateur qui peut faire échouer le canal sous-jacent à la demande.
// Sert à démontrer la résilience du Dispatcher : un canal en panne ne doit pas affecter les autres.
//
// Pattern Décorateur — l'extérieur ne voit qu'une `Channel`, le wrapping reste invisible
// pour le Dispatcher qui n'est pas modifié.
type FaultInjectingChannel struct {
	inner   domain.Channel
	failing atomic.Bool
}

func NewFaultInjectingChannel(inner domain.Channel) *FaultInjectingChannel {
	return &FaultInjectingChannel{inner: inner}
}

func (c *FaultInjectingChannel) Name() string {
	return c.inner.Name()
}

func (c *FaultInjectingChannel) Send(ctx context.Context, n *domain.Notification) error {
	if c.failing.Load() {
		return errors.New("simulated channel failure")
	}
	return c.inner.Send(ctx, n)
}

func (c *FaultInjectingChannel) SetFailing(failing bool) {
	c.failing.Store(failing)
}

func (c *FaultInjectingChannel) IsFailing() bool {
	return c.failing.Load()
}
