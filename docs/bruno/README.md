# Collection Bruno — TaskFlow API

Collection [Bruno](https://www.usebruno.com/) couvrant l'ensemble des endpoints exposés par `taskflow-api` (v1 + v2).

## Utilisation

1. Ouvrir Bruno → **Open Collection** → sélectionner ce dossier (`docs/bruno`).
2. Sélectionner l'environnement **Local** (en haut à droite).
3. Lancer l'API : `docker compose up` (ou `go run ./cmd/api` depuis `taskflow-api/`).
4. Exécuter **Auth / Register** puis **Auth / Login** — le token JWT est automatiquement stocké dans `{{token}}` via un script post-response.
5. Les autres requêtes envoient automatiquement le header `Authorization: Bearer {{token}}`.

## Variables d'environnement

| Variable          | Usage                                                 |
| ----------------- | ----------------------------------------------------- |
| `baseUrl`         | Racine de l'API v1 (`http://localhost:8080/api/v1`)   |
| `baseUrlV2`       | Racine de l'API v2                                    |
| `wsBaseUrl`       | Racine WebSocket                                      |
| `token`           | JWT renseigné automatiquement après `Login`           |
| `projectId`       | ID de projet utilisé par les requêtes Tasks/Members   |
| `taskId`          | ID de tâche utilisé par `Move Task`                   |
| `userId`          | ID utilisateur (pour ajout de membre)                 |
| `notificationId`  | ID de notification (pour `Mark as read`)              |

## Organisation

- **Auth/** — register, login, me
- **Users/** — recherche & lookup
- **Projects/** — CRUD projets + ajout de membre
- **Tasks/** — listing, création, déplacement (Kanban)
- **Notifications/** — listing, lecture, préférences
- **Audit/** — consultation du journal
- **Admin/** — gestion des canaux de notification & retry des échecs (chantier résilience)
- **WebSocket/** — endpoint temps réel (à ouvrir dans un client WS)
- **v2/** — endpoints v2 (réponses enveloppées `{data, meta}`)
