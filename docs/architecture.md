# Architecture — TaskFlow

Le projet est un monorepo avec le backend Go et le frontend Nuxt côte à côte. Chaque sous-projet a son propre Dockerfile et se build indépendamment.

## Structure globale

```
taskflow-popomagico/
├── taskflow-api/           ← backend Go (DDD + Hexagonal + Event-Driven)
├── taskflow-web/           ← frontend Nuxt (Vue.js + Nuxt UI)
├── docs/                   ← ADR, architecture, glossaire
├── docker-compose.yml      ← orchestration : api, web, postgres
├── .github/workflows/      ← CI GitHub Actions
└── .env.example
```

## Backend — `taskflow-api/`

Le code est organisé par **bounded context** (DDD). Chaque contexte contient trois couches : domaine, application, infrastructure.

```
taskflow-api/
├── cmd/api/main.go                  ← point d'entrée, câblage des dépendances
├── internal/
│   ├── project/                     ← contexte "Project"
│   │   ├── domain/                  ← entités, value objects, interfaces, events
│   │   ├── application/             ← use cases, DTOs, tests
│   │   └── infrastructure/          ← adaptateur GORM, modèles, mappers
│   ├── task/                        ← contexte "Task" (même structure)
│   ├── notification/                ← contexte "Notification"
│   ├── audit/                       ← contexte "Audit"
│   └── shared/                      ← interfaces transversales (EventBus, DomainEvent),
│                                      adaptateurs partagés (DB, HTTP, config)
├── api/
│   └── v1/                          ← handlers HTTP, DTOs requête/réponse, mappers
│       ├── handlers/
│       ├── dto/
│       └── mapper/
├── go.mod
└── Dockerfile
```

**Ce que contient chaque couche :**

| Couche | Contenu | Règle |
|--------|---------|-------|
| `domain/` | Entités, value objects, interfaces repository (ports), domain events | Pur Go, aucune dépendance externe |
| `application/` | Services (use cases), DTOs, tests unitaires | Dépend uniquement de `domain/` |
| `infrastructure/` | Implémentations GORM, modèles BDD, mappers | Implémente les interfaces de `domain/` |

## Frontend — `taskflow-web/`

Le frontend est un adaptateur entrant : il consomme l'API REST, pas de logique métier dedans.

```
taskflow-web/
├── app/
│   ├── pages/                  ← routes (index, projects/[id])
│   ├── components/
│   │   ├── kanban/             ← KanbanBoard, KanbanColumn, KanbanCard
│   │   ├── project/            ← ProjectList, ProjectForm
│   │   └── task/               ← TaskForm, TaskDetail
│   ├── composables/            ← useProjects, useTasks, useApi
│   └── layouts/
├── nuxt.config.ts
├── package.json
└── Dockerfile
```

## Règles de dépendance

Les dépendances pointent toujours vers le centre (le domaine). Jamais l'inverse.

```
  Présentation (api/v1/)     →  Application  →  Domain  ←  Infrastructure
  (handlers HTTP)               (use cases)     (pur Go)    (GORM)
```

1. `domain/` n'importe jamais les autres couches
2. `application/` n'importe jamais `infrastructure/` ni `api/`
3. `infrastructure/` implémente les interfaces de `domain/`
4. `api/` appelle `application/` et utilise les types de `domain/`
5. Seul `cmd/api/main.go` connaît tout — c'est là que se fait l'injection de dépendances

## Flux événementiel

Quand un use case fait quelque chose de significatif, il publie un event sur le bus interne de l'application. Les handlers consomment ces events de manière indépendante et découplée.

```
Requête HTTP → Handler → Service → persist en BDD → publie event → Bus interne
                                                                  ↓
                                              NotificationHandler / AuditHandler / ConsoleHandler
```

Un handler qui plante ne fait pas échouer le use case. Les retries sont gérés dans l'application (cf. ADR-003).

## Docker

3 services dans le `docker-compose.yml` :

| Service | Image | Port |
|---------|-------|------|
| `api` | Go multi-stage build | 8080 |
| `web` | Node.js Alpine | 3000 |
| `postgres` | postgres:16-alpine | 5432 |
