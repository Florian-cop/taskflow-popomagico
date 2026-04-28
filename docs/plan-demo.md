# Plan de démonstration live — Rendu 3 (5 minutes)

> Objectif : convaincre le jury que **l'architecture absorbe les disruptions** sans réécriture du domaine.
> Format : démo live, 1 narrateur + 1 conducteur clavier (binôme), pas de slides.

## 0. Pré-démo (à faire avant de présenter)

```bash
# Repo propre, base vide
git checkout main
docker compose down -v
docker compose up --build -d
# attendre l'API : docker compose logs -f api → "serveur démarré sur :8080"

# Préparer 2 onglets navigateur incognito + 1 onglet normal
# Onglet 1 : http://localhost:3000 — Alice
# Onglet 2 : http://localhost:3000 — Bob (privé)
# Onglet 3 : terminal pour curl + websocat

# Préparer dans des post-it ou un éditeur :
#   ALICE_TOKEN=, BOB_TOKEN=, PROJECT_ID=, TASK_ID=
```

---

## 1. (30 s) — Démarrage from scratch

**Narration :**
> « Le PO exige un démarrage en une commande depuis un clone propre. »

**Action :**
```bash
docker compose down -v && docker compose up --build
```

**Preuve montrée :** logs `[Event] xxx` qui défilent quand on commence à interagir, healthcheck postgres OK.

---

## 2. (45 s) — Auth JWT (Disruption #1, chantier 1)

**Narration :**
> « Plus de header `X-User-Id` simulé : vraie inscription, JWT HS256, middleware dédié. Le domaine `task` et `project` n'importent pas une seule ligne de JWT — vérifiable par `grep`. »

**Actions :**
1. Onglet 1 : `/login` → onglet **Inscription** → créer Alice (`alice@taskflow.io`)
2. Onglet 2 : créer Bob (`bob@taskflow.io`)
3. **Pendant ce temps**, terminal :
   ```bash
   grep -r "jwt\|golang-jwt" taskflow-api/internal/task taskflow-api/internal/project
   # → aucun résultat
   ```

**Punch :** « Migrer en SSO = écrire une nouvelle implémentation de `TokenGenerator`. Zéro ligne du domaine ne bouge. »

---

## 3. (1 min) — Temps réel + Notifications + Audit (Disruption #1, chantiers 2-4)

**Narration :**
> « Une seule action métier — déplacer une tâche — déclenche **trois consommateurs découplés** abonnés au même event bus. »

**Actions :**
1. Alice crée un projet « Démo Rendu 3 ».
2. Sur la page projet, **Ajouter un membre** → autocomplete `bob` → sélectionner Bob.
3. Bob refresh → il voit le projet apparaître.
4. Alice crée une tâche « Préparer la soutenance ».
5. **Côté Bob** : la tâche apparaît instantanément (WS).
6. Alice clique « In Progress → ».
7. **Côté Bob** : la tâche bouge en temps réel + cloche notif passe à `1` (in-app).
8. Bob ouvre **Journal d'audit** dans le menu utilisateur → 4 entrées (project.created, member.added, task.created, task.moved).

**Punch :** « Trois bounded contexts, chacun consommateur du bus, chacun isolé du domaine — `grep -r "audit" internal/task` = 0. »

---

## 4. (1 min 15) — Résilience (Disruption #2, chantier 1)

**Narration :**
> « Hier matin, l'API email est tombée. Conséquence attendue : aucune. »

**Actions :**
1. Menu utilisateur Alice → **Admin résilience** (`/admin/notifications`).
2. Toggle **Simuler une panne** sur le canal `email` → ON.
3. Observer le toast « Panne simulée sur email ».
4. Alice retourne sur le projet, clique « In Progress → Done » sur une tâche.
5. **Vérifier** :
   - L'UI confirme « Task moved » sans erreur visible (Dispatcher avale l'erreur).
   - Bob voit la tâche bouger en temps réel (WS) → **Realtime continue**.
   - Cloche de Bob s'incrémente (in-app) → **InApp continue**.
   - Audit montre la nouvelle entrée → **Audit continue**.
6. Retour Admin → la liste **Messages échoués** affiche maintenant un message email avec l'erreur `simulated channel failure`.
7. Toggle OFF la panne email.
8. Cliquer **Rejouer** → toast « Notification rejouée », message disparaît de la liste.

**Punch :** « Le canal défaillant est isolé, persisté pour retraitement, et les trois autres consommateurs ne savent même pas qu'il y a eu un problème. »

---

## 5. (50 s) — Versioning API (Disruption #2, chantier 2)

**Narration :**
> « Un partenaire mobile ne peut pas suivre nos évolutions. Solution : `/api/v2` coexiste avec `/api/v1`, **sans dupliquer un seul service métier**. »

**Actions** (terminal, à projeter) :
```bash
# v1 : retour brut
curl -s -H "Authorization: Bearer $ALICE_TOKEN" \
  http://localhost:8080/api/v1/projects | jq
# → [ { "id": "...", "name": "Démo Rendu 3", ... } ]

# v2 : même donnée, présentation enveloppée
curl -s -H "Authorization: Bearer $ALICE_TOKEN" \
  http://localhost:8080/api/v2/projects | jq
# → { "data": [ { ... } ], "meta": { "apiVersion": "v2", "generatedAt": "...", "count": 1 } }

# Preuve de réutilisation : créer un projet via v2
curl -s -X POST -H "Authorization: Bearer $ALICE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Via v2","description":"Test"}' \
  http://localhost:8080/api/v2/projects | jq

# Le frontend (qui tape v1) voit immédiatement le projet apparaître → MÊME service métier
```

**Punch :** « `api/v2/handlers/` ne dépend que de `internal/{project,task}/application` — la couche présentation change, le domaine est immobile. »

---

## 6. (45 s) — CLI réutilise le domaine (Disruption #1, chantier 5)

**Narration :**
> « Le PO veut une CLI pour générer un jeu de démo. Sans dupliquer la logique. »

**Actions :**
```bash
docker compose exec api /cli generate-demo
```

**Preuve :**
- 1 projet + 9 tâches apparaissent dans l'UI (refresh).
- Logs `api` : 9× `[Event] task.created`, 3× `[Event] task.moved`, etc.
- `cmd/cli/main.go` importe **exactement** les mêmes services applicatifs que `cmd/api/main.go`.

**Punch :** « La CLI est un adaptateur entrant comme l'API REST. Aucun copier-coller. »

---

## 7. (40 s) — Multi-workspace : démo de l'analyse (Disruption #2, chantier 3)

**Narration :**
> « Le PO veut isoler les données par entreprise. On nous demande l'analyse, pas l'implémentation. »

**Action :**
- Ouvrir `docs/multi-workspace-impact.md` côté écran.
- Lire à voix haute la section §6 :
  > « **Pas un seul service métier ne change.** `ProjectService.CreateProject(ctx, dto)` ne reçoit pas de workspaceId dans son DTO. Le repo lit le workspaceId depuis ctx au moment d'écrire en BDD. »

**Punch :** « C'est exactement le test de notre architecture. Si le PO nous demandait de coder ça lundi, on toucherait : 1 champ sur `Project`, 2 méthodes de repo, 1 middleware. Zéro service métier. »

---

## 8. (15 s) — Conclusion

**Narration :**
> « Trois disruptions absorbées en ajoutant des **adaptateurs** et des **consommateurs d'events**. Le domaine `task` et `project` n'a pas changé d'une ligne entre le rendu 1 et aujourd'hui. C'est ce que vise l'architecture hexagonale. »

**Démontrer :**
```bash
git diff rendu-1..HEAD -- taskflow-api/internal/task/domain taskflow-api/internal/project/domain
# → diff minimal voire vide hors fix de typage
```

---

## Annexe A — Commandes utiles pendant la démo

```bash
# Tokens (copier la sortie, filtre jq pour récupérer)
ALICE_TOKEN=$(curl -s -X POST localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"alice@taskflow.io","password":"<password>"}' | jq -r .token)

# Forcer un échec, pour la démo si l'UI buggue :
curl -X PUT -H "Authorization: Bearer $ALICE_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"failing":true}' \
  localhost:8080/api/v1/admin/notifications/channels/email

# Lister les failed
curl -s -H "Authorization: Bearer $ALICE_TOKEN" \
  localhost:8080/api/v1/admin/notifications/failed | jq

# Reset complet
docker compose down -v && docker compose up -d
```

## Annexe B — Plan B si une étape rate

| Étape | Plan B |
|---|---|
| Frontend ne charge pas | Tout démontrer au curl + websocat |
| WS ne se connecte pas | Montrer les logs API qui prouvent que les events sont publiés ; expliquer que le bus continue même sans consommateur WS |
| Email toggle ne s'active pas | Le faire au curl directement (cf. annexe A) |
| `generate-demo` plante | Créer 2-3 tâches manuellement via UI, montrer les events dans les logs |

## Annexe C — Fichiers à avoir ouverts en arrière-plan

- `docs/multi-workspace-impact.md` (étape 7)
- `docs/checklist-rendu-3.md` (référence si question)
- `taskflow-api/internal/notification/application/dispatcher.go` (étape 4 si question sur la résilience)
- `taskflow-api/api/v2/handlers/project_handler.go` (étape 5 si question sur la duplication)
