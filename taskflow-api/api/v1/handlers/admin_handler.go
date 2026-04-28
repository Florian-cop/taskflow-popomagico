package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	notifApp "taskflow-api/internal/notification/application"
)

type AdminHandler struct {
	service *notifApp.AdminService
}

func NewAdminHandler(s *notifApp.AdminService) *AdminHandler {
	return &AdminHandler{service: s}
}

// GET /api/v1/admin/notifications/failed?limit=100
func (h *AdminHandler) ListFailed(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	items, err := h.service.ListFailed(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

// POST /api/v1/admin/notifications/failed/{id}/retry
func (h *AdminHandler) RetryFailed(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.RetryFailed(r.Context(), id); err != nil {
		if errors.Is(err, notifApp.ErrChannelUnknown) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GET /api/v1/admin/notifications/channels
func (h *AdminHandler) ListChannels(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.service.ListChannels())
}

// PUT /api/v1/admin/notifications/channels/{name}
// body: { "failing": true|false }
func (h *AdminHandler) SetChannelFailing(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	var req struct {
		Failing bool `json:"failing"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if err := h.service.SetChannelFailing(name, req.Failing); err != nil {
		if errors.Is(err, notifApp.ErrChannelUnknown) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
