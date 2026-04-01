package domain

import (
	"time"

	shared "github.com/Floxtouille/taskflow-popomagico/taskflow-api/internal/shared/domain"
)

// --- task.created ---

type TaskCreatedEvent struct {
	id         string
	taskID     string
	projectID  string
	title      string
	occurredAt time.Time
}

func NewTaskCreatedEvent(taskID, projectID, title string) TaskCreatedEvent {
	return TaskCreatedEvent{
		id:         shared.NewID(),
		taskID:     taskID,
		projectID:  projectID,
		title:      title,
		occurredAt: time.Now(),
	}
}

func (e TaskCreatedEvent) EventID() string      { return e.id }
func (e TaskCreatedEvent) EventName() string     { return "task.created" }
func (e TaskCreatedEvent) AggregateID() string   { return e.taskID }
func (e TaskCreatedEvent) OccurredAt() time.Time { return e.occurredAt }
func (e TaskCreatedEvent) ProjectID() string     { return e.projectID }
func (e TaskCreatedEvent) Title() string         { return e.title }

// --- task.moved ---

type TaskMovedEvent struct {
	id         string
	taskID     string
	projectID  string
	fromStatus TaskStatus
	toStatus   TaskStatus
	occurredAt time.Time
}

func NewTaskMovedEvent(taskID, projectID string, from, to TaskStatus) TaskMovedEvent {
	return TaskMovedEvent{
		id:         shared.NewID(),
		taskID:     taskID,
		projectID:  projectID,
		fromStatus: from,
		toStatus:   to,
		occurredAt: time.Now(),
	}
}

func (e TaskMovedEvent) EventID() string        { return e.id }
func (e TaskMovedEvent) EventName() string      { return "task.moved" }
func (e TaskMovedEvent) AggregateID() string    { return e.taskID }
func (e TaskMovedEvent) OccurredAt() time.Time  { return e.occurredAt }
func (e TaskMovedEvent) FromStatus() TaskStatus { return e.fromStatus }
func (e TaskMovedEvent) ToStatus() TaskStatus   { return e.toStatus }

// --- task.assigned ---

type TaskAssignedEvent struct {
	id         string
	taskID     string
	projectID  string
	assigneeID string
	occurredAt time.Time
}

func NewTaskAssignedEvent(taskID, projectID, assigneeID string) TaskAssignedEvent {
	return TaskAssignedEvent{
		id:         shared.NewID(),
		taskID:     taskID,
		projectID:  projectID,
		assigneeID: assigneeID,
		occurredAt: time.Now(),
	}
}

func (e TaskAssignedEvent) EventID() string      { return e.id }
func (e TaskAssignedEvent) EventName() string     { return "task.assigned" }
func (e TaskAssignedEvent) AggregateID() string   { return e.taskID }
func (e TaskAssignedEvent) OccurredAt() time.Time { return e.occurredAt }
func (e TaskAssignedEvent) AssigneeID() string    { return e.assigneeID }
