package domain

import "time"

type Task struct {
	ID          string
	Title       string
	Description string
	Status      TaskStatus
	AssigneeID  *string
	ProjectID   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewTask crée une nouvelle tâche avec le statut Todo par défaut.
func NewTask(id, title, description, projectID string) *Task {
	now := time.Now()
	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      StatusTodo,
		ProjectID:   projectID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// MoveTo déplace la tâche vers un nouveau statut si la transition est valide.
func (t *Task) MoveTo(newStatus TaskStatus) (TaskMovedEvent, error) {
	if !t.Status.CanTransitionTo(newStatus) {
		return TaskMovedEvent{}, ErrInvalidTransition
	}
	from := t.Status
	t.Status = newStatus
	t.UpdatedAt = time.Now()
	return NewTaskMovedEvent(t.ID, t.ProjectID, from, newStatus), nil
}

// AssignTo assigne la tâche à un membre.
func (t *Task) AssignTo(memberID string) TaskAssignedEvent {
	t.AssigneeID = &memberID
	t.UpdatedAt = time.Now()
	return NewTaskAssignedEvent(t.ID, t.ProjectID, memberID)
}
