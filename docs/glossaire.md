# Glossaire — Langage Ubiquitaire

Le DDD impose un langage commun entre le code et le métier. Ce glossaire liste les termes à utiliser partout : noms de classes, méthodes, événements, routes.

## Termes métier

| Concept | A utiliser | A éviter |
| --- | --- | --- |
| Déplacer une tâche | `MoveTask`, `task.moved` | `UpdateStatus`, `ChangeState` |
| Membre d'un projet | `Member`, `member.added` | `User`, `Participant` |
| Créer un projet | `CreateProject`, `project.created` | `AddProject`, `NewProject` |
| Assigner une tâche | `AssignTask`, `task.assigned` | `SetAssignee`, `UpdateAssignee` |
| Projet | `Project` | `Board`, `Workspace` |
| Tâche | `Task` | `Ticket`, `Issue`, `Card` |
| Colonne Kanban | `TaskStatus` (Todo, InProgress, Done) | `State`, `Phase`, `Column` |

## Convention de nommage des événements

Format : **`entité.action`** au passé implicite.

- `task.created` — une tâche a été créée
- `task.moved` — une tâche a été déplacée
- `task.assigned` — une tâche a été assignée
- `member.added` — un membre a été ajouté
- `project.created` — un projet a été créé
