package handlers

import (
	"net/http"
	"strconv"

	"taskflow-api/api/v1/dto"
	auditApp "taskflow-api/internal/audit/application"
)

type AuditHandler struct {
	service *auditApp.AuditService
}

func NewAuditHandler(s *auditApp.AuditService) *AuditHandler {
	return &AuditHandler{service: s}
}

// GET /api/v1/audit/logs?aggregateType=&aggregateId=&userId=&limit=
func (h *AuditHandler) Query(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))

	entries, err := h.service.Query(r.Context(), auditApp.QueryDTO{
		AggregateType: q.Get("aggregateType"),
		AggregateID:   q.Get("aggregateId"),
		UserID:        q.Get("userId"),
		Limit:         limit,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]dto.AuditLogResponse, len(entries))
	for i, e := range entries {
		resp[i] = dto.AuditLogResponse{
			ID: e.ID, UserID: e.UserID, EventName: e.EventName,
			AggregateType: e.AggregateType, AggregateID: e.AggregateID,
			Payload: e.Payload, OccurredAt: e.OccurredAt,
		}
	}
	writeJSON(w, http.StatusOK, resp)
}
