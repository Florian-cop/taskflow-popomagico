package infrastructure

import (
	"context"

	projectApp "taskflow-api/internal/project/application"
)

// ProjectMemberFinder adapte le bounded context project à l'interface MemberFinder
// du bounded context notification. Ce coupling cross-context est isolé à la couche
// infra et ne touche jamais les domaines.
type ProjectMemberFinder struct {
	svc *projectApp.ProjectService
}

func NewProjectMemberFinder(s *projectApp.ProjectService) *ProjectMemberFinder {
	return &ProjectMemberFinder{svc: s}
}

func (f *ProjectMemberFinder) FindMembers(ctx context.Context, projectID string) ([]string, error) {
	p, err := f.svc.GetProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(p.Members))
	for i, m := range p.Members {
		ids[i] = m.UserID
	}
	return ids, nil
}
