package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	projectApp "taskflow-api/internal/project/application"
	projectInfra "taskflow-api/internal/project/infrastructure"
	"taskflow-api/internal/shared/infrastructure/config"
	"taskflow-api/internal/shared/infrastructure/console"
	"taskflow-api/internal/shared/infrastructure/memory"
	"taskflow-api/internal/shared/infrastructure/persistence"
	taskApp "taskflow-api/internal/task/application"
	taskInfra "taskflow-api/internal/task/infrastructure"
)

// CLI d'administration TaskFlow.
// Réutilise EXACTEMENT les mêmes services applicatifs que l'API REST :
// aucune duplication de logique métier. C'est un adaptateur entrant alternatif
// (au même titre que l'API HTTP), câblé sur les mêmes ports.
//
// Usage :
//   cli project create --name "Mon projet" --owner "user-id"
//   cli task create --project <id> --title "Titre"
//   cli generate-demo [--owner <id>]
func main() {
	if len(os.Args) < 2 {
		printUsageAndExit()
	}

	ctx := context.Background()
	cfg := config.Load()

	db := persistence.NewDatabase(cfg.DSN())
	persistence.Migrate(db,
		&projectInfra.ProjectModel{},
		&projectInfra.MemberModel{},
		&taskInfra.TaskModel{},
	)

	bus := memory.NewInMemoryEventBus()
	for _, name := range []string{
		"project.created", "member.added",
		"task.created", "task.moved", "task.assigned",
	} {
		bus.Subscribe(name, console.Handle)
	}

	projectSvc := projectApp.NewProjectService(projectInfra.NewGormProjectRepository(db), bus)
	taskSvc := taskApp.NewTaskService(taskInfra.NewGormTaskRepository(db), bus)

	switch os.Args[1] {
	case "project":
		runProject(ctx, projectSvc)
	case "task":
		runTask(ctx, taskSvc)
	case "generate-demo":
		runGenerateDemo(ctx, projectSvc, taskSvc)
	default:
		printUsageAndExit()
	}
}

func runProject(ctx context.Context, svc *projectApp.ProjectService) {
	if len(os.Args) < 3 || os.Args[2] != "create" {
		fmt.Fprintln(os.Stderr, "usage: cli project create --name <name> [--description <d>] --owner <id>")
		os.Exit(1)
	}
	fs := flag.NewFlagSet("project create", flag.ExitOnError)
	name := fs.String("name", "", "project name")
	desc := fs.String("description", "", "project description")
	owner := fs.String("owner", "", "owner user id")
	_ = fs.Parse(os.Args[3:])

	if *name == "" || *owner == "" {
		fmt.Fprintln(os.Stderr, "--name et --owner sont requis")
		os.Exit(1)
	}

	p, err := svc.CreateProject(ctx, projectApp.CreateProjectDTO{
		Name: *name, Description: *desc, OwnerID: *owner,
	})
	if err != nil {
		log.Fatalf("erreur création projet: %v", err)
	}
	fmt.Printf("projet créé: id=%s name=%q\n", p.ID, p.Name)
}

func runTask(ctx context.Context, svc *taskApp.TaskService) {
	if len(os.Args) < 3 || os.Args[2] != "create" {
		fmt.Fprintln(os.Stderr, "usage: cli task create --project <id> --title <title> [--description <d>]")
		os.Exit(1)
	}
	fs := flag.NewFlagSet("task create", flag.ExitOnError)
	project := fs.String("project", "", "project id")
	title := fs.String("title", "", "task title")
	desc := fs.String("description", "", "task description")
	_ = fs.Parse(os.Args[3:])

	if *project == "" || *title == "" {
		fmt.Fprintln(os.Stderr, "--project et --title sont requis")
		os.Exit(1)
	}

	t, err := svc.CreateTask(ctx, taskApp.CreateTaskDTO{
		Title: *title, Description: *desc, ProjectID: *project,
	})
	if err != nil {
		log.Fatalf("erreur création tâche: %v", err)
	}
	fmt.Printf("tâche créée: id=%s title=%q status=%s\n", t.ID, t.Title, t.Status)
}

func runGenerateDemo(ctx context.Context, projectSvc *projectApp.ProjectService, taskSvc *taskApp.TaskService) {
	fs := flag.NewFlagSet("generate-demo", flag.ExitOnError)
	owner := fs.String("owner", "demo-user", "owner id for demo project")
	_ = fs.Parse(os.Args[2:])

	project, err := projectSvc.CreateProject(ctx, projectApp.CreateProjectDTO{
		Name: "Projet de démonstration", Description: "Jeu de données de démo", OwnerID: *owner,
	})
	if err != nil {
		log.Fatalf("erreur création projet demo: %v", err)
	}
	fmt.Printf("projet demo créé: %s\n", project.ID)

	targets := []struct {
		title  string
		status string
	}{
		{"Analyser le besoin client", "todo"},
		{"Rédiger les user stories", "todo"},
		{"Prototyper la maquette", "todo"},
		{"Implémenter l'authentification", "in_progress"},
		{"Mettre en place le Kanban temps réel", "in_progress"},
		{"Configurer la CI", "in_progress"},
		{"Écrire les tests unitaires initiaux", "done"},
		{"Docker Compose de base", "done"},
		{"ADR architecture hexagonale", "done"},
	}

	for _, t := range targets {
		created, err := taskSvc.CreateTask(ctx, taskApp.CreateTaskDTO{
			Title: t.title, ProjectID: project.ID,
		})
		if err != nil {
			log.Fatalf("erreur création tâche %q: %v", t.title, err)
		}

		// Progression via les transitions autorisées : todo → in_progress → done
		switch t.status {
		case "in_progress":
			if _, err := taskSvc.MoveTask(ctx, taskApp.MoveTaskDTO{
				TaskID: created.ID, NewStatus: "in_progress",
			}); err != nil {
				log.Fatalf("erreur move→in_progress %s: %v", created.ID, err)
			}
		case "done":
			if _, err := taskSvc.MoveTask(ctx, taskApp.MoveTaskDTO{
				TaskID: created.ID, NewStatus: "in_progress",
			}); err != nil {
				log.Fatalf("erreur move→in_progress %s: %v", created.ID, err)
			}
			if _, err := taskSvc.MoveTask(ctx, taskApp.MoveTaskDTO{
				TaskID: created.ID, NewStatus: "done",
			}); err != nil {
				log.Fatalf("erreur move→done %s: %v", created.ID, err)
			}
		}
	}

	fmt.Printf("démo générée: projet=%s, 9 tâches (3 todo / 3 in_progress / 3 done)\n", project.ID)
}

func printUsageAndExit() {
	fmt.Fprintln(os.Stderr, `TaskFlow CLI — administration
Usage:
  cli project create --name <name> [--description <d>] --owner <user-id>
  cli task create --project <id> --title <title> [--description <d>]
  cli generate-demo [--owner <user-id>]`)
	os.Exit(1)
}
