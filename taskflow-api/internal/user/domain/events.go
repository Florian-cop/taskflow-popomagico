package domain

import (
	shared "taskflow-api/internal/shared/domain"
	"time"
)

type UserCreatedEvent struct {
	id         string
	userID     string
	email      string
	firstName  string
	lastName   string
	occurredAt time.Time
}

func NewUserCreatedEvent(userID, email, firstName, lastName string) UserCreatedEvent {
	return UserCreatedEvent{
		id:         shared.NewID(),
		userID:     userID,
		email:      email,
		firstName:  firstName,
		lastName:   lastName,
		occurredAt: time.Now(),
	}
}

func (e UserCreatedEvent) EventID() string       { return e.id }
func (e UserCreatedEvent) EventName() string     { return "user.created" }
func (e UserCreatedEvent) AggregateID() string   { return e.userID }
func (e UserCreatedEvent) OccurredAt() time.Time { return e.occurredAt }
func (e UserCreatedEvent) Email() string         { return e.email }
func (e UserCreatedEvent) FirstName() string     { return e.firstName }
func (e UserCreatedEvent) LastName() string      { return e.lastName }
