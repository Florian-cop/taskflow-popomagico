# Grille d'évaluation — TaskFlow

---

L'évaluation distingue deux axes :

- **Qualité architecturale** : séparation des responsabilités, stabilité du domaine, découplage, tests, ADR
- **Complétude et maturité du système** : fonctionnalités livrées, démonstration, qualité opérationnelle, robustesse

Une couverture fonctionnelle incomplète mais portée par une architecture cohérente peut être mieux évaluée qu'un rendu plus complet mais fortement couplé ou réécrit à chaque disruption.

---

## Rendu 1 — Fondations *(coefficient 20 %)*

| Critère | Détail | Points |
|---------|--------|--------|
| Séparation controller / service / repository | Le controller ne contient pas de logique métier. Le service ne connaît pas l'ORM directement. | /5 |
| Abstraction repository effective | Une interface définit le contrat. L'implémentation ORM est injectée. Remplacer l'ORM ne touche pas les services. | /5 |
| Tests sans dépendance infrastructure | Les tests unitaires utilisent des mocks. Aucun appel réel à la BDD. Les transitions de statut et la publication d'events sont couvertes. | /5 |
| Domain events dès la Phase 1 | Au moins deux events sont publiés (`task.created`, `task.moved`). Un `ConsoleHandler` y souscrit et affiche event + taskId + horodatage. Le service ne connaît pas ses consommateurs. | /5 |
| **Total** | | **/20** |

---

## Rendu 2 — Évolution *(coefficient 25 %)*

| Critère | Détail | Points |
|---------|--------|--------|
| Capacité d'absorption sans réécriture | Les couches stables (domaine, services existants) n'ont pas été réécrites. L'analyse d'impact le démontre. | /5 |
| Audit trail non-intrusif | L'audit est ajouté sans modifier la logique métier des services. Il se branche sur les domain events existants. | /5 |
| Docker Compose fonctionnel | `docker compose up` depuis un clone propre démarre tous les services sans intervention manuelle. | /5 |
| Pipeline CI/CD | La pipeline est verte sur `main`. Les tests et le build sont automatisés. | /5 |
| **Total** | | **/20** |

---

## Rendu 3 — Résilience *(coefficient 30 %)*

| Critère | Détail | Points |
|---------|--------|--------|
| Résilience démontrée | La panne d'un canal n'arrête pas le système. Les messages échoués sont mis en file. La démonstration live le prouve. | /5 |
| Multi-workspace | L'isolation des données entre workspaces est garantie par l'architecture (guard + scoping des repositories), pas uniquement côté frontend. | /5 |
| API rétrocompatible | Les deux versions de l'API coexistent. Un client v1 fonctionne encore après l'ajout de la v2. | /5 |
| Maturité globale | Cohérence de l'ensemble, qualité des ADR, test des scénarios de panne, maturité opérationnelle du système. | /5 |
| **Total** | | **/20** |

---

## Soutenance *(coefficient 25 %)*

| Critère | Détail | Points |
|---------|--------|--------|
| Démo live | `docker compose up` fonctionne depuis zéro. Le Kanban temps réel est visible sur deux navigateurs. | /5 |
| Preuve de stabilité | Des fichiers précis sont montrés comme n'ayant jamais été modifiés entre les phases. | /5 |
| Analyse critique | L'équipe identifie ses bons et mauvais choix. Les ADR sont défendus avec conviction. | /5 |
| Réponses aux questions | Maîtrise des concepts : repository pattern, Open/Closed, circuit breaker, interceptor. | /5 |
| **Total** | | **/20** |

---

## Bonus et malus

| | Condition | Points |
|-|-----------|--------|
| **Bonus** | ADR de qualité exceptionnelle (analyse, alternatives réelles, conséquences détaillées) | +1 |
| **Bonus** | Métriques `git diff --stat` impressionnantes (moins de 5 fichiers *existants* modifiés pour absorber une disruption — les nouveaux fichiers ne comptent pas) | +1 |
| **Bonus** | Tests couvrant les scénarios de panne (circuit breaker, mock de canal défaillant) | +1 |
| **Malus** | ORM importé directement dans les services | −2 |
| **Malus** | Réécriture complète suite à une disruption (aucune couche stable) | −2 |
| **Malus** | Aucun test unitaire dans le projet | −2 |
| **Malus** | `docker compose up` ne démarre pas lors de la soutenance | −2 |

---

## Note finale

$$\text{Note} = R1 \times 0{,}20 + R2 \times 0{,}25 + R3 \times 0{,}30 + \text{Soutenance} \times 0{,}25 + \text{Bonus/Malus}$$
