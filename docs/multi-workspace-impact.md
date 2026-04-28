# Analyse d'impact — Multi-workspace (Disruption #2, message 3)

> **Statut : analyse uniquement, pas d'implémentation.**
> Le PO demande à vérifier que l'architecture issue de la Disruption #1 peut absorber l'isolation par entreprise sans réécrire `TaskService` ni `ProjectService`.

---

## 1. Modèle métier cible

- Une entité racine **`Workspace`** (id, nom, propriétaire, créé le).
- Un **`WorkspaceMembership`** (workspaceId × userId × role : owner / admin / member).
- Un `Project` n'existe **plus dans le vide** : il appartient à un et un seul `Workspace`. Donc `Project.WorkspaceID string` (foreign key, non-null).
- Une `Task` hérite implicitement du workspace via son `Project`.

Workflow utilisateur :
1. login → `GET /workspaces` (mes workspaces)
2. choisir un workspace → token enrichi *ou* contexte HTTP qui transporte `workspaceId`
3. `GET /projects`, `GET /tasks/...` filtrés sur ce `workspaceId`

---

## 2. Couches impactées et raison

| Couche | Impact | Pourquoi |
|---|---|---|
| **Domain** `task`, `project` | **Aucun changement direct.** L'entité `Project` reçoit un champ `WorkspaceID`. Aucune méthode métier (`MoveTo`, `AddMember`, transitions) ne change. | Le workspaceId est une donnée d'**identité**, pas de logique. Les invariants restent locaux à l'agrégat. |
| **Domain** nouveau `workspace` | Création complète : entité, membership, repository (port), events `workspace.created`, `workspace.member.added`. | Nouveau bounded context, pas de dépendance vers `task`/`project`. |
| **Application** services existants | **Aucun changement de signature.** Les services lisent le `workspaceId` depuis le `context.Context` (helper `WorkspaceIDFromContext`, jumeau du `UserIDFromContext` actuel). | Évite la pollution des DTO et garde les services agnostiques de la HTTP layer. |
| **Application** repositories | Les méthodes `FindByID`, `FindAll`, `FindByProjectID` doivent **toujours** prendre en compte le workspaceId du context et filtrer. | C'est le point d'enforcement de l'isolation : un repo qui oublie le filtre = fuite de données. |
| **Infrastructure** GORM repos | Modifier les requêtes pour ajouter `WHERE workspace_id = ?` (sur projects, et indirectement sur tasks via project). | Filtre systématique au plus bas niveau, toujours appliqué. |
| **Présentation** API v1 + v2 | Nouveau **middleware `WorkspaceContext`** qui lit l'ID du workspace (header `X-Workspace-Id` ou claim `workspaceId` du JWT) et injecte dans le ctx. Vérifie aussi que l'utilisateur est membre du workspace. | Première ligne de défense : refuser toute requête sans workspace valide. |
| **Présentation** WebSocket | Le scoping doit aussi être vérifié au workspace, pas seulement au projet. | Un client connecté ne doit recevoir que les events de son workspace courant. |
| **Frontend** | Sélecteur de workspace après login + envoi systématique de l'header / sélection enregistrée en cookie. | UX : naviguer entre workspaces. |

---

## 3. Où placer le `workspaceId`

### Décision : dans le `context.Context`, alimenté par un middleware HTTP.

```go
// internal/shared/domain/context.go (extension)
func WithWorkspaceID(ctx context.Context, id string) context.Context { ... }
func WorkspaceIDFromContext(ctx context.Context) (string, bool)        { ... }
```

Le middleware (analogue à `JWTAuth`) :
- lit le workspaceId (header ou claim JWT)
- vérifie que l'utilisateur authentifié est membre du workspace (via `WorkspaceMembershipRepository`)
- injecte dans le ctx

**Pourquoi pas un argument explicite dans chaque DTO ?**
Cela imposerait de modifier la signature de **toutes** les méthodes de service (`CreateProject`, `MoveTask`, `Dispatch`, etc.), donc de toucher aux services métier. Or le PO interdit explicitement cela. Le `context.Context` est exactement le mécanisme prévu par Go pour ces préoccupations transverses (request-scoped).

**Pourquoi pas dans le JWT en dur ?**
Un utilisateur peut appartenir à plusieurs workspaces. Imposer un workspace dans le token oblige à déconnecter / reconnecter pour switcher. Le pattern qui marche bien :
- JWT contient `userId` + (optionnel) `defaultWorkspaceId`
- Header HTTP `X-Workspace-Id` override le défaut à chaque requête
- Le middleware fait l'arbitrage

---

## 4. Comment garantir l'isolation côté API

L'isolation est **enforced à la couche repository**, pas à la couche présentation.

```go
// project repository (modification minimale)
func (r *GormProjectRepository) FindAll(ctx context.Context) ([]*domain.Project, error) {
    workspaceID, ok := sharedDomain.WorkspaceIDFromContext(ctx)
    if !ok {
        return nil, sharedDomain.ErrUnauthorized
    }
    var models []ProjectModel
    err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Find(&models).Error
    // ...
}
```

**Trois lignes de défense :**
1. **Middleware** rejette les requêtes sans workspaceId valide → 401/403 avant tout traitement.
2. **Repository** applique le filtre SQL → impossible d'oublier en application code.
3. **Tests d'intégration** : un user du workspace A appelle `GET /projects/{id-du-workspace-B}` → doit retourner 404, jamais 200.

---

## 5. Pourquoi un simple filtre frontend serait insuffisant

Trois raisons concrètes :

1. **Bypass trivial via curl ou Postman.** Le frontend filtre ses appels par workspace sélectionné, mais rien n'empêche l'utilisateur d'appeler directement `GET /api/v1/projects/{id-d-un-autre-workspace}`. Si le backend ne filtre pas, les données fuient.
2. **Bug frontend = fuite de données.** Une condition de course, une variable mal scoped, et l'UI affiche les projets du mauvais workspace. La sécurité ne peut pas reposer sur la qualité du code Vue.
3. **Audit trail compromis.** L'event `task.moved` consommé par `audit` n'a aucun moyen de savoir quel workspace est concerné s'il n'est pas porté par l'event ou le contexte. Le `audit_logs` deviendrait inutilisable pour répondre à « qui a fait quoi *dans le workspace X* ? ».

L'isolation doit donc vivre **côté serveur**, et précisément au niveau du **repository**, parce que c'est le dernier point où l'on peut empêcher un SELECT cross-workspace.

---

## 6. Comment éviter de réécrire `TaskService` et `ProjectService`

Ce point est le test critique de l'architecture hexagonale livrée à la Disruption #1. La réponse est :

### Pas un seul service métier ne change.

Concrètement :
- `ProjectService.CreateProject(ctx, dto)` ne reçoit **pas** de `workspaceId` dans son DTO.
- À l'intérieur, il appelle `s.repo.Save(ctx, project)`.
- Le repo, dans `Save`, lit le `workspaceId` depuis `ctx` et **l'injecte au moment d'écrire** en BDD : `project.WorkspaceID = workspaceID; db.Create(project)`.
- L'invariant « tout projet a un workspaceId » est ainsi forcé à la couche infra, pas dans le domaine.

Pourquoi ça marche :
- Le **domaine** `project` a une nouvelle propriété `WorkspaceID string` (champ pur, sans logique). Migration Go : 1 ligne.
- Le **service** ne sait pas d'où vient le workspaceId — c'est conforme au principe de Cockburn (« le service ne sait pas s'il est appelé depuis HTTP, CLI ou test »).
- Le **repo** est l'enforcement boundary. Il était déjà l'enforcement boundary pour les autres règles d'accès BDD ; ajouter un filtre est cohérent.

### Cas particulier : la CLI et les tests
- La CLI peut soit poser le workspaceId via un flag (`--workspace`), soit utiliser un workspace par défaut. Le binaire pose `ctx = WithWorkspaceID(ctx, ...)` avant d'appeler les services.
- Les tests unitaires des services posent un workspaceId dans le ctx mock, exactement comme ils posent déjà un userId.

### Cas particulier : les events
- Les events domaine (`TaskCreatedEvent`, etc.) n'ont pas besoin de `workspaceId` dans leur payload tant qu'ils transportent déjà un `projectId` — l'audit handler peut résoudre `projectId → workspaceId` via le repo si besoin. Alternative plus propre : ajouter un champ `WorkspaceID` aux events, à décider lors de l'implémentation.

---

## 7. Plan d'implémentation (si on devait le faire — non requis pour ce rendu)

Découpage minimal, ordonné :

1. Bounded context `workspace` complet (domain + application + infra + API basic).
2. Migration : ajouter `workspace_id` sur `projects` (nullable au début), backfill, puis NOT NULL.
3. Helper `WithWorkspaceID` / `WorkspaceIDFromContext` dans `internal/shared/domain/context.go`.
4. Middleware `WorkspaceContext` (lit header, valide membership, injecte ctx).
5. Repos `project` et `task` : filtrer sur `workspace_id`.
6. Frontend : page de sélection de workspace, store cookie, intercepteur qui ajoute `X-Workspace-Id`.
7. Tests d'intégration cross-workspace (user A dans workspace 1 essaie de lire un projet du workspace 2 → 404).
8. ADR-010 — décision multi-workspace.

**Estimation honnête** : 2 jours pour un binôme expérimenté en Go + Nuxt, avec migrations BDD et tests sérieux.

---

## 8. Ce que cette analyse démontre

L'architecture issue de la Disruption #1 absorbe l'isolation multi-workspace **sans toucher** :
- aucun service métier (`ProjectService`, `TaskService`, `NotificationService`, `AuditService`)
- aucune entité du domaine `task`
- aucun handler HTTP existant (le middleware s'ajoute, les handlers ne bougent pas)
- aucun frontend existant (un sélecteur s'ajoute, les composants existants restent)

Les seules modifications **inévitables** sont :
- `Project.WorkspaceID` (1 champ ajouté, pas de logique)
- 2 méthodes de repo qui ajoutent un `WHERE`
- 1 middleware HTTP

C'est l'inverse exact de ce qui se passerait si la logique métier vivait dans les controllers : il faudrait alors propager le workspaceId dans chaque endpoint et chaque appel.
