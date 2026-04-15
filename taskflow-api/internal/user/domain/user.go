package domain

import "time"

type User struct {
	ID           string
	FirstName    string
	LastName     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(id, firstName, lastName, email, passwordHash string) *User {
	return &User{
		ID:           id,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
