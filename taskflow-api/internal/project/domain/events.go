package domain

import (
	"time"

	shared "github.com/Floxtouille/taskflow-popomagico/taskflow-api/internal/shared/domain"
)

// --- project.created ---

type ProjectCreatedEvent struct {
	id         string
	projectID  string
	name       string
	ownerID    string
	occurredAt time.Time
}

func NewProjectCreatedEvent(projectID, name, ownerID string) ProjectCreatedEvent {
	return ProjectCreatedEvent{
		id:         shared.NewID(),
		projectID:  projectID,
		name:       name,
		ownerID:    ownerID,
		occurredAt: time.Now(),
	}
}

func (e ProjectCreatedEvent) EventID() string      { return e.id }
func (e ProjectCreatedEvent) EventName() string     { return "project.created" }
func (e ProjectCreatedEvent) AggregateID() string   { return e.projectID }
func (e ProjectCreatedEvent) OccurredAt() time.Time { return e.occurredAt }
func (e ProjectCreatedEvent) OwnerID() string       { return e.ownerID }

// --- member.added ---

type MemberAddedEvent struct {
	id         string
	projectID  string
	userID     string
	occurredAt time.Time
}

func NewMemberAddedEvent(projectID, userID string) MemberAddedEvent {
	return MemberAddedEvent{
		id:         shared.NewID(),
		projectID:  projectID,
		userID:     userID,
		occurredAt: time.Now(),
	}
}

func (e MemberAddedEvent) EventID() string      { return e.id }
func (e MemberAddedEvent) EventName() string     { return "member.added" }
func (e MemberAddedEvent) AggregateID() string   { return e.projectID }
func (e MemberAddedEvent) OccurredAt() time.Time { return e.occurredAt }
func (e MemberAddedEvent) UserID() string        { return e.userID }
