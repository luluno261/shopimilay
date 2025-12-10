# Migrations de base de données

Fichiers de migration pour gérer l'évolution du schéma de base de données.

## Structure

Les migrations sont organisées par service :
- `auth-service/` - Migrations pour le service d'authentification
- `checkout-service/` - Migrations pour le service de checkout
- `catalogue-service/` - Migrations pour le service de catalogue

## Format

Les migrations peuvent être :
- Fichiers SQL (`.sql`)
- Migrations Go (utilisant golang-migrate)

## Utilisation

Voir les scripts dans `infra/devops_scripts/run-migrations.sh` pour exécuter les migrations.

## Outils recommandés

- **golang-migrate** : Pour les services Go
- **Knex.js** : Pour les services Node.js
- **Flyway** : Alternative SQL pure

