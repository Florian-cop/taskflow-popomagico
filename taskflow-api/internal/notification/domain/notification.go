package domain

import "time"

type Notification struct {
	ID        string
	UserID    string
	Type      string
	Title     string
	Body      string
	ReadAt    *time.Time
	CreatedAt time.Time
}

func NewNotification(id, userID, notifType, title, body string) *Notification {
	return &Notification{
		ID:        id,
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Body:      body,
		CreatedAt: time.Now(),
	}
}

func (n *Notification) MarkAsRead() {
	if n.ReadAt != nil {
		return
	}
	t := time.Now()
	n.ReadAt = &t
}
