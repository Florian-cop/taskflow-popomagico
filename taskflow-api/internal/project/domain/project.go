package domain

import (
	"errors"
	"time"
)

var ErrMemberAlreadyExists = errors.New("member already exists in project")

type Project struct {
	ID          string
	Name        string
	Description string
	Members     []Member
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewProject crée un projet et ajoute le créateur comme owner.
func NewProject(id, name, description, ownerID string) *Project {
	now := time.Now()
	return &Project{
		ID:          id,
		Name:        name,
		Description: description,
		Members:     []Member{NewMember(ownerID, RoleOwner)},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// AddMember ajoute un membre au projet s'il n'en fait pas déjà partie.
func (p *Project) AddMember(userID string) (MemberAddedEvent, error) {
	for _, m := range p.Members {
		if m.UserID == userID {
			return MemberAddedEvent{}, ErrMemberAlreadyExists
		}
	}
	p.Members = append(p.Members, NewMember(userID, RoleMember))
	p.UpdatedAt = time.Now()
	return NewMemberAddedEvent(p.ID, userID), nil
}

// HasMember vérifie si un utilisateur est membre du projet.
func (p *Project) HasMember(userID string) bool {
	for _, m := range p.Members {
		if m.UserID == userID {
			return true
		}
	}
	return false
}
