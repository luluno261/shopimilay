# API Gateway

Service API Gateway centralisé pour OmniSphere qui route les requêtes vers les microservices backend.

## Fonctionnalités

- Routage des requêtes vers les services backend
- Authentification JWT centralisée
- Rate limiting
- CORS
- Health checks

## Développement

```bash
go mod download
go run main.go
```

## Configuration

Variables d'environnement :
- `PORT` - Port d'écoute (défaut: 8080)
- `JWT_SECRET` - Clé secrète pour valider les tokens JWT
- `AUTH_SERVICE_HOST` - Host du service auth (défaut: localhost)
- `AUTH_SERVICE_PORT` - Port du service auth (défaut: 8080)
- `CHECKOUT_SERVICE_HOST` - Host du service checkout (défaut: localhost)
- `CHECKOUT_SERVICE_PORT` - Port du service checkout (défaut: 8081)
- `CATALOGUE_SERVICE_HOST` - Host du service catalogue (défaut: localhost)
- `CATALOGUE_SERVICE_PORT` - Port du service catalogue (défaut: 8082)
- `MARKETING_ENGINE_HOST` - Host du service marketing (défaut: localhost)
- `MARKETING_ENGINE_PORT` - Port du service marketing (défaut: 8083)

## Endpoints

Tous les endpoints sont préfixés par `/api/v1`.

### Routes publiques
- `POST /auth/login` - Connexion
- `POST /auth/register` - Inscription
- `POST /auth/refresh` - Rafraîchir le token
- `GET /products` - Liste des produits
- `GET /products/:id` - Détail d'un produit
- `POST /search` - Recherche de produits

### Routes protégées
Toutes les autres routes nécessitent un header `Authorization: Bearer <token>`

