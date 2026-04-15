package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"taskflow-api/api/v1/dto"
	notifApp "taskflow-api/internal/notification/application"
	sharedDomain "taskflow-api/internal/shared/domain"
)

type NotificationHandler struct {
	service *notifApp.NotificationService
}

func NewNotificationHandler(s *notifApp.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: s}
}

// GET /api/v1/notifications
func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := sharedDomain.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	items, err := h.service.ListByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]dto.NotificationResponse, len(items))
	for i, n := range items {
		resp[i] = dto.NotificationResponse{
			ID: n.ID, Type: n.Type, Title: n.Title, Body: n.Body,
			ReadAt: n.ReadAt, CreatedAt: n.CreatedAt,
		}
	}
	writeJSON(w, http.StatusOK, resp)
}

// PATCH /api/v1/notifications/{id}/read
func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	userID, ok := sharedDomain.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	id := chi.URLParam(r, "id")
	n, err := h.service.MarkAsRead(r.Context(), id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, dto.NotificationResponse{
		ID: n.ID, Type: n.Type, Title: n.Title, Body: n.Body,
		ReadAt: n.ReadAt, CreatedAt: n.CreatedAt,
	})
}

// GET /api/v1/notifications/preferences
func (h *NotificationHandler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	userID, ok := sharedDomain.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	p, err := h.service.GetPreferences(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, dto.PreferencesResponse{Enabled: p.Enabled})
}

// PUT /api/v1/notifications/preferences
func (h *NotificationHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	userID, ok := sharedDomain.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req dto.UpdatePreferencesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	p, err := h.service.UpdatePreferences(r.Context(), notifApp.UpdatePreferencesDTO{
		UserID: userID, Enabled: req.Enabled,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, dto.PreferencesResponse{Enabled: p.Enabled})
}
