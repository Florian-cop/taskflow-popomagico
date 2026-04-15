package infrastructure

import "time"

type AuditLogModel struct {
	ID            string `gorm:"primaryKey"`
	UserID        string `gorm:"index"`
	EventName     string `gorm:"index"`
	AggregateType string `gorm:"index"`
	AggregateID   string `gorm:"index"`
	Payload       string
	OccurredAt    time.Time
}

func (AuditLogModel) TableName() string { return "audit_logs" }
