package application

import (
	"context"
	"testing"

	shared "taskflow-api/internal/shared/application"
	"taskflow-api/internal/task/domain"
)

func newTestService() (*TaskService, *MockTaskRepository, *shared.MockEventBus) {
	repo := NewMockTaskRepository()
	bus := shared.NewMockEventBus()
	service := NewTaskService(repo, bus)
	return service, repo, bus
}

func TestCreateTask(t *testing.T) {
	service, repo, bus := newTestService()
	ctx := context.Background()

	task, err := service.CreateTask(ctx, CreateTaskDTO{
		Title:       "Ma tâche",
		Description: "Description",
		ProjectID:   "project-1",
	})

	if err != nil {
		t.Fatalf("erreur inattendue: %v", err)
	}
	if task.Title != "Ma tâche" {
		t.Errorf("titre attendu 'Ma tâche', reçu '%s'", task.Title)
	}
	if task.Status != "todo" {
		t.Errorf("statut attendu 'todo', reçu '%s'", task.Status)
	}
	if len(repo.Tasks) != 1 {
		t.Errorf("1 tâche attendue dans le repo, reçu %d", len(repo.Tasks))
	}
	if len(bus.Published) != 1 {
		t.Fatalf("1 event attendu, reçu %d", len(bus.Published))
	}
	if bus.Published[0].EventName() != "task.created" {
		t.Errorf("event attendu 'task.created', reçu '%s'", bus.Published[0].EventName())
	}
}

func TestMoveTask_ValidTransition(t *testing.T) {
	service, repo, bus := newTestService()
	ctx := context.Background()

	// Créer une tâche en statut Todo
	repo.Tasks["task-1"] = domain.NewTask("task-1", "Test", "Desc", "project-1")

	task, err := service.MoveTask(ctx, MoveTaskDTO{
		TaskID:    "task-1",
		NewStatus: "in_progress",
	})

	if err != nil {
		t.Fatalf("erreur inattendue: %v", err)
	}
	if task.Status != "in_progress" {
		t.Errorf("statut attendu 'in_progress', reçu '%s'", task.Status)
	}
	if len(bus.Published) != 1 {
		t.Fatalf("1 event attendu, reçu %d", len(bus.Published))
	}
	if bus.Published[0].EventName() != "task.moved" {
		t.Errorf("event attendu 'task.moved', reçu '%s'", bus.Published[0].EventName())
	}
}

func TestMoveTask_InvalidTransition(t *testing.T) {
	service, repo, _ := newTestService()
	ctx := context.Background()

	// Tâche en statut Todo → Done (interdit)
	repo.Tasks["task-1"] = domain.NewTask("task-1", "Test", "Desc", "project-1")

	_, err := service.MoveTask(ctx, MoveTaskDTO{
		TaskID:    "task-1",
		NewStatus: "done",
	})

	if err == nil {
		t.Fatal("une erreur était attendue pour une transition invalide")
	}
}

func TestMoveTask_NotFound(t *testing.T) {
	service, _, _ := newTestService()
	ctx := context.Background()

	_, err := service.MoveTask(ctx, MoveTaskDTO{
		TaskID:    "inexistant",
		NewStatus: "in_progress",
	})

	if err == nil {
		t.Fatal("une erreur était attendue pour une tâche inexistante")
	}
}

func TestGetTasksByProject(t *testing.T) {
	service, repo, _ := newTestService()
	ctx := context.Background()

	repo.Tasks["t1"] = domain.NewTask("t1", "Tâche 1", "", "project-1")
	repo.Tasks["t2"] = domain.NewTask("t2", "Tâche 2", "", "project-1")
	repo.Tasks["t3"] = domain.NewTask("t3", "Tâche 3", "", "project-2")

	tasks, err := service.GetTasksByProject(ctx, "project-1")
	if err != nil {
		t.Fatalf("erreur inattendue: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("2 tâches attendues, reçu %d", len(tasks))
	}
}
