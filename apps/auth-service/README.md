# Auth Service

Service d'authentification pour OmniSphere.

## Fonctionnalités

- Authentification JWT
- Gestion des rôles et permissions
- Gestion des utilisateurs de la plateforme
- Refresh tokens

## Développement

```bash
# Installer les dépendances
go mod download

# Lancer le service
go run main.go

# Builder
go build -o auth-service main.go
```

## Endpoints

- `GET /health` - Health check
- `GET /ready` - Readiness check
- `POST /api/v1/auth/login` - Connexion
- `POST /api/v1/auth/register` - Inscription
- `POST /api/v1/auth/refresh` - Rafraîchir le token
- `GET /api/v1/auth/me` - Informations utilisateur connecté

## Configuration

Voir `config.yaml` pour la configuration du service.
