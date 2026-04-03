package infrastructure

import (
	"encoding/json"

	auditDomain "taskflow-api/internal/audit/domain"
)

func toModel(entry *auditDomain.AuditEntry) *AuditEntryModel {
	detailsJSON, _ := json.Marshal(entry.Details)
	return &AuditEntryModel{
		ID:         entry.ID,
		Actor:      entry.Actor,
		Action:     entry.Action,
		EntityType: entry.EntityType,
		EntityID:   entry.EntityID,
		Details:    string(detailsJSON),
		OccurredAt: entry.OccurredAt,
	}
}

func toDomain(model *AuditEntryModel) *auditDomain.AuditEntry {
	details := make(map[string]string)
	_ = json.Unmarshal([]byte(model.Details), &details)
	return &auditDomain.AuditEntry{
		ID:         model.ID,
		Actor:      model.Actor,
		Action:     model.Action,
		EntityType: model.EntityType,
		EntityID:   model.EntityID,
		Details:    details,
		OccurredAt: model.OccurredAt,
	}
}
