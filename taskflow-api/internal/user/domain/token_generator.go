package domain

import "context"

type TokenGenerator interface {
	Validate(ctx context.Context, token string) (*User, error)
	Generate(ctx context.Context, user *User) (string, error)
}
