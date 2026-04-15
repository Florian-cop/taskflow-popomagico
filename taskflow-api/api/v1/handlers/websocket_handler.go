package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"

	projectApp "taskflow-api/internal/project/application"
	realtimeInfra "taskflow-api/internal/realtime/infrastructure"
	sharedDomain "taskflow-api/internal/shared/domain"
	userDomain "taskflow-api/internal/user/domain"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WebSocketHandler struct {
	broadcaster    *realtimeInfra.WSBroadcaster
	projectService *projectApp.ProjectService
	tokens         userDomain.TokenGenerator
}

func NewWebSocketHandler(
	b *realtimeInfra.WSBroadcaster,
	p *projectApp.ProjectService,
	t userDomain.TokenGenerator,
) *WebSocketHandler {
	return &WebSocketHandler{broadcaster: b, projectService: p, tokens: t}
}

// HandleConnect GET /api/v1/projects/{id}/ws
// Scoping : l'utilisateur doit être membre du projet.
// Auth : token JWT soit via header Authorization (si le client le permet),
// soit via query param ?token=xxx pour les WebSocket navigateur.
func (h *WebSocketHandler) HandleConnect(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	userID, err := h.authenticate(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	project, err := h.projectService.GetProject(r.Context(), projectID)
	if err != nil {
		http.Error(w, "project not found", http.StatusNotFound)
		return
	}

	member := false
	for _, m := range project.Members {
		if m.UserID == userID {
			member = true
			break
		}
	}
	if !member {
		http.Error(w, "forbidden: not a member of this project", http.StatusForbidden)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &realtimeInfra.Client{
		ID:        sharedDomain.NewID(),
		ProjectID: projectID,
		Conn:      conn,
		Send:      make(chan []byte, 16),
	}
	h.broadcaster.Register(client)

	go h.writePump(client)
	h.readPump(client)
}

func (h *WebSocketHandler) authenticate(r *http.Request) (string, error) {
	token := r.URL.Query().Get("token")
	if token == "" {
		if v, ok := sharedDomain.UserIDFromContext(r.Context()); ok {
			return v, nil
		}
		return "", websocket.ErrBadHandshake
	}
	user, err := h.tokens.Validate(r.Context(), token)
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

func (h *WebSocketHandler) writePump(c *realtimeInfra.Client) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (h *WebSocketHandler) readPump(c *realtimeInfra.Client) {
	defer h.broadcaster.Unregister(c.ProjectID, c.ID)
	c.Conn.SetReadLimit(512)
	c.Conn.SetPongHandler(func(string) error { return nil })
	for {
		if _, _, err := c.Conn.ReadMessage(); err != nil {
			return
		}
	}
}
