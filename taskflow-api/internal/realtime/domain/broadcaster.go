package domain

import "context"

// Broadcaster est le port sortant pour la diffusion temps réel.
// Implémentations possibles : WebSocket, SSE, Centrifugo, etc.
// La logique métier ne dépend que de cette interface.
type Broadcaster interface {
	Broadcast(ctx context.Context, projectID string, message []byte) error
}
