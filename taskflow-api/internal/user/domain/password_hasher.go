package domain

import "context"

type PasswordHasher interface {
	Validate(ctx context.Context, hashedPassword string, plainPassword string) error
	Hash(ctx context.Context, password string) (string, error)
}
