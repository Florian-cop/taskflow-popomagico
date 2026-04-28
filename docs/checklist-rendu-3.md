# Checklist d'évaluation — Rendu 3

> Synthèse exhaustive de toutes les exigences (Phase 1, Disruption #1, Disruption #2, règles README) confrontées à l'état du code à ce jour.

Légende : ✅ implémenté · ⚠️ partiel ou raccourci documenté · ❌ non fait · 📄 doc seule

---

## A. Phase 1 — Fondations (Rendu 1)

| Exigence | État | Preuve / fichiers |
|---|---|---|
| Module `project` (entité, service, interface repo, impl ORM) | ✅ | `internal/project/{domain,application,infrastructure}` |
| Module `task` (entité, service, interface repo, impl ORM) | ✅ | `internal/task/{domain,application,infrastructure}` |
| Value Object `TaskStatus` avec transitions valides (todo → in_progress → done) | ✅ | `internal/task/domain/task_status.go` |
| Au moins 2 domain events publiés (`task.created`, `task.moved`) | ✅ + bonus | `task.assigned`, `project.created`, `member.added`, `user.created` aussi |
| ConsoleHandler abonné aux events | ✅ | `internal/shared/infrastructure/console/handler.go` |
| Frontend Kanban avec déplacement de tâche | ✅ | `taskflow-web/app/components/{KanbanBoard,KanbanColumn,TaskCard}.vue` |
| Tests unitaires sans BDD | ✅ | `internal/project/application/*_test.go`, idem task ; `mock_event_bus.go` |
| Identifiant utilisateur simulé (Phase 1) | ✅ remplacé | Originellement `X-User-Id`, maintenant remplacé par JWT (Disruption #1) |
| Procédure de démarrage documentée | ✅ | `docs/demarrage.md` |
| Min. 3 ADR | ✅ | `docs/ADR-001` à `ADR-004` (et 005 à 009 pour Disruption #1) |

---

## B. Règles permanentes du README

| Règle | État | Preuve |
|---|---|---|
| 1. Un ADR par décision technique significative | ✅ | 9 ADR dans `docs/` (000 à 009) |
| 2. Aucune logique métier dans les controllers | ✅ | Handlers HTTP très minces : decode → service → encode |
| 3. Aucun accès direct ORM depuis les services | ✅ | Services dépendent uniquement d'interfaces `*Repository` |
| 4. Tests unitaires sans BDD | ✅ | Mock repos + `MockEventBus` (`internal/shared/application/mock_event_bus.go`) |
| 5. `docker compose up` depuis un clone propre | ✅ | `cp .env.example .env && docker compose up --build` (frontend + api + cli + postgres) |
| 6. Commits réguliers | ✅ | `git log` de l'équipe |
| 7. Tags `rendu-1`, `rendu-2`, `rendu-3` | ⚠️ | `rendu-2` créé en local, `rendu-3` à créer après ce rendu |

---

## C. Disruption #1 — Évolution

### C.1 Authentification JWT
| Exigence | État | Preuve |
|---|---|---|
| Inscription / connexion / token | ✅ | `POST /api/v1/auth/register`, `POST /api/v1/auth/login`, `GET /auth/me` |
| Migration SSO sans toucher la logique métier | ✅ | Interfaces `TokenGenerator` + `PasswordHasher` dans `internal/user/domain/` ; aucun import JWT/bcrypt dans `internal/task` ou `internal/project` (vérifié par grep) |
| Documentation du choix | ✅ | `docs/ADR-005.md` |

### C.2 Temps réel Kanban
| Exigence | État | Preuve |
|---|---|---|
| Mise à jour instantanée sans recharger | ✅ | `internal/realtime/` + `pages/projects/[id].vue` + `useRealtime.ts` |
| Scoping par projet | ✅ | `WSBroadcaster.rooms` indexée par `projectID` (`internal/realtime/infrastructure/ws_broadcaster.go`) ; check membership dans le handler |
| Techno interchangeable sans impact métier | ✅ | Interface `Broadcaster` (1 méthode) ; `task_service`/`project_service` ne référencent jamais le bounded context realtime |
| Documentation + alternatives | ✅ | `docs/ADR-006.md` |

### C.3 Notifications événementielles
| Exigence | État | Preuve |
|---|---|---|
| Notifs sur `task.assigned` (assignee) et `task.moved` (membres du projet) | ✅ | `internal/notification/infrastructure/event_handlers.go` |
| Canaux email + in-app | ✅ | `email_channel.go` (log) + `in_app_channel.go` (BDD) |
| Préférences utilisateur via API (toggle email/in-app) | ✅ | `GET/PUT /api/v1/notifications/preferences` + UI `pages/settings/notifications.vue` |
| Ajouter Slack/Teams/SMS sans toucher les services métier | ✅ | Strategy `Channel` interface ; ajouter un canal = 1 fichier infra + 1 ligne dans `main.go` |
| Documentation extensibilité | ✅ | `docs/ADR-007.md` |
| Email simulé par log | ⚠️ assumé | Documenté dans ADR-007 et `disruption-1-impact.md` |

### C.4 Audit trail
| Exigence | État | Preuve |
|---|---|---|
| Trace des écritures sur Task et Project | ✅ | `internal/audit/infrastructure/event_handlers.go` abonné à 6 events |
| Sans pollution de la logique métier | ✅ | Pur consommateur d'events ; `grep -r "audit" internal/task internal/project` → 0 résultat |
| Backend de stockage interchangeable | ✅ | Interface `AuditRepository` ; impl unique `GormAuditRepository` |
| Documentation isolation | ✅ | `docs/ADR-008.md` |

### C.5 CLI d'administration
| Exigence | État | Preuve |
|---|---|---|
| Créer projet | ✅ | `cli project create --name X --owner Y` |
| Créer tâche | ✅ | `cli task create --project ID --title T` |
| Jeu de démo (TODO/IN_PROGRESS/DONE) | ✅ | `cli generate-demo` → 9 tâches réparties |
| Pas de duplication de logique | ✅ | `cmd/cli/main.go` importe `internal/project/application` et `internal/task/application` |
| Documentation interface | ✅ | `docs/ADR-009.md`, `docs/demarrage.md` |

### C.6 Mise en production initiale
| Exigence | État | Preuve |
|---|---|---|
| `docker compose up` depuis clone propre | ✅ | Stack complète (web + api + postgres) |
| `.env` non versionné, `.env.example` documenté | ✅ | `.gitignore` ignore `.env` ; `.env.example` racine + `taskflow-api/.env.example` |
| Pas d'installation manuelle | ✅ | Tout via Docker, healthcheck postgres avant démarrage api |

---

## D. Disruption #2 — Résilience

### D.1 Résilience des consommateurs de notifications
| Exigence | État | Preuve |
|---|---|---|
| Un canal en panne ne casse pas les autres | ✅ | `Dispatcher.Dispatch` itère et swallow chaque erreur de canal (`internal/notification/application/dispatcher.go`) |
| Aucune erreur visible côté utilisateur | ✅ | `Dispatch` retourne `nil` côté API ; les events restent consommés |
| Échec du canal email stocké pour retraitement | ✅ | Table `failed_notifications` + `Dispatcher.recordFailure` |
| Démontrer la simulation d'exception EmailChannel | ✅ | `FaultInjectingChannel` + `PUT /admin/notifications/channels/email` + UI `/admin/notifications` |
| Retraitement manuel | ✅ | `POST /admin/notifications/failed/{id}/retry` + bouton "Rejouer" |
| InApp / Realtime / Audit continuent malgré la panne email | ✅ | Vérifiable via démo : déplacer une tâche avec email en panne → toast UI passe, in-app notif arrive, WS notifie, audit_log s'incrémente |

### D.2 Versioning API
| Exigence | État | Preuve |
|---|---|---|
| `/api/v1/...` continue de fonctionner | ✅ | Routes inchangées, clients existants (frontend) toujours compatibles |
| `/api/v2/...` expose les mêmes cas d'usage | ✅ | `api/v2/handlers/{project,task}_handler.go` |
| Format de réponse différencié | ✅ | v2 enveloppe systématique `{ data, meta: { apiVersion, generatedAt, count? } }` |
| Réutilisation de la logique métier | ✅ | Les v2 handlers appellent `projectService` et `taskService` — exactement les mêmes instances que v1 |
| Démontrer que l'API est une couche d'entrée distincte du domaine | ✅ | Aucun fichier sous `internal/` ne référence `v1` ou `v2`, ni l'enveloppe |

### D.3 Multi-workspace (analyse seulement)
| Exigence | État | Preuve |
|---|---|---|
| Quelles couches modifiées | 📄 | `docs/multi-workspace-impact.md` §2 |
| Où placer le `workspaceId` | 📄 | §3 — context.Context via middleware |
| Comment garantir l'isolation côté API | 📄 | §4 — middleware + repo, trois lignes de défense |
| Pourquoi un filtre frontend insuffisant | 📄 | §5 — bypass curl, bug front, audit compromis |
| Comment éviter de réécrire les services | 📄 | §6 — aucun service ne change, repos appliquent le filtre |

### D.4 Critères Rendu 3 (README)
| Exigence | État | Preuve |
|---|---|---|
| ADR pour chaque choix important | ⚠️ | ADR-005 à 009 existent ; à compléter par ADR-010 (résilience) et ADR-011 (versioning) avant figeage |
| Tableau scénarios de panne | ❌ | À créer dans `docs/scenarios-panne.md` |
| Tag `rendu-3` | ❌ | À créer après checklist OK et tag pushé |
| Démo live `docker compose up` | ✅ | Validé : démarre sans config manuelle |
| Simuler une panne et prouver que le système continue | ✅ | Voir D.1 |
| Montrer deux versions de l'API coexister | ✅ | Voir D.2 |

---

## E. Trous identifiés à combler avant le Rendu 3

1. **ADR-010 — Résilience des notifications** — formaliser le décorateur `FaultInjectingChannel`, la dead-letter `failed_notifications`, le retry manuel synchrone, la décision de ne pas faire d'auto-retry pour cette livraison.
2. **ADR-011 — Versioning API par URL** — formaliser le choix `/api/v1` vs `/api/v2`, comparer aux alternatives (header `X-API-Version`, media type), expliquer la non-duplication des services.
3. **`docs/scenarios-panne.md`** — tableau exigé par le README Rendu 3 (situation, comportement attendu, comportement observé). À minima : panne email, panne BDD pendant un retry, déconnexion WS pendant déplacement de tâche.
4. **`docs/disruption-2-impact.md`** — analogue à `disruption-1-impact.md` pour la D2 (fichiers stables, ajoutés, modifiés).
5. **Tag `rendu-3`** — une fois les 4 points ci-dessus terminés.
6. **Frontend** : pas d'écran v2 — non requis par le PO, mais pourrait densifier la démo en montrant la même page Kanban consommant `/api/v2/...` via un toggle. Optionnel.

---

## F. Score d'auto-évaluation

- Phase 1 : 100 %
- Disruption #1 : 100 % du périmètre demandé, 2 raccourcis explicitement assumés (email = log, JWT HS256 secret partagé)
- Disruption #2 : 95 % — manque les 4 éléments de section E
- Règles transverses : 90 % — manque ADR-010, ADR-011, scénarios-panne, tag rendu-3

Tout le code compile (`go build ./...`, `pnpm exec nuxt typecheck`), les tests existants passent (`go test ./...`).
