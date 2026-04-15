package dto

import "time"

type NotificationResponse struct {
	ID        string     `json:"id"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	ReadAt    *time.Time `json:"readAt"`
	CreatedAt time.Time  `json:"createdAt"`
}

type PreferencesResponse struct {
	Enabled map[string]bool `json:"enabled"`
}

type UpdatePreferencesRequest struct {
	Enabled map[string]bool `json:"enabled"`
}
