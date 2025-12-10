# Catalogue Service

Service de gestion des produits et des stocks pour OmniSphere.

## Fonctionnalités

- Gestion des produits
- Gestion des stocks
- Synchronisation avec Elasticsearch
- Recherche de produits

## Développement

```bash
go mod download
go run main.go
```

## Endpoints

- `GET /health` - Health check
- `GET /ready` - Readiness check
- `GET /api/v1/products` - Liste des produits
- `GET /api/v1/products/:id` - Détails d'un produit
- `POST /api/v1/products` - Créer un produit
- `PUT /api/v1/products/:id` - Mettre à jour un produit
- `DELETE /api/v1/products/:id` - Supprimer un produit
- `GET /api/v1/inventory/:productId` - Récupérer le stock
- `PUT /api/v1/inventory/:productId` - Mettre à jour le stock
- `POST /api/v1/search` - Rechercher des produits

## Configuration

Variables d'environnement:
- `ELASTICSEARCH_URL` - URL d'Elasticsearch (défaut: http://localhost:9200)

