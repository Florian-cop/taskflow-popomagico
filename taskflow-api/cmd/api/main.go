package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"taskflow-api/api/v1/handlers"
	apiMiddleware "taskflow-api/api/v1/middleware"
	auditApp "taskflow-api/internal/audit/application"
	auditInfra "taskflow-api/internal/audit/infrastructure"
	notifApp "taskflow-api/internal/notification/application"
	notifDomain "taskflow-api/internal/notification/domain"
	notifInfra "taskflow-api/internal/notification/infrastructure"
	projectApp "taskflow-api/internal/project/application"
	projectInfra "taskflow-api/internal/project/infrastructure"
	realtimeApp "taskflow-api/internal/realtime/application"
	realtimeInfra "taskflow-api/internal/realtime/infrastructure"
	"taskflow-api/internal/shared/infrastructure/config"
	"taskflow-api/internal/shared/infrastructure/console"
	"taskflow-api/internal/shared/infrastructure/memory"
	"taskflow-api/internal/shared/infrastructure/persistence"
	taskApp "taskflow-api/internal/task/application"
	taskInfra "taskflow-api/internal/task/infrastructure"
	userApp "taskflow-api/internal/user/application"
	userInfra "taskflow-api/internal/user/infrastructure"
)

func main() {
	// 1. Configuration
	cfg := config.Load()

	// 2. Base de données
	db := persistence.NewDatabase(cfg.DSN())
	persistence.Migrate(db,
		&userInfra.UserModel{},
		&projectInfra.ProjectModel{},
		&projectInfra.MemberModel{},
		&taskInfra.TaskModel{},
		&notifInfra.NotificationModel{},
		&notifInfra.PreferencesModel{},
		&auditInfra.AuditLogModel{},
	)

	// 3. Event Bus + handlers console
	eventBus := memory.NewInMemoryEventBus()
	for _, name := range []string{
		"user.created",
		"project.created", "member.added",
		"task.created", "task.moved", "task.assigned",
	} {
		eventBus.Subscribe(name, console.Handle)
	}

	// 4. Repositories
	userRepo := userInfra.NewGormUserRepository(db)
	projectRepo := projectInfra.NewGormProjectRepository(db)
	taskRepo := taskInfra.NewGormTaskRepository(db)

	// 5. Adaptateurs auth (infrastructure)
	hasher := userInfra.NewBcryptPasswordHasher(0)
	tokens := userInfra.NewJWTTokenGenerator(cfg.JWTSecret, 24*time.Hour)

	// 6. Services applicatifs
	userService := userApp.NewUserService(userRepo, hasher, tokens, eventBus)
	projectService := projectApp.NewProjectService(projectRepo, eventBus)
	taskService := taskApp.NewTaskService(taskRepo, eventBus)

	// 7. Temps réel : broadcaster + event bridge branché sur le bus
	wsBroadcaster := realtimeInfra.NewWSBroadcaster()
	bridge := realtimeApp.NewEventBridge(wsBroadcaster)
	bridge.Register(eventBus,
		"task.created", "task.moved", "task.assigned",
		"member.added",
	)

	// 8. Notifications : canaux + dispatcher + event handlers
	notifRepo := notifInfra.NewGormNotificationRepository(db)
	prefsRepo := notifInfra.NewGormPreferencesRepository(db)
	channels := []notifDomain.Channel{
		notifInfra.NewEmailChannel(),
		notifInfra.NewInAppChannel(notifRepo),
	}
	dispatcher := notifApp.NewDispatcher(channels, prefsRepo)
	memberFinder := notifInfra.NewProjectMemberFinder(projectService)
	notifHandlers := notifInfra.NewEventHandlers(dispatcher, memberFinder)
	notifHandlers.Register(eventBus)
	notificationService := notifApp.NewNotificationService(notifRepo, prefsRepo)

	// 9. Audit : repository + service + event handler universel
	auditRepo := auditInfra.NewGormAuditRepository(db)
	auditService := auditApp.NewAuditService(auditRepo)
	auditHandlers := auditInfra.NewEventHandlers(auditService)
	auditHandlers.Register(eventBus,
		"user.created",
		"project.created", "member.added",
		"task.created", "task.moved", "task.assigned",
	)

	// 10. Handlers HTTP
	authHandler := handlers.NewAuthHandler(userService)
	projectHandler := handlers.NewProjectHandler(projectService)
	taskHandler := handlers.NewTaskHandler(taskService)
	wsHandler := handlers.NewWebSocketHandler(wsBroadcaster, projectService, tokens)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	auditHandler := handlers.NewAuditHandler(auditService)

	// 8. Routeur
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}))

	r.Route("/api/v1", func(r chi.Router) {
		// --- Public ---
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)

		// WS : auth via query ?token=xxx pour compat navigateur, donc pas de middleware ici
		r.Get("/projects/{id}/ws", wsHandler.HandleConnect)

		// --- Protégé par JWT ---
		r.Group(func(r chi.Router) {
			r.Use(apiMiddleware.JWTAuth(tokens))

			r.Get("/auth/me", authHandler.Me)
			r.Get("/users", authHandler.SearchUsers)
			r.Get("/users/by-email", authHandler.LookupByEmail)

			r.Get("/projects", projectHandler.GetAllProjects)
			r.Post("/projects", projectHandler.CreateProject)
			r.Get("/projects/{id}", projectHandler.GetProject)
			r.Post("/projects/{id}/members", projectHandler.AddMember)

			r.Get("/projects/{id}/tasks", taskHandler.GetTasksByProject)
			r.Post("/projects/{id}/tasks", taskHandler.CreateTask)
			r.Put("/tasks/{id}/move", taskHandler.MoveTask)

			r.Get("/notifications", notificationHandler.List)
			r.Patch("/notifications/{id}/read", notificationHandler.MarkAsRead)
			r.Get("/notifications/preferences", notificationHandler.GetPreferences)
			r.Put("/notifications/preferences", notificationHandler.UpdatePreferences)

			r.Get("/audit/logs", auditHandler.Query)
		})
	})

	// 9. Démarrage
	log.Printf("serveur démarré sur :%s", cfg.APIPort)
	if err := http.ListenAndServe(":"+cfg.APIPort, r); err != nil {
		log.Fatalf("erreur serveur: %v", err)
	}
}
