# Disruption #1 — Évolutions produit

**De :** product@taskflow.io  
**À :** équipe technique  
**Objet :** Nouvelles priorités — 5 chantiers produit + contrainte ops  

---

Bonjour,

Suite au retour de nos premiers utilisateurs bêta, nous avons cinq chantiers à lancer en parallèle. Je vous demande de les intégrer dans la prochaine livraison.

---

**1. Authentification**

L'identifiant utilisateur simulé utilisé jusqu'ici doit être remplacé par une vraie authentification JWT (inscription, connexion, token).

Exigence de durabilité : nous envisageons de migrer vers un SSO d'entreprise dans 12 à 18 mois. Ce changement de mécanisme d'authentification **ne doit pas nécessiter de modifier la logique métier** de l'application (gestion des tâches, des projets, des membres). Documentez votre choix d'implémentation en conséquence.

---

**2. Temps réel sur le Kanban**

Quand un membre déplace une tâche, tous les autres membres du projet doivent voir le changement instantanément, sans recharger la page.

Exigence de durabilité : nous ne savons pas encore quelle technologie temps réel sera la plus adaptée à nos contraintes d'infrastructure. **Ce choix technique ne doit pas impacter la logique de déplacement de tâche** — si nous changeons de mécanisme dans 6 mois, les règles métier ne doivent pas être réécrites. Documentez votre choix et les alternatives envisagées.

Exigence produit : **les mises à jour temps réel sont scopées par projet**. Un membre connecté sur le projet A ne doit pas recevoir les événements du projet B.

---

**3. Système de notifications événementiel**

Nous voulons notifier les membres lors des événements métier importants :

- `task.assigned` → notifier le membre assigné
- `task.moved` → notifier tous les membres du projet

Les canaux initiaux sont : **email** et **notification in-app**.

Exigence de durabilité : nous prévoyons d'ajouter Slack, Teams et SMS dans les 6 prochains mois, au fil des demandes clients. **Chaque ajout d'un nouveau canal ne doit pas nécessiter de modifier les services qui déclenchent les notifications** (gestion des tâches, des projets). Documentez comment votre conception rend cela possible.

Exigence produit : chaque membre peut **activer ou désactiver les canaux** de son choix (ex. : désactiver les emails mais garder les notifications in-app). Un endpoint API pour gérer ces préférences suffit — pas besoin d'interface frontend pour cette livraison.

---

**4. Audit trail**

Nous devons pouvoir répondre à la question : "Qui a fait quoi, sur quelle entité, et à quelle heure ?"

Toute action d'écriture sur les entités `Task` et `Project` doit être tracée automatiquement. Cette traçabilité ne doit pas polluer la logique métier des services.

Exigence de durabilité : nous envisageons de changer le backend de stockage des logs (base relationnelle, service dédié, fichiers). **Ce changement ne doit pas toucher la logique métier** de gestion des tâches et des projets. Expliquez dans votre documentation comment votre implémentation garantit cela.

---

**5. CLI d'administration**

Pour faciliter les démonstrations et la génération de données de test, nous avons besoin d'une interface en ligne de commande permettant de créer des projets, des tâches, et un jeu de données de démonstration sans passer par l'API REST.

Les opérations attendues :

- Créer un projet (nom)
- Créer une tâche dans un projet (titre, projectId)
- Générer un jeu de données de démonstration : un projet avec plusieurs tâches réparties dans les trois colonnes du Kanban (TODO, IN_PROGRESS, DONE)

L'interface d'invocation (commande, script, binaire…) est à votre choix et doit être documentée dans le README.

Exigence de durabilité : l'ajout de cette CLI **ne doit pas nécessiter de dupliquer la logique de création de projets ou de tâches**. Si votre architecture vous oblige à réécrire cette logique pour l'exposer en CLI, documentez pourquoi et ce que cela révèle de votre conception.

---

**6. Mise en production initiale**

La prochaine livraison doit pouvoir être démarrée depuis un clone propre du dépôt avec une seule commande. L'équipe ops ne doit pas avoir à installer Node, configurer la base de données manuellement ou modifier des fichiers de configuration.

Contrainte de livraison : `docker compose up` doit démarrer l'ensemble du système (backend, base de données, frontend). Les variables d'environnement sensibles passent par un `.env` non versionné, documenté par un `.env.example`.

---

Merci de bien documenter vos choix techniques. Nous aurons besoin de les relire dans 6 mois.

**Une dernière chose :** nous n'attendons pas une implémentation production-ready sur cette livraison. Si vous faites des choix simplifiés — stocker les logs en base plutôt qu'un vrai service d'audit, simuler l'envoi d'email par un log console, implémenter un mécanisme temps réel minimal — documentez-le et expliquez ce qui devrait évoluer en production. Un raccourci assumé et documenté vaut mieux qu'une implémentation complète mal structurée.

Cordialement,
Product Owner — TaskFlow
