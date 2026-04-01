package application

import (
	"context"

	"taskflow-api/internal/shared/domain"
)

// MockEventBus enregistre les events publiés pour les vérifier dans les tests.
type MockEventBus struct {
	Published []domain.DomainEvent
}

func NewMockEventBus() *MockEventBus {
	return &MockEventBus{}
}

func (m *MockEventBus) Publish(ctx context.Context, event domain.DomainEvent) error {
	m.Published = append(m.Published, event)
	return nil
}

func (m *MockEventBus) Subscribe(eventName string, handler EventHandler) {}
