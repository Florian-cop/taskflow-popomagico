package domain

import "context"

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	SearchByEmail(ctx context.Context, query string, limit int) ([]*User, error)
	Update(ctx context.Context, user *User) error
	Save(ctx context.Context, user *User) error
}
