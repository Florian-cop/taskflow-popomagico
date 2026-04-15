package infrastructure

import (
	"context"

	"taskflow-api/internal/notification/domain"
)

// InAppChannel persiste la notification pour consultation via l'API.
type InAppChannel struct {
	repo domain.NotificationRepository
}

func NewInAppChannel(repo domain.NotificationRepository) *InAppChannel {
	return &InAppChannel{repo: repo}
}

func (c *InAppChannel) Name() string { return "in_app" }

func (c *InAppChannel) Send(ctx context.Context, n *domain.Notification) error {
	return c.repo.Save(ctx, n)
}
