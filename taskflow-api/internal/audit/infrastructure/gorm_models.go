package infrastructure

import "time"

type AuditEntryModel struct {
	ID         string    `gorm:"primaryKey"`
	Actor      string    `gorm:"index"`
	Action     string
	EntityType string    `gorm:"index"`
	EntityID   string    `gorm:"index"`
	Details    string    // JSON serialise
	OccurredAt time.Time `gorm:"index"`
}

func (AuditEntryModel) TableName() string {
	return "audit_log"
}
