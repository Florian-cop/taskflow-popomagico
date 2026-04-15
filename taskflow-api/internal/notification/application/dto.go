package application

import "time"

type NotificationDTO struct {
	ID        string
	UserID    string
	Type      string
	Title     string
	Body      string
	ReadAt    *time.Time
	CreatedAt time.Time
}

type PreferencesDTO struct {
	UserID  string
	Enabled map[string]bool
}

type UpdatePreferencesDTO struct {
	UserID  string
	Enabled map[string]bool
}
