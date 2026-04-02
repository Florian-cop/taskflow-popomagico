package main

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	projectApp "taskflow-api/internal/project/application"
	projectInfra "taskflow-api/internal/project/infrastructure"
	"taskflow-api/internal/shared/infrastructure/config"
	"taskflow-api/internal/shared/infrastructure/console"
	"taskflow-api/internal/shared/infrastructure/memory"
	"taskflow-api/internal/shared/infrastructure/persistence"
	taskApp "taskflow-api/internal/task/application"
	taskInfra "taskflow-api/internal/task/infrastructure"
)

var generateDemoCmd = &cobra.Command{
	Use:   "generate-demo",
	Short: "Generate a demo dataset",
	Long:  "Create a demo project with tasks distributed across TODO, IN_PROGRESS, and DONE",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Bootstrap
		cfg := config.Load()
		db := persistence.NewDatabase(cfg.DSN())
		persistence.Migrate(db,
			&projectInfra.ProjectModel{},
			&projectInfra.MemberModel{},
			&taskInfra.TaskModel{},
		)

		eventBus := memory.NewInMemoryEventBus()
		eventBus.Subscribe("project.created", console.Handle)
		eventBus.Subscribe("task.created", console.Handle)
		eventBus.Subscribe("task.moved", console.Handle)

		projectRepo := projectInfra.NewGormProjectRepository(db)
		taskRepo := taskInfra.NewGormTaskRepository(db)
		projectService := projectApp.NewProjectService(projectRepo, eventBus)
		taskService := taskApp.NewTaskService(taskRepo, eventBus)

		// Créer le projet démo
		project, err := projectService.CreateProject(ctx, projectApp.CreateProjectDTO{
			Name:    "Projet Démo",
			OwnerID: "cli-user",
		})
		if err != nil {
			log.Fatalf("erreur création projet: %v", err)
		}
		fmt.Printf("Projet créé: %s (%s)\n", project.Name, project.ID)

		// Tâches à créer avec leur statut cible
		tasks := []struct {
			title  string
			status string // "" = reste en todo, "in_progress", "done"
		}{
			{"Rédiger le cahier des charges", "done"},
			{"Configurer le CI/CD", "done"},
			{"Implémenter l'authentification", "in_progress"},
			{"Créer les endpoints REST", "in_progress"},
			{"Écrire les tests unitaires", ""},
			{"Déployer en production", ""},
		}

		for _, t := range tasks {
			task, err := taskService.CreateTask(ctx, taskApp.CreateTaskDTO{
				Title:     t.title,
				ProjectID: project.ID,
			})
			if err != nil {
				log.Fatalf("erreur création tâche %q: %v", t.title, err)
			}

			// Déplacer si nécessaire (todo -> in_progress -> done)
			if t.status == "in_progress" || t.status == "done" {
				task, err = taskService.MoveTask(ctx, taskApp.MoveTaskDTO{
					TaskID:    task.ID,
					NewStatus: "in_progress",
				})
				if err != nil {
					log.Fatalf("erreur déplacement tâche %q vers in_progress: %v", t.title, err)
				}
			}
			if t.status == "done" {
				_, err = taskService.MoveTask(ctx, taskApp.MoveTaskDTO{
					TaskID:    task.ID,
					NewStatus: "done",
				})
				if err != nil {
					log.Fatalf("erreur déplacement tâche %q vers done: %v", t.title, err)
				}
			}

			status := task.Status
			if t.status != "" {
				status = t.status
			}
			fmt.Printf("  [%-12s] %s\n", status, t.title)
		}

		fmt.Println("\nDataset démo généré avec succès !")
	},
}

func init() {
	rootCmd.AddCommand(generateDemoCmd)
}
