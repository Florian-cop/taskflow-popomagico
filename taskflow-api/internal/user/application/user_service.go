package application

import (
	"context"
	shared "taskflow-api/internal/shared/application"
	userDomain "taskflow-api/internal/user/domain"
)

type UserService struct {
	repo     userDomain.UserRepository
	hasher   userDomain.PasswordHasher
	eventBus shared.EventBus
}

func NewUserService(repo userDomain.UserRepository, eventBus shared.EventBus) *UserService {
	return &UserService{repo: repo, eventBus: eventBus}
}

func (s *UserService) Register(ctx context.Context, dto UserDTO) (*UserDTO, error) {
	s.repo.FindByEmail(ctx, dto.Email)
	s.hasher.Hash(ctx, dto.)
}

func toDTO(t *userDomain.User) *UserDTO {
	return &UserDTO{
		ID:        t.ID,
		FirstName: t.FirstName,
		LastName:  t.LastName,
		Email:     t.Email,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
