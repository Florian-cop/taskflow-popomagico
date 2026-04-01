# Glossaire — Langage Ubiquitaire

Le DDD impose un langage commun entre le code et le métier. Ce glossaire liste les termes à utiliser partout : noms de classes, méthodes, événements, routes.

## Termes métier

| Concept | A utiliser | Définition du concept |
| --- | --- | --- |
| Déplacer une tâche | `MoveTask`, `task.moved` | Transition d'une tâche d'un `TaskStatus` à un autre selon les règles métier. |
| Membre d'un projet | `Member`, `member.added` | Personne rattachée à un projet, pouvant être assignée aux tâches de ce projet. |
| Créer un projet | `CreateProject`, `project.created` | Initialisation d'un nouveau projet avec ses informations de base. |
| Assigner une tâche | `AssignTask`, `task.assigned` | Association d'une tâche à un membre responsable de son exécution. |
| Projet | `Project` | Agrégat racine qui regroupe les membres et les tâches. |
| Tâche | `Task` | Unité de travail suivie dans le flux Kanban avec statut et responsable. |
| Colonne Kanban | `TaskStatus` (Todo, InProgress, Done) | Représentation visuelle d'un statut de tâche dans le tableau Kanban. |

## Convention de nommage des événements

Format : **`entité.action`** au passé implicite.

- `task.created` — une tâche a été créée
- `task.moved` — une tâche a été déplacée
- `task.assigned` — une tâche a été assignée
- `member.added` — un membre a été ajouté
- `project.created` — un projet a été créé
