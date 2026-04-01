package infrastructure

import "time"

type ProjectModel struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (ProjectModel) TableName() string { return "projects" }

type MemberModel struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	ProjectID string `gorm:"index"`
	UserID    string
	Role      string
	JoinedAt  time.Time
}

func (MemberModel) TableName() string { return "members" }
