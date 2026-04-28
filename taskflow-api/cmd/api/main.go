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
	v2handlers "taskflow-api/api/v2/handlers"
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
		&notifInfra.FailedNotificationModel{},
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
	failedRepo := notifInfra.NewGormFailedNotificationRepository(db)

	// Le canal email est wrappé dans un FaultInjectingChannel pour pouvoir simuler
	// une panne via l'API admin (chantier 1 disruption #2).
	emailToggle := notifInfra.NewFaultInjectingChannel(notifInfra.NewEmailChannel())
	channels := []notifDomain.Channel{
		emailToggle,
		notifInfra.NewInAppChannel(notifRepo),
	}
	dispatcher := notifApp.NewDispatcher(channels, prefsRepo, failedRepo)
	memberFinder := notifInfra.NewProjectMemberFinder(projectService)
	notifHandlers := notifInfra.NewEventHandlers(dispatcher, memberFinder)
	notifHandlers.Register(eventBus)
	notificationService := notifApp.NewNotificationService(notifRepo, prefsRepo)
	adminService := notifApp.NewAdminService(dispatcher, failedRepo, emailToggle)

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
	adminHandler := handlers.NewAdminHandler(adminService)

	// Handlers v2 — consomment les MÊMES services applicatifs que v1.
	v2ProjectHandler := v2handlers.NewProjectHandler(projectService)
	v2TaskHandler := v2handlers.NewTaskHandler(taskService)

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

			// --- Admin (chantier 1 disruption #2 — résilience) ---
			r.Get("/admin/notifications/channels", adminHandler.ListChannels)
			r.Put("/admin/notifications/channels/{name}", adminHandler.SetChannelFailing)
			r.Get("/admin/notifications/failed", adminHandler.ListFailed)
			r.Post("/admin/notifications/failed/{id}/retry", adminHandler.RetryFailed)
		})
	})

	// --- API v2 (chantier 2 disruption #2) ---
	// Coexistence stricte avec v1 : mêmes services métier, présentation différenciée.
	r.Route("/api/v2", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(apiMiddleware.JWTAuth(tokens))

			r.Get("/projects", v2ProjectHandler.GetAll)
			r.Post("/projects", v2ProjectHandler.Create)
			r.Get("/projects/{id}", v2ProjectHandler.Get)
			r.Post("/projects/{id}/members", v2ProjectHandler.AddMember)

			r.Get("/projects/{id}/tasks", v2TaskHandler.ListByProject)
			r.Post("/projects/{id}/tasks", v2TaskHandler.Create)
			r.Put("/tasks/{id}/move", v2TaskHandler.Move)
		})
	})

	// 9. Démarrage
	log.Printf("serveur démarré sur :%s", cfg.APIPort)
	if err := http.ListenAndServe(":"+cfg.APIPort, r); err != nil {
		log.Fatalf("erreur serveur: %v", err)
	}
}
