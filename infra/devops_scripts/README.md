# Scripts DevOps

Scripts pour le CI/CD, les migrations de base de données et le déploiement en production.

## Scripts disponibles

### deploy-prod.sh

Script de déploiement en production pour les services OmniSphere.

```bash
# Déployer un service spécifique
./deploy-prod.sh auth-service

# Déployer tous les services
./deploy-prod.sh all
```

### run-migrations.sh

Script pour exécuter les migrations de base de données.

```bash
# Exécuter les migrations d'un service
./run-migrations.sh auth-service up

# Exécuter toutes les migrations
./run-migrations.sh all up

# Rollback
./run-migrations.sh auth-service down
```

## Prérequis

- AWS CLI configuré
- Accès au cluster ECS
- Accès à la base de données
- Outils de migration (golang-migrate, etc.)

## Configuration

Les scripts utilisent des variables d'environnement :
- `AWS_REGION` - Région AWS
- `DB_HOST` - Hôte de la base de données
- `DB_PORT` - Port de la base de données
- `DB_NAME` - Nom de la base de données
- `DB_USER` - Utilisateur de la base de données
- `DB_PASSWORD` - Mot de passe de la base de données

