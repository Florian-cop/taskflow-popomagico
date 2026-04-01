package domain

import "time"

type Member struct {
	UserID   string
	Role     Role
	JoinedAt time.Time
}

func NewMember(userID string, role Role) Member {
	return Member{
		UserID:   userID,
		Role:     role,
		JoinedAt: time.Now(),
	}
}
