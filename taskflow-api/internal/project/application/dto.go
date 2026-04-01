package application

import "time"

type CreateProjectDTO struct {
	Name        string
	Description string
	OwnerID     string
}

type AddMemberDTO struct {
	ProjectID string
	UserID    string
}

type MemberDTO struct {
	UserID   string
	Role     string
	JoinedAt time.Time
}

type ProjectDTO struct {
	ID          string
	Name        string
	Description string
	Members     []MemberDTO
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
