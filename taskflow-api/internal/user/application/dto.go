package application

import "time"

type CreateUserDTO struct {
}

type UserDTO struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
