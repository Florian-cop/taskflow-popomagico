package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"taskflow-api/api/v1/handlers"
	projectApp "taskflow-api/internal/project/application"
	projectInfra "taskflow-api/internal/project/infrastructure"
	"taskflow-api/internal/shared/infrastructure/config"
	"taskflow-api/internal/shared/infrastructure/console"
	"taskflow-api/internal/shared/infrastructure/memory"
	"taskflow-api/internal/shared/infrastructure/persistence"
	taskApp "taskflow-api/internal/task/application"
	taskInfra "taskflow-api/internal/task/infrastructure"
)

func main() {
	// 1. Configuration
	cfg := config.Load()

	// 2. Base de données
	db := persistence.NewDatabase(cfg.DSN())
	persistence.Migrate(db,
		&projectInfra.ProjectModel{},
		&projectInfra.MemberModel{},
		&taskInfra.TaskModel{},
	)

	// 3. Event Bus + handlers
	eventBus := memory.NewInMemoryEventBus()
	eventBus.Subscribe("task.created", console.Handle)
	eventBus.Subscribe("task.moved", console.Handle)

	// 4. Repositories
	projectRepo := projectInfra.NewGormProjectRepository(db)
	taskRepo := taskInfra.NewGormTaskRepository(db)

	// 5. Services applicatifs
	projectService := projectApp.NewProjectService(projectRepo, eventBus)
	taskService := taskApp.NewTaskService(taskRepo, eventBus)

	// 6. Handlers HTTP
	projectHandler := handlers.NewProjectHandler(projectService)
	taskHandler := handlers.NewTaskHandler(taskService)

	// 7. Routeur
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "X-User-Id"},
	}))

	r.Route("/api/v1", func(r chi.Router) {
		// Projects
		r.Get("/projects", projectHandler.GetAllProjects)
		r.Post("/projects", projectHandler.CreateProject)
		r.Get("/projects/{id}", projectHandler.GetProject)
		r.Post("/projects/{id}/members", projectHandler.AddMember)

		// Tasks
		r.Get("/projects/{id}/tasks", taskHandler.GetTasksByProject)
		r.Post("/projects/{id}/tasks", taskHandler.CreateTask)
		r.Put("/tasks/{id}/move", taskHandler.MoveTask)
	})

	// 8. Démarrage
	log.Printf("serveur démarré sur :%s", cfg.APIPort)
	if err := http.ListenAndServe(":"+cfg.APIPort, r); err != nil {
		log.Fatalf("erreur serveur: %v", err)
	}
}
