# Scénarios de panne — Rendu 3

> Tableau exigé par le README pour le Rendu 3 : pour chaque scénario, la situation, le comportement attendu, le comportement observé, et le mécanisme d'absorption qui le rend possible.

Chaque scénario est reproductible localement. Les commandes utilisent les conventions de `docs/demarrage.md`.

---

## Scénario 1 — Le canal email tombe en pleine journée

| Situation | Un utilisateur déplace une tâche. Le canal email est en panne (API tierce HS). |
|---|---|
| Comportement attendu | Le déplacement réussit. Les autres canaux (`in_app`) fonctionnent. Le WebSocket diffuse l'event. L'audit log enregistre l'opération. L'utilisateur ne voit aucune erreur. |
| Comportement observé | ✅ Conforme. Le `Dispatcher` swallow l'erreur du canal email, persiste un `failed_notifications`, et continue avec le canal `in_app`. L'event `task.moved` est consommé indépendamment par `EventBridge` (WS) et `AuditHandlers` qui ne dépendent pas du canal email. |
| Mécanisme d'absorption | `Dispatcher` n'expose jamais une erreur de canal au caller. Les consommateurs `Realtime`, `Audit` sont abonnés directement au bus, pas au canal email. Le décorateur `FaultInjectingChannel` permet de simuler la panne sans déployer. |
| Reproduction | `PUT /api/v1/admin/notifications/channels/email` body `{"failing":true}` puis déplacer une tâche. Vérifier `[EMAIL]` absent des logs, `[Event] task.moved` présent, `failed_notifications` contient une nouvelle ligne. |

---

## Scénario 2 — Le canal email est rétabli, on rejoue les messages en attente

| Situation | Suite au scénario 1, le canal email est de nouveau opérationnel. Plusieurs messages sont en attente. |
|---|---|
| Comportement attendu | Un opérateur peut consulter la liste des messages échoués et les rejouer un par un. Les messages rejoués avec succès passent en status `retried` et disparaissent de la liste pending. |
| Comportement observé | ✅ Conforme. `GET /api/v1/admin/notifications/failed` liste les messages pending. `POST /api/v1/admin/notifications/failed/{id}/retry` retente le canal d'origine. La méthode `Dispatcher.RetryFailed` met à jour le compteur `RetryCount` et le statut. |
| Mécanisme d'absorption | Persistence des échecs dans `failed_notifications` (table relationnelle, survie au restart). API admin pour orchestrer manuellement. UI `/admin/notifications` pour les utilisateurs non-techniques. |
| Reproduction | Créer 2-3 échecs (scénario 1), repasser le canal en `failing:false`, rejouer chaque message via l'UI. Vérifier les logs `[EMAIL] to=...` qui réapparaissent. |
| Limite assumée | Le retry est synchrone et manuel. Pas d'auto-retry avec back-off exponentiel — voir ADR-010. |

---

## Scénario 3 — Le client mobile (v2) appelle l'API pendant qu'un déploiement met v1 à jour

| Situation | Le frontend interne consomme `/api/v1/projects`. Un client mobile externe consomme `/api/v2/projects`. Une mise à jour modifie le format de réponse v1. |
|---|---|
| Comportement attendu | Le client mobile continue de fonctionner sans interruption. Le frontend interne se met à jour pour suivre le nouveau format v1. |
| Comportement observé | ✅ Conforme par construction. Les deux versions sont des **packages séparés** (`api/v1/handlers`, `api/v2/handlers`). Modifier l'enveloppe v1 ne touche pas le code v2. Les deux passent par les mêmes services applicatifs, donc la logique métier reste cohérente. |
| Mécanisme d'absorption | Versioning par URL (ADR-011). Coexistence stricte. Aucun service métier ne sait quelle version l'a appelé. |
| Reproduction | `curl /api/v1/projects` puis `curl /api/v2/projects` — comparer la structure (flat vs `{data, meta}`). Modifier `api/v1/dto/project.go` (ex: renommer `members` en `users`) → `curl v1` change, `curl v2` est identique. |

---

## Scénario 4 — Un client se connecte au WebSocket d'un projet dont il n'est pas membre

| Situation | User A envoie un GET `/api/v1/projects/{id-du-projet-de-B}/ws?token=<son-token>`. |
|---|---|
| Comportement attendu | La connexion WS est refusée. Aucun event du projet ne lui est diffusé. |
| Comportement observé | ✅ Conforme. Le `WebSocketHandler` valide le JWT, charge le projet via `ProjectService.GetProject`, vérifie la liste des membres. Si l'user n'est pas membre → `403 forbidden`, l'upgrade WS n'a pas lieu. |
| Mécanisme d'absorption | Scoping du `WSBroadcaster` par `projectID` (clé de la `sync.Map`). Un client n'est inscrit dans la `room` du projet que s'il a passé le check membership en amont. |
| Reproduction | Créer un projet sous A, copier l'ID. Avec le token de B, ouvrir `wscat -c "ws://.../projects/<id>/ws?token=$TOKEN_B"` → erreur 403 dans les logs API, connexion refusée. |

---

## Scénario 5 — La base PostgreSQL devient indisponible pendant un retry

| Situation | Pendant `POST /admin/notifications/failed/{id}/retry`, la BDD tombe (timeout, OOM, restart). |
|---|---|
| Comportement attendu | L'opération retourne une erreur claire à l'API admin. Le `failed_notifications` reste cohérent (statut pending, retryCount éventuellement incrémenté ou non selon le moment de la panne). Aucune corruption métier. |
| Comportement observé | ⚠️ Partiellement testé. GORM remonte une erreur claire, l'API renvoie 500. La cohérence dépend de la nature de la panne : la mise à jour de `RetryCount` et l'envoi par le canal ne sont **pas dans une transaction** (le canal email est externe à la BDD). |
| Mécanisme d'absorption | Le canal continue de tenter d'envoyer (ou échoue) indépendamment. La BDD étant la source de vérité du `failed_notifications`, son indisponibilité bloque seulement l'orchestration admin. Les events métier publiés sur le bus ne sont pas affectés (bus en mémoire). |
| Limite | En cas de coupure pendant la persistence d'un nouvel échec, l'incident peut être perdu (le `Dispatcher.recordFailure` log l'erreur de save mais n'a pas de fallback). À corriger avec un buffer en mémoire + retry asynchrone, hors scope de la livraison. |
| Reproduction | `docker compose stop postgres` pendant un retry → 500. `docker compose start postgres` → reprise normale. |

---

## Scénario 6 — Plusieurs déplacements de tâche en concurrence

| Situation | Deux clients (A et B, tous deux membres) déplacent simultanément la même tâche : A clique « Done », B clique « In Progress → Done » au même instant. |
|---|---|
| Comportement attendu | Une seule transition réussit. L'autre reçoit une erreur cohérente (`invalid status transition` ou `404`). Aucun état corrompu. |
| Comportement observé | ⚠️ Comportement « last-write-wins » sans contrôle d'optimistic locking. Le `TaskService.MoveTask` charge → mute → save sans version. La transition de A peut être écrasée par celle de B selon l'ordre d'exécution. Les events sont publiés dans l'ordre des saves. |
| Mécanisme d'absorption (partiel) | Les transitions sont vérifiées par le Value Object `TaskStatus.CanTransitionTo` au moment du domaine. Si B mute après A, le domaine valide la transition à partir de l'état lu (potentiellement obsolète). |
| Limite assumée | Pas d'optimistic locking (champ `version` GORM) pour cette livraison. À ajouter pour de la prod (ADR ultérieur). En revanche, **les events publiés restent ordonnés et cohérents** par rapport à la BDD au moment du save. |
| Reproduction | Difficile à reproduire à la main fiablement ; nécessite un load-test. À noter dans la liste des dettes techniques. |

---

## Synthèse

| # | Scénario | Statut |
|---|---|---|
| 1 | Email en panne | ✅ Géré |
| 2 | Retry après rétablissement | ✅ Géré |
| 3 | Client mobile pendant déploiement v1 | ✅ Géré par construction |
| 4 | WS sur projet non-membre | ✅ Refusé |
| 5 | BDD indisponible pendant retry | ⚠️ Géré, raccourci documenté |
| 6 | Race condition sur déplacement | ⚠️ Last-write-wins, dette identifiée |

Aucun scénario ne nécessite de modification du domaine `task` ou `project`. Les ajustements futurs (auto-retry, optimistic locking, pagination de la dead-letter) se feront tous en couche **infrastructure** ou **application**, pas dans le domaine — c'est ce que valide l'architecture hexagonale livrée à la Disruption #1.
