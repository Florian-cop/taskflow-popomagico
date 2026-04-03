package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"taskflow-api/api/v1/handlers"
	customMiddleware "taskflow-api/api/v1/middleware"
	auditApp "taskflow-api/internal/audit/application"
	auditInfra "taskflow-api/internal/audit/infrastructure"
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
		&auditInfra.AuditEntryModel{},
	)

	// 3. Event Bus + handlers
	eventBus := memory.NewInMemoryEventBus()
	eventBus.Subscribe("project.created", console.Handle)
	eventBus.Subscribe("member.added", console.Handle)
	eventBus.Subscribe("task.created", console.Handle)
	eventBus.Subscribe("task.moved", console.Handle)
	eventBus.Subscribe("task.assigned", console.Handle)

	// 4. Repositories
	projectRepo := projectInfra.NewGormProjectRepository(db)
	taskRepo := taskInfra.NewGormTaskRepository(db)
	auditRepo := auditInfra.NewGormAuditRepository(db)

	// 4b. Audit handler — ecoute tous les events d'ecriture
	auditHandler := auditInfra.NewAuditHandler(auditRepo)
	eventBus.Subscribe("project.created", auditHandler.Handle)
	eventBus.Subscribe("member.added", auditHandler.Handle)
	eventBus.Subscribe("task.created", auditHandler.Handle)
	eventBus.Subscribe("task.moved", auditHandler.Handle)
	eventBus.Subscribe("task.assigned", auditHandler.Handle)

	// 5. Services applicatifs
	projectService := projectApp.NewProjectService(projectRepo, eventBus)
	taskService := taskApp.NewTaskService(taskRepo, eventBus)
	auditService := auditApp.NewAuditService(auditRepo)

	// 6. Handlers HTTP
	projectHandler := handlers.NewProjectHandler(projectService)
	taskHandler := handlers.NewTaskHandler(taskService)
	auditHTTPHandler := handlers.NewAuditHandler(auditService)

	// 7. Routeur
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(customMiddleware.UserContext)
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

		// Audit
		r.Get("/audit", auditHTTPHandler.GetAll)
		r.Get("/audit/{entityId}", auditHTTPHandler.GetByEntity)
	})

	// 8. Démarrage
	log.Printf("serveur démarré sur :%s", cfg.APIPort)
	if err := http.ListenAndServe(":"+cfg.APIPort, r); err != nil {
		log.Fatalf("erreur serveur: %v", err)
	}
}
