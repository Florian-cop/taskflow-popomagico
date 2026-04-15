package application

import (
	"context"
	"errors"

	shared "taskflow-api/internal/shared/application"
	sharedDomain "taskflow-api/internal/shared/domain"
	userDomain "taskflow-api/internal/user/domain"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type UserService struct {
	repo     userDomain.UserRepository
	hasher   userDomain.PasswordHasher
	tokens   userDomain.TokenGenerator
	eventBus shared.EventBus
}

func NewUserService(
	repo userDomain.UserRepository,
	hasher userDomain.PasswordHasher,
	tokens userDomain.TokenGenerator,
	eventBus shared.EventBus,
) *UserService {
	return &UserService{repo: repo, hasher: hasher, tokens: tokens, eventBus: eventBus}
}

func (s *UserService) Register(ctx context.Context, dto RegisterUserDTO) (*AuthResultDTO, error) {
	if existing, err := s.repo.FindByEmail(ctx, dto.Email); err == nil && existing != nil {
		return nil, sharedDomain.ErrConflict
	} else if err != nil && !errors.Is(err, sharedDomain.ErrNotFound) {
		return nil, err
	}

	hashed, err := s.hasher.Hash(ctx, dto.Password)
	if err != nil {
		return nil, err
	}

	user := userDomain.NewUser(sharedDomain.NewID(), dto.FirstName, dto.LastName, dto.Email, hashed)
	if err := s.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	token, err := s.tokens.Generate(ctx, user)
	if err != nil {
		return nil, err
	}

	event := userDomain.NewUserCreatedEvent(user.ID, user.Email, user.FirstName, user.LastName)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		return nil, err
	}

	return &AuthResultDTO{Token: token, User: *toDTO(user)}, nil
}

func (s *UserService) Login(ctx context.Context, dto LoginDTO) (*AuthResultDTO, error) {
	user, err := s.repo.FindByEmail(ctx, dto.Email)
	if err != nil {
		if errors.Is(err, sharedDomain.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := s.hasher.Validate(ctx, user.PasswordHash, dto.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.tokens.Generate(ctx, user)
	if err != nil {
		return nil, err
	}

	return &AuthResultDTO{Token: token, User: *toDTO(user)}, nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (*UserDTO, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDTO(user), nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*UserDTO, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return toDTO(user), nil
}

func (s *UserService) SearchUsersByEmail(ctx context.Context, query string, limit int) ([]*UserDTO, error) {
	users, err := s.repo.SearchByEmail(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	out := make([]*UserDTO, len(users))
	for i, u := range users {
		out[i] = toDTO(u)
	}
	return out, nil
}

func toDTO(u *userDomain.User) *UserDTO {
	return &UserDTO{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
