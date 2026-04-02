package main

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	projectInfra "taskflow-api/internal/project/infrastructure"
	"taskflow-api/internal/shared/infrastructure/config"
	"taskflow-api/internal/shared/infrastructure/console"
	"taskflow-api/internal/shared/infrastructure/memory"
	"taskflow-api/internal/shared/infrastructure/persistence"
	taskApp "taskflow-api/internal/task/application"
	taskInfra "taskflow-api/internal/task/infrastructure"
)

var createTaskCmd = &cobra.Command{
	Use:   "create-task",
	Short: "Create a new task",
	Long:  "Create a new task in an existing project",
	Run: func(cmd *cobra.Command, args []string) {
		title, _ := cmd.Flags().GetString("title")
		projectID, _ := cmd.Flags().GetString("project-id")

		cfg := config.Load()
		db := persistence.NewDatabase(cfg.DSN())
		persistence.Migrate(db,
			&projectInfra.ProjectModel{},
			&projectInfra.MemberModel{},
			&taskInfra.TaskModel{},
		)

		eventBus := memory.NewInMemoryEventBus()
		eventBus.Subscribe("task.created", console.Handle)

		taskRepo := taskInfra.NewGormTaskRepository(db)
		taskService := taskApp.NewTaskService(taskRepo, eventBus)

		dto := taskApp.CreateTaskDTO{
			Title:     title,
			ProjectID: projectID,
		}

		task, err := taskService.CreateTask(context.Background(), dto)
		if err != nil {
			log.Fatalf("erreur lors de la création de la tâche: %v", err)
		}

		fmt.Printf("Tâche créée avec succès !\n  ID:      %s\n  Titre:   %s\n  Projet:  %s\n  Statut:  %s\n", task.ID, task.Title, task.ProjectID, task.Status)
	},
}

func init() {
	createTaskCmd.Flags().String("title", "", "Titre de la tâche")
	createTaskCmd.Flags().String("project-id", "", "ID du projet")
	createTaskCmd.MarkFlagRequired("title")
	createTaskCmd.MarkFlagRequired("project-id")
	rootCmd.AddCommand(createTaskCmd)
}
