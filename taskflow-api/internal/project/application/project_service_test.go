package application

import (
	"context"
	"testing"

	shared "taskflow-api/internal/shared/application"
)

func newTestService() (*ProjectService, *MockProjectRepository, *shared.MockEventBus) {
	repo := NewMockProjectRepository()
	bus := shared.NewMockEventBus()
	service := NewProjectService(repo, bus)
	return service, repo, bus
}

func TestCreateProject(t *testing.T) {
	service, repo, bus := newTestService()
	ctx := context.Background()

	project, err := service.CreateProject(ctx, CreateProjectDTO{
		Name:        "Mon projet",
		Description: "Description",
		OwnerID:     "user-1",
	})

	if err != nil {
		t.Fatalf("erreur inattendue: %v", err)
	}
	if project.Name != "Mon projet" {
		t.Errorf("nom attendu 'Mon projet', reçu '%s'", project.Name)
	}
	if len(project.Members) != 1 {
		t.Fatalf("1 membre attendu, reçu %d", len(project.Members))
	}
	if project.Members[0].Role != "owner" {
		t.Errorf("rôle attendu 'owner', reçu '%s'", project.Members[0].Role)
	}
	if len(repo.Projects) != 1 {
		t.Errorf("1 projet attendu dans le repo, reçu %d", len(repo.Projects))
	}
	if len(bus.Published) != 1 {
		t.Fatalf("1 event attendu, reçu %d", len(bus.Published))
	}
	if bus.Published[0].EventName() != "project.created" {
		t.Errorf("event attendu 'project.created', reçu '%s'", bus.Published[0].EventName())
	}
}

func TestAddMember(t *testing.T) {
	service, _, bus := newTestService()
	ctx := context.Background()

	// Créer un projet d'abord
	project, _ := service.CreateProject(ctx, CreateProjectDTO{
		Name:    "Projet",
		OwnerID: "user-1",
	})

	// Ajouter un membre
	updated, err := service.AddMember(ctx, AddMemberDTO{
		ProjectID: project.ID,
		UserID:    "user-2",
	})

	if err != nil {
		t.Fatalf("erreur inattendue: %v", err)
	}
	if len(updated.Members) != 2 {
		t.Errorf("2 membres attendus, reçu %d", len(updated.Members))
	}
	// project.created + member.added
	if len(bus.Published) != 2 {
		t.Fatalf("2 events attendus, reçu %d", len(bus.Published))
	}
	if bus.Published[1].EventName() != "member.added" {
		t.Errorf("event attendu 'member.added', reçu '%s'", bus.Published[1].EventName())
	}
}

func TestAddMember_AlreadyExists(t *testing.T) {
	service, _, _ := newTestService()
	ctx := context.Background()

	project, _ := service.CreateProject(ctx, CreateProjectDTO{
		Name:    "Projet",
		OwnerID: "user-1",
	})

	// Ajouter user-1 qui est déjà owner
	_, err := service.AddMember(ctx, AddMemberDTO{
		ProjectID: project.ID,
		UserID:    "user-1",
	})

	if err == nil {
		t.Fatal("une erreur était attendue pour un membre déjà existant")
	}
}

func TestGetProject_NotFound(t *testing.T) {
	service, _, _ := newTestService()
	ctx := context.Background()

	_, err := service.GetProject(ctx, "inexistant")

	if err == nil {
		t.Fatal("une erreur était attendue pour un projet inexistant")
	}
}

func TestGetAllProjects(t *testing.T) {
	service, _, _ := newTestService()
	ctx := context.Background()

	service.CreateProject(ctx, CreateProjectDTO{Name: "Projet 1", OwnerID: "u1"})
	service.CreateProject(ctx, CreateProjectDTO{Name: "Projet 2", OwnerID: "u2"})

	projects, err := service.GetAllProjects(ctx)
	if err != nil {
		t.Fatalf("erreur inattendue: %v", err)
	}
	if len(projects) != 2 {
		t.Errorf("2 projets attendus, reçu %d", len(projects))
	}
}
