package application

import (
	"context"

	projectDomain "taskflow-api/internal/project/domain"
	shared "taskflow-api/internal/shared/application"
	"taskflow-api/internal/shared/domain"
)

type ProjectService struct {
	repo     projectDomain.ProjectRepository
	eventBus shared.EventBus
}

func NewProjectService(repo projectDomain.ProjectRepository, eventBus shared.EventBus) *ProjectService {
	return &ProjectService{repo: repo, eventBus: eventBus}
}

func (s *ProjectService) CreateProject(ctx context.Context, dto CreateProjectDTO) (*ProjectDTO, error) {
	project := projectDomain.NewProject(domain.NewID(), dto.Name, dto.Description, dto.OwnerID)

	if err := s.repo.Save(ctx, project); err != nil {
		return nil, err
	}

	event := projectDomain.NewProjectCreatedEvent(project.ID, project.Name, dto.OwnerID)
	if err := s.eventBus.Publish(ctx, event); err != nil {
		return nil, err
	}

	return toDTO(project), nil
}

func (s *ProjectService) AddMember(ctx context.Context, dto AddMemberDTO) (*ProjectDTO, error) {
	project, err := s.repo.FindByID(ctx, dto.ProjectID)
	if err != nil {
		return nil, err
	}

	event, err := project.AddMember(dto.UserID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, project); err != nil {
		return nil, err
	}

	if err := s.eventBus.Publish(ctx, event); err != nil {
		return nil, err
	}

	return toDTO(project), nil
}

func (s *ProjectService) GetProject(ctx context.Context, id string) (*ProjectDTO, error) {
	project, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDTO(project), nil
}

func (s *ProjectService) GetAllProjects(ctx context.Context) ([]*ProjectDTO, error) {
	projects, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]*ProjectDTO, len(projects))
	for i, p := range projects {
		dtos[i] = toDTO(p)
	}
	return dtos, nil
}

func toDTO(p *projectDomain.Project) *ProjectDTO {
	members := make([]MemberDTO, len(p.Members))
	for i, m := range p.Members {
		members[i] = MemberDTO{
			UserID:   m.UserID,
			Role:     string(m.Role),
			JoinedAt: m.JoinedAt,
		}
	}
	return &ProjectDTO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Members:     members,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
