package infrastructure

import "time"

type NotificationModel struct {
	ID        string `gorm:"primaryKey"`
	UserID    string `gorm:"index"`
	Type      string
	Title     string
	Body      string
	ReadAt    *time.Time
	CreatedAt time.Time
}

func (NotificationModel) TableName() string { return "notifications" }

type PreferencesModel struct {
	UserID  string `gorm:"primaryKey"`
	Enabled string // JSON serialized map[string]bool
}

func (PreferencesModel) TableName() string { return "notification_preferences" }
