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
	taskInfra "taskflow-api/internal/task/infrastructure"
)

var createProjectCmd = &cobra.Command{
	Use:   "create-project",
	Short: "Create a new project",
	Long:  "Create a new project in the TaskFlow API",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")

		// Bootstrap (même wiring que l'API)
		cfg := config.Load()
		db := persistence.NewDatabase(cfg.DSN())
		persistence.Migrate(db,
			&projectInfra.ProjectModel{},
			&projectInfra.MemberModel{},
			&taskInfra.TaskModel{},
		)

		eventBus := memory.NewInMemoryEventBus()
		eventBus.Subscribe("project.created", console.Handle)

		projectRepo := projectInfra.NewGormProjectRepository(db)
		projectService := projectApp.NewProjectService(projectRepo, eventBus)

		// Appel du service applicatif
		dto := projectApp.CreateProjectDTO{
			Name:    name,
			OwnerID: "cli-user",
		}

		project, err := projectService.CreateProject(context.Background(), dto)
		if err != nil {
			log.Fatalf("erreur lors de la création du projet: %v", err)
		}

		fmt.Printf("Projet créé avec succès !\n  ID:   %s\n  Nom:  %s\n", project.ID, project.Name)
	},
}

func init() {
	createProjectCmd.Flags().String("name", "", "Nom du projet")
	createProjectCmd.MarkFlagRequired("name")
	rootCmd.AddCommand(createProjectCmd)
}
