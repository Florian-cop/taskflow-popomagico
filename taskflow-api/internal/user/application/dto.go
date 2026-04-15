package application

import "time"

type RegisterUserDTO struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

type LoginDTO struct {
	Email    string
	Password string
}

type UserDTO struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AuthResultDTO struct {
	Token string
	User  UserDTO
}
