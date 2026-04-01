package infrastructure

import (
	"taskflow-api/internal/project/domain"
)

func toModel(p *domain.Project) (*ProjectModel, []MemberModel) {
	model := &ProjectModel{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}

	members := make([]MemberModel, len(p.Members))
	for i, m := range p.Members {
		members[i] = MemberModel{
			ProjectID: p.ID,
			UserID:    m.UserID,
			Role:      string(m.Role),
			JoinedAt:  m.JoinedAt,
		}
	}

	return model, members
}

func toDomain(model *ProjectModel, memberModels []MemberModel) *domain.Project {
	members := make([]domain.Member, len(memberModels))
	for i, m := range memberModels {
		members[i] = domain.Member{
			UserID:   m.UserID,
			Role:     domain.Role(m.Role),
			JoinedAt: m.JoinedAt,
		}
	}

	return &domain.Project{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		Members:     members,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}
