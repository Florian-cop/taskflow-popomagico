# Disruption #2 — Résilience & maturité produit

## Message 1 — Incident de production

**Canal :** Slack — #incidents  
**De :** ops@taskflow.io  
**Heure :** 08h47

L'API email est tombée ce matin. Résultat : le handler de notifications email plante et entraîne des erreurs en cascade dans tout le système de traitement d'événements.

**Exigence immédiate à implémenter :**

- si un canal de notification est défaillant, les autres canaux doivent continuer à fonctionner ;
- le reste du système ne doit pas remonter d'erreur visible à l'utilisateur ;
- l'échec du canal email doit être conservé dans une file de messages à retraiter plus tard.

Pour la démonstration, vous devez pouvoir simuler une exception dans `EmailChannel` et montrer que `InApp`, `Realtime` et `Audit` continuent de fonctionner.

## Message 2 — Nouvelle contrainte client

**Canal :** Email  
**De :** product@taskflow.io  
**Objet :** Application mobile tierce

Nous avons signé avec un partenaire qui va développer une application mobile consommant notre API. Ils ne pourront pas mettre à jour leur application aussi vite que nous faisons évoluer la nôtre.

**Exigence à implémenter :** versionner l'API REST avec deux versions qui coexistent.

Pour cette livraison, une version simple suffit :

- `/api/v1/...` continue de fonctionner pour les clients existants ;
- `/api/v2/...` expose les mêmes cas d'usage, éventuellement avec un format de réponse légèrement différent ;
- les deux versions doivent réutiliser la logique métier existante.

L'objectif n'est pas de dupliquer les services, mais de montrer que l'API est une couche d'entrée distincte du domaine métier.

## Message 3 — Évolution produit à analyser : multi-workspace

**Canal :** Email  
**De :** product@taskflow.io  
**Objet :** Isolation des données par entreprise

Notre modèle commercial évolue. Chaque entreprise cliente doit avoir son propre espace isolé (workspace). Un utilisateur peut appartenir à plusieurs workspaces. Les données d'un workspace ne doivent jamais être visibles depuis un autre.

Un utilisateur se connecterait, verrait la liste de ses workspaces, sélectionnerait l'un d'eux, puis accéderait à ses projets et tâches.

**Exigence pour cette séance : analyse d'impact uniquement.**

Vous ne devez pas implémenter complètement le multi-workspace dans les 2 heures. Vous devez expliquer :

- quelles couches seraient modifiées ;
- où placer le `workspaceId` ;
- comment garantir l'isolation côté API ;
- pourquoi un simple filtre frontend serait insuffisant ;
- comment éviter de réécrire `TaskService` et `ProjectService`.

L'objectif est de vérifier que l'architecture issue de la disruption 1 peut absorber cette évolution sans déplacer la logique métier dans les controllers ou les handlers techniques.

## Livrables attendus pour cette disruption

1. Résilience des consommateurs de notifications.
2. Versioning API simple avec `/api/v1` et `/api/v2`.
3. Analyse d'impact écrite pour le multi-workspace.

Les choix simplifiés sont acceptés s'ils sont explicitement documentés.
