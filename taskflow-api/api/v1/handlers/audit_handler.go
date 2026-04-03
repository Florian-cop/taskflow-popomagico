package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"taskflow-api/internal/audit/application"
)

type AuditHandler struct {
	service *application.AuditService
}

func NewAuditHandler(service *application.AuditService) *AuditHandler {
	return &AuditHandler{service: service}
}

// GetAll GET /api/v1/audit
func (h *AuditHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	entries, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, entries)
}

// GetByEntity GET /api/v1/audit/{entityId}
func (h *AuditHandler) GetByEntity(w http.ResponseWriter, r *http.Request) {
	entityID := chi.URLParam(r, "entityId")

	entries, err := h.service.GetByEntityID(r.Context(), entityID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, entries)
}
