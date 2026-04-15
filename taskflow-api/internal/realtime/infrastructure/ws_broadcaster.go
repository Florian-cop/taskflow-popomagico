package infrastructure

import (
	"context"
	"sync"

	"github.com/gorilla/websocket"
)

// Client représente une connexion WebSocket inscrite à un projet.
type Client struct {
	ID        string
	ProjectID string
	Conn      *websocket.Conn
	Send      chan []byte
}

// WSBroadcaster maintient les clients groupés par projet et diffuse les messages.
// sync.Map est retenu pour sa simplicité sur un cas read-heavy + writes occasionnels.
type WSBroadcaster struct {
	// projectID -> *projectRoom
	rooms sync.Map
}

type projectRoom struct {
	mu      sync.RWMutex
	clients map[string]*Client
}

func NewWSBroadcaster() *WSBroadcaster {
	return &WSBroadcaster{}
}

func (b *WSBroadcaster) Register(client *Client) {
	raw, _ := b.rooms.LoadOrStore(client.ProjectID, &projectRoom{clients: map[string]*Client{}})
	room := raw.(*projectRoom)
	room.mu.Lock()
	room.clients[client.ID] = client
	room.mu.Unlock()
}

func (b *WSBroadcaster) Unregister(projectID, clientID string) {
	raw, ok := b.rooms.Load(projectID)
	if !ok {
		return
	}
	room := raw.(*projectRoom)
	room.mu.Lock()
	if c, exists := room.clients[clientID]; exists {
		close(c.Send)
		delete(room.clients, clientID)
	}
	room.mu.Unlock()
}

// Broadcast implémente realtime/domain.Broadcaster.
func (b *WSBroadcaster) Broadcast(_ context.Context, projectID string, message []byte) error {
	raw, ok := b.rooms.Load(projectID)
	if !ok {
		return nil
	}
	room := raw.(*projectRoom)
	room.mu.RLock()
	defer room.mu.RUnlock()
	for _, c := range room.clients {
		select {
		case c.Send <- message:
		default:
			// si le buffer est plein, on drop — le client est lent
		}
	}
	return nil
}
