# Analyse d'impact — Disruption #2

> Diff entre le tag `rendu-2` et la livraison actuelle (Rendu 3).

---

## Fichiers stables (aucune modification)

Le test critique de l'architecture issue de la Disruption #1 :

- `internal/project/domain/*` — entités, repository (port), events
- `internal/project/application/project_service.go`
- `internal/project/infrastructure/*`
- `internal/task/domain/*`
- `internal/task/application/task_service.go`
- `internal/task/infrastructure/*`
- `internal/user/*` (entier)
- `internal/audit/*` (entier — résilience ne le concerne pas, versioning non plus)
- `internal/realtime/*` (entier)
- `internal/shared/*` (event bus, context, persistence)
- Toute l'API v1 sous `api/v1/*` reste compatible byte-for-byte avec rendu-2

**Aucune logique métier ni aucun service applicatif existant n'a été modifié pour livrer la Disruption #2.**

---

## Fichiers ajoutés

### Chantier 1 — Résilience des notifications (ADR-010)

Backend :
- `internal/notification/domain/failed_notification.go` — entité `FailedNotification` + status pending/retried
- `internal/notification/infrastructure/fault_injecting_channel.go` — décorateur `Channel` avec `SetFailing/IsFailing`
- `internal/notification/infrastructure/gorm_failed_notification_repository.go` — table `failed_notifications`
- `internal/notification/application/admin_service.go` — service admin (list/retry/toggle)
- `internal/notification/application/errors.go` — `ErrChannelUnknown` (ajout)
- `api/v1/handlers/admin_handler.go` — handlers admin

Frontend :
- `taskflow-web/app/composables/useAdmin.ts`
- `taskflow-web/app/pages/admin/notifications.vue`

### Chantier 2 — Versioning API (ADR-011)

- `api/v2/dto/envelope.go` — structure `Response` avec `data` + `meta`
- `api/v2/handlers/common.go` — helpers `writeEnveloped` / `writeListEnveloped`
- `api/v2/handlers/project_handler.go` — handlers v2 réutilisant `ProjectService`
- `api/v2/handlers/task_handler.go` — handlers v2 réutilisant `TaskService`

### Chantier 3 — Multi-workspace (analyse uniquement)

- `docs/multi-workspace-impact.md` — analyse complète, pas de code

### Documentation transverse

- `docs/ADR-010.md` — résilience
- `docs/ADR-011.md` — versioning API
- `docs/scenarios-panne.md` — tableau exigé par le README
- `docs/checklist-rendu-3.md` — exigences vs réalité
- `docs/plan-demo.md` — plan de démo 5 minutes

---

## Fichiers modifiés

| Fichier | Raison |
|---|---|
| `internal/notification/domain/repository.go` | Ajout de l'interface `FailedNotificationRepository` |
| `internal/notification/application/dispatcher.go` | Ajout du paramètre `failedRepo`, méthodes `RetryFailed`, `Channels`, `recordFailure` |
| `cmd/api/main.go` | Wiring du décorateur `FaultInjectingChannel`, du `failedRepo`, de l'`AdminService`, des handlers v2, des routes admin et `/api/v2` |

Hors de ces 3 fichiers, **aucun code métier n'a été touché**.

---

## Vérification des contraintes du PO (Disruption #2)

| Exigence | Preuve |
|---|---|
| Un canal en panne ne casse pas les autres | `internal/notification/application/dispatcher.go` ligne 44-47 : la boucle continue après `c.Send` qui échoue ; l'erreur n'est pas remontée |
| Aucune erreur visible côté utilisateur | `Dispatch` retourne `nil` ; `MoveTask` côté API retourne 200 ; UI ne reçoit jamais d'erreur de notif |
| Stockage des messages échoués | Table `failed_notifications`, ligne créée à chaque échec via `recordFailure` |
| Démontrer la simulation d'exception sur EmailChannel | `FaultInjectingChannel` + UI `/admin/notifications` ou `PUT /admin/notifications/channels/email` |
| `/api/v1` continue de fonctionner | Code v1 inchangé, frontend v1 sans modification, tests existants passent |
| `/api/v2` expose les mêmes cas d'usage | `api/v2/handlers/{project,task}_handler.go` couvre projects, members, tasks |
| Format de réponse différencié | v1 = JSON plat ; v2 = `{data, meta:{apiVersion,generatedAt,count?}}` |
| Réutilisation de la logique métier | v2 handlers sont instanciés avec **les mêmes** `projectService` et `taskService` que v1 (lignes 113 et 121 de `cmd/api/main.go`) |
| Analyse multi-workspace | `docs/multi-workspace-impact.md` — 8 sections |

---

## Test des principes architecturaux (grep)

```bash
# Le domaine ne doit rien savoir de la résilience ou du versioning :
grep -r "fault\|failed_notification\|FaultInjecting" \
  taskflow-api/internal/task taskflow-api/internal/project \
  taskflow-api/internal/audit taskflow-api/internal/realtime
# → 0 résultat

grep -r "v2\|api/v2\|envelope" \
  taskflow-api/internal/
# → 0 résultat

# Les handlers v2 ne dépendent que de la couche application :
grep "import" taskflow-api/api/v2/handlers/*.go | grep "internal" | sort -u
# → seulement internal/project/application, internal/task/application, internal/shared/domain
```

---

## Raccourcis assumés (ajout à ceux de la Disruption #1)

| Raccourci | Justification | Remplacement en production |
|---|---|---|
| Retry des notifications synchrone et manuel | Pas de scheduler ni de worker dans le pilote | Worker NATS/Kafka avec back-off exponentiel et idempotence |
| Toggle de panne email accessible à tout user authentifié | Pas de RBAC dans le pilote | Rôle `admin` avec policy + audit des actions admin |
| Versioning v2 limité à projects/tasks/members | Suffisant pour démontrer le pattern | Étendre v2 aux notifs/audit/admin si un client v2 le demande |
| Pas de header `Sunset` sur v1 | v1 n'est pas dépréciée | Ajouter quand la deprecation policy sera décidée |
| `FaultInjectingChannel` à état mémoire | Suffit pour la démo, simple | Configuration externalisée (Redis ou feature flags) pour panne distribuée |

---

## Ce que cette analyse démontre

L'architecture hexagonale livrée à la Disruption #1 absorbe la Disruption #2 **sans toucher** :
- aucun service métier
- aucune entité du domaine
- aucun handler v1 existant
- aucun composable frontend existant (l'admin est une nouvelle page autonome)

Les seules modifications **inévitables** sont en couche infra/application :
- `Dispatcher` : un nouveau paramètre, une nouvelle méthode, le swallow d'erreur déjà en place
- `notification/domain/repository.go` : une interface ajoutée
- `cmd/api/main.go` : du câblage

C'est le critère exact que le PO voulait vérifier : que la disruption peut être absorbée par **ajout d'adaptateurs et de handlers**, sans réécriture du domaine.
